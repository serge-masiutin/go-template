package httpapp

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	_ "embed"
	"encoding/hex"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
	inertia "github.com/romsar/gonertia/v3"
	"github.com/serge-masiutin/go-template/internal/accounts"
	"github.com/serge-masiutin/go-template/internal/assistant"
	"github.com/serge-masiutin/go-template/internal/background"
	"github.com/serge-masiutin/go-template/internal/config"
	"github.com/serge-masiutin/go-template/internal/database"
	"github.com/serge-masiutin/go-template/internal/diagnostics"
	"github.com/serge-masiutin/go-template/internal/notemail"
	"github.com/serge-masiutin/go-template/internal/notes"
	"github.com/serge-masiutin/go-template/internal/observability"
	"gorm.io/gorm"
)

//go:embed root.html
var rootHTML string

type App struct {
	cfg        config.Config
	emails     *notemail.Service
	assistant  *assistant.Service
	pool       *database.DB
	accounts   *accounts.Store
	notes      *notes.Store
	sessions   *scs.SessionManager
	inertia    *inertia.Inertia
	logger     *slog.Logger
	loginLimit *loginLimiter
}

func New(cfg config.Config, pool *database.DB, sessions *scs.SessionManager, logger *slog.Logger) (http.Handler, error) {
	users, err := accounts.New(pool)
	if err != nil {
		return nil, err
	}
	hotFile := "tmp/vite.hot"
	options := []inertia.Option{inertia.WithFlashProvider(flashProvider{sessions}), inertia.WithEncryptHistory()}
	if cfg.Production || cfg.Environment == "test" {
		hotFile = "" // Production and browser tests use built assets, even alongside bin/dev.
	}
	_, hotErr := os.Stat(hotFile)
	if hotFile == "" || errors.Is(hotErr, os.ErrNotExist) {
		options = append(options, inertia.WithVersionFromFile("web/build/manifest.json"))
	} else if hotErr != nil {
		return nil, hotErr
	} else {
		contents, err := os.ReadFile(hotFile)
		if err != nil {
			return nil, err
		}
		if strings.TrimSpace(string(contents)) != "http://localhost:5173" {
			return nil, errors.New("unexpected Vite development origin")
		}
	}
	engine, err := inertia.New(rootHTML, options...)
	if err != nil {
		return nil, err
	}
	_, err = inertia.NewVite(engine, inertia.WithHotFile(hotFile), inertia.WithBuildManifest("web/build/manifest.json"), inertia.WithBuildDir("/build/"))
	if err != nil {
		return nil, err
	}
	queue, err := background.Producer(pool, logger)
	if err != nil {
		return nil, err
	}
	metrics := observability.New(pool, logger)
	app := &App{cfg: cfg, emails: notemail.New(pool, queue, users), assistant: assistant.New(pool, queue), pool: pool, accounts: users, notes: notes.New(pool), sessions: sessions, inertia: engine, logger: logger, loginLimit: newLoginLimiter()}
	sessions.ErrorFunc = func(w http.ResponseWriter, r *http.Request, err error) { app.fail(w, r, err) }
	pages := http.NewServeMux()
	register := func(pattern string, handler http.HandlerFunc) { pages.Handle(pattern, metrics.HTTP(pattern, handler)) }
	register("GET /{$}", app.home)
	register("GET /login", app.loginPage)
	register("POST /login", app.login)
	register("POST /logout", app.logout)
	register("POST /notes", app.createNote)
	register("DELETE /notes/{id}", app.deleteNote)
	register("GET /admin", app.admin)
	register("GET /tools", app.toolsPage)
	register("POST /tools/email", app.emailNotes)
	register("POST /tools/assistant", app.askAssistant)
	mux := http.NewServeMux()
	if cfg.MetricsToken != "" {
		mux.Handle("GET /metrics", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if subtle.ConstantTimeCompare([]byte(r.Header.Get("Authorization")), []byte("Bearer "+cfg.MetricsToken)) != 1 {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			metrics.Handler().ServeHTTP(w, r)
		}))
	}
	mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"status\":\"ok\"}\n"))
	})
	mux.HandleFunc("GET /health/ready", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := pool.Ping(ctx); err != nil {
			logger.ErrorContext(ctx, "database readiness failed")
			http.Error(w, "unavailable", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"status\":\"ok\"}\n"))
	})
	mux.Handle("GET /build/", http.StripPrefix("/build/", http.FileServer(http.Dir("web/build"))))
	mux.Handle("GET /fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("web/public/fonts"))))
	mux.Handle("/", sessions.LoadAndSave(app.recoverPanic(app.csrf(engine.Middleware(pages)))))
	return app.security(cfg, http.NewCrossOriginProtection().Handler(mux)), nil
}

func (a *App) csrf(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Precognition") != "" {
			http.Error(w, "Precognition is not supported", http.StatusBadRequest)
			return
		}
		token := a.sessions.GetString(r.Context(), "csrf")
		if token == "" {
			bytes := make([]byte, 32)
			if _, err := rand.Read(bytes); err != nil {
				a.fail(w, r, err)
				return
			}
			token = hex.EncodeToString(bytes)
			a.sessions.Put(r.Context(), "csrf", token)
		}
		r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
		if r.Method != "GET" && r.Method != "HEAD" && r.Method != "OPTIONS" {
			supplied := r.Header.Get("X-CSRF-Token")
			if subtle.ConstantTimeCompare([]byte(supplied), []byte(token)) != 1 {
				http.Error(w, "invalid CSRF token", 403)
				return
			}
		}
		ctx := inertia.SetProp(r.Context(), "csrfToken", inertia.Always(token))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *App) security(cfg config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "same-origin")
		w.Header().Set("Cache-Control", "no-store")
		if cfg.Production {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000")
			w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'")
		}
		next.ServeHTTP(w, r)
	})
}

func (a *App) user(w http.ResponseWriter, r *http.Request) (accounts.User, bool) {
	id := a.sessions.GetInt64(r.Context(), "userID")
	if id == 0 {
		a.inertia.Redirect(w, r, "/login", http.StatusSeeOther)
		return accounts.User{}, false
	}
	user, err := a.accounts.Find(r.Context(), id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := a.sessions.Destroy(r.Context()); err != nil {
			a.fail(w, r, err)
			return accounts.User{}, false
		}
		a.inertia.Redirect(w, r, "/login", http.StatusSeeOther)
		return accounts.User{}, false
	}
	if err != nil {
		a.fail(w, r, err)
		return accounts.User{}, false
	}
	return user, true
}

func (a *App) render(w http.ResponseWriter, r *http.Request, page string, props inertia.Props) {
	if err := a.inertia.Render(w, r, page, props); err != nil {
		a.fail(w, r, err)
	}
}

func (a *App) fail(w http.ResponseWriter, r *http.Request, err error) {
	a.logger.With("route", r.Pattern).ErrorContext(r.Context(), "request failed", diagnostics.Attributes(err)...)
	http.Error(w, "internal server error", 500)
}

// Gonertia treats returned flash-provider errors as optional. An invalid persisted
// flash value is a broken session contract; fail the request instead of losing it.
func (a *App) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				if recovered == http.ErrAbortHandler {
					panic(recovered)
				}
				a.logger.ErrorContext(r.Context(), "request panic", "stack", string(debug.Stack()))
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
