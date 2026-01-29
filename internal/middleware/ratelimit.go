package middleware

import (
	"strconv"
	"time"
	"sync"
	"net/http"
)

type clientRecord struct {
	count int
	windowStart time.Time
}

var (
	mu sync.Mutex
	clients = make(map[string]*clientRecord)
)

const (
	maxRequests = 5
	windowSize = 10 * time.Second
)

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

		// userID := r.Header.Get("X-User-ID")
		// if userID == "" {
		// 	http.Error(w, "user not Identified", http.StatusUnauthorized)
		// 	return
		// }

		userID, ok := r.Context().Value(userIDKey).(string)
		if !ok || userID == "" {
			http.Error(w, "user not identified", http.StatusUnauthorized)
			return
		}

		mu.Lock()
		record, exists := clients[userID]
		now := time.Now()

		if !exists {
			clients[userID] = &clientRecord{
				count: 1,
				windowStart: now,
			}
			mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		if now.Sub(record.windowStart) > windowSize {
			record.count = 1
		    record.windowStart = now
			mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		if record.count >= maxRequests {
			retryAfter := int(windowSize.Seconds() - time.Since(record.windowStart).Seconds())
			if retryAfter < 0 {
				retryAfter = 0
			}

			w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			mu.Unlock()
			
			return
		}
		record.count++
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}