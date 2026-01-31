package server

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
	"context"

	"api-gateway/internal/middleware"
	"api-gateway/internal/metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const maxRetries = 1
const retryBackoff = 100 * time.Millisecond


func New() *http.Server {
	mux := http.NewServeMux()

	// Public endpoints (NO auth, NO rate limit)
	mux.Handle("/metrics", promhttp.Handler())

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("gateway is healthy"))
	})

	// Protected routes (will go through middleware)
	mux.Handle("/users", userServiceProxy())
	mux.Handle("/users/", userServiceProxy())

	// Root handler decides what bypasses middleware
	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Metrics must be public
		if r.URL.Path == "/metrics" || r.URL.Path == "/health" {
			mux.ServeHTTP(w, r)
			return
		}

		// Everything else goes through middleware chain
		middleware.AuthMiddleware(
			middleware.RateLimitMiddleware(
				loggingMiddleware(mux),
			),
		).ServeHTTP(w, r)
	})

	return &http.Server{
		Addr:         ":8080",
		Handler:      rootHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.
			WithLabelValues(r.URL.Path, r.Method).
			Observe(duration)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func userServiceProxy() http.Handler {
	target := "http://localhost:8081"

	proxy := httputil.NewSingleHostReverseProxy(mustParseURL(target))
	//originalDirector := proxy.Director


	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = "localhost:8081"
		req.URL.Path = "/users"
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		retries, _ := r.Context().Value("retries").(int)

		if retries < maxRetries {
		//	log.Println("retrying upstream after backoff")
			metrics.RetryAttemptsTotal.
				WithLabelValues(r.URL.Path).
				Inc()
			time.Sleep(retryBackoff)

			ctx := context.WithValue(r.Context(), "retries", retries+1)
			reqCopy := r.Clone(ctx)

			proxy.ServeHTTP(w, reqCopy)
			return
		}

		http.Error(
			w,
			"upstream service unavailable",
			http.StatusBadGateway,
		)
	}


	return proxy
}

func mustParseURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return u
}



