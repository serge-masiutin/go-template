# Extract Query Logic to Query Object

Pull a long database query out of a handler into a parameterized query object.

## Before

```go
// BAD: HTTP handler builds report SQL and interprets persistence rows.
rows, err := pool.Query(r.Context(), reportSQL,
    r.URL.Query().Get("start"), r.URL.Query().Get("end"), r.URL.Query().Get("region"))
if err != nil { return err }
defer rows.Close()
// Transport parsing and reporting concerns are coupled.
```

## After

```go
type SalesFilter struct { From, Until time.Time; Region string }
type DailySales struct { Day time.Time; OrderCount int64; RevenueCents int64 }
type SalesReport struct { pool *pgxpool.Pool }
func (q *SalesReport) Daily(ctx context.Context, filter SalesFilter) ([]DailySales, error) {
    const statement = `SELECT date_trunc('day', o.created_at),
        count(DISTINCT o.id), sum(i.quantity * i.price_cents)
        FROM orders o JOIN customers c ON c.id=o.customer_id
        JOIN line_items i ON i.order_id=o.id
        WHERE o.status='completed' AND o.created_at >= $1 AND o.created_at < $2
        AND ($3='' OR c.region=$3)
        GROUP BY 1 ORDER BY 1 DESC`
    rows, err := q.pool.Query(ctx, statement, filter.From, filter.Until, filter.Region)
    if err != nil { return nil, err }
    defer rows.Close()
    result := make([]DailySales, 0)
    for rows.Next() {
        var day DailySales
        if err := rows.Scan(&day.Day, &day.OrderCount, &day.RevenueCents); err != nil { return nil, err }
        result = append(result, day)
    }
    return result, rows.Err()
}
```
