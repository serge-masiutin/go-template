package observability

import (
	"context"
	"crypto/rand"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/serge-masiutin/go-template/internal/database"
)

type Metrics struct {
	registry *prometheus.Registry
	requests *prometheus.HistogramVec
	logger   *slog.Logger
}

func New(db *database.DB, logger *slog.Logger) *Metrics {
	registry := prometheus.NewRegistry()
	requests := prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "http_request_duration_seconds", Help: "HTTP request duration by registered route and response code.", Buckets: prometheus.DefBuckets}, []string{"route", "status"})
	registry.MustRegister(requests, collectors.NewGoCollector(), collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	registry.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{Name: "database_connections", Help: "Total application pgx connections."}, func() float64 { return float64(db.Stat().TotalConns()) }))
	registry.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{Name: "database_connections_in_use", Help: "Acquired application pgx connections."}, func() float64 { return float64(db.Stat().AcquiredConns()) }))
	registry.MustRegister(queueCollector{db: db, jobs: prometheus.NewDesc("background_jobs", "Persisted River jobs by queue and state.", []string{"queue", "state"}, nil)})
	return &Metrics{registry: registry, requests: requests, logger: logger}
}
func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}
func (m *Metrics) HTTP(route string, next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := rand.Text()
		w.Header().Set("X-Request-ID", requestID)
		measured := httpsnoop.CaptureMetrics(next, w, r)
		m.requests.WithLabelValues(route, strconv.Itoa(measured.Code)).Observe(measured.Duration.Seconds())
		m.logger.InfoContext(r.Context(), "request", "request_id", requestID, "route", route, "status", measured.Code, "duration_ms", measured.Duration.Milliseconds())
	})
}

type queueCollector struct {
	db   *database.DB
	jobs *prometheus.Desc
}

func (c queueCollector) Describe(ch chan<- *prometheus.Desc) { ch <- c.jobs }
func (c queueCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	rows, err := c.db.Query(ctx, "SELECT CASE WHEN queue IN ('mail','ai','default') THEN queue ELSE 'other' END, state::text, count(*) FROM river_job GROUP BY 1,2")
	if err != nil {
		ch <- prometheus.NewInvalidMetric(c.jobs, errors.New("queue metrics query failed"))
		return
	}
	defer rows.Close()
	for rows.Next() {
		var queue, state string
		var count float64
		if err := rows.Scan(&queue, &state, &count); err != nil {
			ch <- prometheus.NewInvalidMetric(c.jobs, err)
			return
		}
		ch <- prometheus.MustNewConstMetric(c.jobs, prometheus.GaugeValue, count, queue, state)
	}
	if err := rows.Err(); err != nil {
		ch <- prometheus.NewInvalidMetric(c.jobs, err)
	}
}
