package internalhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server struct { // TODO
	router *mux.Router
	logger *zap.Logger
	app    Application
}

type Application interface { // TODO
}

func NewServer(logger *zap.Logger, app Application) *Server {
	router := mux.NewRouter()
	s := &Server{
		router: router,
		logger: logger,
		app:    app,
	}
	s.configureRouter()

	return s
}

func (s *Server) configureRouter() {
	s.router.Use(s.loggingMiddleware)
	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info(
			"started",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.String("IP", r.RemoteAddr),
			zap.String("user_agent", r.UserAgent()),
		)
		start := time.Now()
		next.ServeHTTP(w, r)
		s.logger.Info(
			"completed request",
			zap.String("method", r.Method),
			zap.String("url", r.RequestURI),
			zap.String("IP", r.RemoteAddr),
			zap.String("user_agent", r.UserAgent()),
			zap.Time("request_datetime", start),
			zap.Duration("duration", time.Since(start)),
		)
	})
}

func (s *Server) Start(ctx context.Context) error {
	// TODO
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return nil
}

// TODO
