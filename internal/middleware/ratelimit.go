package middleware

import (
	// "strconv" -- this was for in memory rate limiting
	"time"
	//"sync" -- this was for in memory rate limiting
	"net/http"
	"fmt"
	"api-gateway/internal/redisclient"
	"api-gateway/internal/metrics"

)
//-------- this is code for in memory rate limiting ---------
// type clientRecord struct {
// 	count int
// 	windowStart time.Time
// }

// var (
// 	mu sync.Mutex
// 	clients = make(map[string]*clientRecord)
// )
//---------end of in memory rate limiting -----------

//-------- this is code for redis based rate limiting ---------



const (
	maxRequests = 5
	windowSize = 10 * time.Second
)

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

	//------- this was when we were fetching userID from header directly ------
		// userID := r.Header.Get("X-User-ID")
		// if userID == "" {
		// 	http.Error(w, "user not Identified", http.StatusUnauthorized)
		// 	return
		// }
	//------- end of that -------



		userID, ok := r.Context().Value(userIDKey).(string)
		if !ok || userID == "" {
			http.Error(w, "user not identified", http.StatusUnauthorized)
			return
		}
		
		key := fmt.Sprintf("rate_limit:%s", userID)

		// increment request count atomically in Redis
		count, err := redisclient.Client.
			Incr(redisclient.Ctx, key).
			Result()
		if err != nil {
			http.Error(w, "rate limit error", http.StatusInternalServerError)
			return
		}

		// first request in window → set expiry
		if count == 1 {
			redisclient.Client.Expire(
				redisclient.Ctx,
				key,
				windowSize,
			)
		}

		// reject if limit exceeded
		if count > maxRequests {

			w.Header().Set(
				"Retry-After",
				fmt.Sprintf("%d", int(windowSize.Seconds())),
			)
			metrics.RateLimitedRequestsTotal.
				WithLabelValues(r.URL.Path).
				Inc()

			metrics.HttpRequestsTotal.
				WithLabelValues(r.URL.Path, r.Method, "429").
				Inc()

			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
		metrics.HttpRequestsTotal.
			WithLabelValues(r.URL.Path, r.Method, "200").
			Inc()


		next.ServeHTTP(w, r)

		/*-------earlier in memory rate limiting code --------

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

		------- end of earlier in memory rate limiting code --------*/
	})
}






