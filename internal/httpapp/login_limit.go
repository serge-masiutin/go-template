package httpapp

import (
	"net"
	"sync"
	"time"
)

type loginWindow struct {
	started time.Time
	count   int
}
type loginLimiter struct {
	mu      sync.Mutex
	windows map[string]loginWindow
}

func newLoginLimiter() *loginLimiter { return &loginLimiter{windows: make(map[string]loginWindow)} }

// Allow limits CPU-expensive password checks by the direct peer. Forwarded
// headers are deliberately not trusted; deployments also limit at the proxy.
func (l *loginLimiter) Allow(remote string, now time.Time) bool {
	host, _, err := net.SplitHostPort(remote)
	if err != nil {
		return false
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	window, found := l.windows[host]
	if !found || now.Sub(window.started) >= time.Minute {
		if !found && len(l.windows) >= 4096 {
			for key, entry := range l.windows {
				if now.Sub(entry.started) >= time.Minute {
					delete(l.windows, key)
				}
			}
			if len(l.windows) >= 4096 {
				return false
			}
		}
		window = loginWindow{started: now}
	}
	if window.count >= 10 {
		return false
	}
	window.count++
	l.windows[host] = window
	return true
}
