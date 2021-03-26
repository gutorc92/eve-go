package handlers

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
	"github.com/gutorc92/api-farm/config"
	"github.com/gutorc92/api-farm/dao"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type key int

const (
	requestIDKey key = 0
)

var (
	listenAddr string
	healthy    int32
)

// Server holds the information needed to run Whisper
type Server struct {
	*config.WebConfig
	FarmApis FarmAPI
	Apis     []API
}

// InitFromWebConfig builds a Server instance
func (s *Server) InitFromWebConfig(wc *config.WebConfig) *Server {
	s.WebConfig = wc
	var dt *dao.DataMongo
	dt, err := dao.NewDataMongo(wc.Uri, wc.Database)
	if err != nil {
		panic(err)
	}
	s.FarmApis = new(DefaultFarmAPI).InitConfig(wc, dt)
	batchApi := new(DefaultBatchAPI).InitConfig(wc, dt)
	cowApi := new(DefaultCowAPI).InitConfig(wc, dt)
	s.Apis = append(s.Apis, batchApi)
	s.Apis = append(s.Apis, cowApi)
	return s
}

func (s *Server) Serve() error {

	flag.StringVar(&listenAddr, "listen-addr", ":5000", "server listen address")
	flag.Parse()

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Println("Server is starting...")

	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/healthz", healthz())
	router.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)
	v1Router := router.PathPrefix("/v1").Subrouter()
	v1Router.Handle("/farm", s.FarmApis.GETHandler()).Methods("GET")
	v1Router.Handle("/farm", s.FarmApis.POSTHandler()).Methods("POST")
	for _, api := range s.Apis {
		v1Router.Handle(api.GetUrl(), api.GETHandler()).Methods("GET")
		v1Router.Handle(api.GetUrl(), api.POSTHandler()).Methods("POST")
	}
	return s.ListenAndServe(router)

}

// ListenAndServe fires up the configured http server
func (s *Server) ListenAndServe(router *mux.Router) error {
	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	server := &http.Server{
		Addr:         listenAddr,
		Handler:      tracing(nextRequestID)(logging(logger)(router)),
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Println("Server is shutting down...")
		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	logger.Println("Server is ready to handle requests at", listenAddr)
	atomic.StoreInt32(&healthy, 1)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}

	<-done
	logger.Println("Server stopped")
	return nil
}

func healthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt32(&healthy) == 1 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	})
}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(requestIDKey).(string)
				if !ok {
					requestID = "unknown"
				}
				logger.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
