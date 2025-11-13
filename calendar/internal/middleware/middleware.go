package middleware

import (
	"fmt"
	"net/http"
	"time"
)

func (l *Logger) Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		handler.ServeHTTP(w, r)

		duration := time.Since(start)

		logStr := fmt.Sprintf("%s %s %s", r.Method, r.URL.Path, duration)

		l.Ch <- logStr
	})
}
