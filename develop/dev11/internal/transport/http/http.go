package http

import (
	"context"
	"github.com/swmh/wbl2/develop/dev11/internal/service"
	"log/slog"
	"net/http"
	"time"
)

type LoggingMiddleware struct {
	logger *slog.Logger
}

func (l *LoggingMiddleware) Handler(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		l.logger.Info(
			"",
			slog.String("addr", r.URL.String()),
			slog.Time("time", time.Now()),
		)
		h(w, r)
	}
}

type Handler func(w http.ResponseWriter, r *http.Request)

func Method(method string, h Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if method != r.Method {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		h(w, r)
	}
}

type Config struct {
	Addr    string
	Service *service.Service
	Logger  *slog.Logger
}

type Server struct {
	server *http.Server
	s      *service.Service
	l      *slog.Logger
}

func New(c Config) Server {
	mux := http.NewServeMux()
	serv := &http.Server{
		Addr:    c.Addr,
		Handler: mux,
	}

	s := Server{
		server: serv,
		s:      c.Service,
		l:      c.Logger,
	}

	l := LoggingMiddleware{logger: c.Logger}

	mux.HandleFunc("/create_event", l.Handler(Method("POST", s.CreateEvent)))
	mux.HandleFunc("/update_event", l.Handler(Method("POST", s.UpdateEvent)))
	mux.HandleFunc("/delete_event", l.Handler(Method("POST", s.DeleteEvent)))
	mux.HandleFunc("/events_for_day", l.Handler(Method("GET", s.GetEventsForDay)))
	mux.HandleFunc("/events_for_week", l.Handler(Method("GET", s.GetEventsForWeek)))
	mux.HandleFunc("/events_for_month", l.Handler(Method("GET", s.GetEventsForMonth)))

	return s
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
