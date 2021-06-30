package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gutorc92/eve-go/metrics"
)

func MetricsMiddleware(next http.Handler, metrics *metrics.Metrics, url string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Println("Executing middlewareOne")
		next.ServeHTTP(w, r)
		log.Println("Executing middlewareOne again")
		metrics.CountApiCall("http", "200", "false", "", "GET", url, time.Since(start).Seconds())
	})
}

func HeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middlewareOne")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
		log.Println("Executing middlewareOne again")
	})
}
