package server

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"api-gateway/internal/middleware"
)

func New() *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("gateway is healthy"))
	})

	mux.Handle("/users", userServiceProxy())
	mux.Handle("/users/", userServiceProxy())

	return &http.Server{
		Addr:         ":8080",
		Handler:      middleware.AuthMiddleware(middleware.RateLimitMiddleware((loggingMiddleware(mux)))),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		next.ServeHTTP(w, r)

		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func userServiceProxy() http.Handler {
	target := "http://localhost:8081"

	proxy := httputil.NewSingleHostReverseProxy(
		mustParseURL(target),
	)

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = "localhost:8081"
		req.URL.Path = "/users"
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
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
