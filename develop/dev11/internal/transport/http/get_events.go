package http

import (
	"errors"
	"github.com/swmh/wbl2/develop/dev11/internal/service"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type GetEventsRequest struct {
	UserId int
	Date   time.Time
}

func ParseGetEvents(r *http.Request) (GetEventsRequest, error) {
	err := r.ParseForm()
	if err != nil {
		return GetEventsRequest{}, err
	}

	userId, err := strconv.Atoi(r.Form.Get("user_id"))
	if err != nil {
		return GetEventsRequest{}, err
	}

	date, err := time.Parse("2006-01-02", r.Form.Get("date"))
	if err != nil {
		return GetEventsRequest{}, err
	}

	return GetEventsRequest{
		UserId: userId,
		Date:   date,
	}, nil
}

func (s *Server) GetEventsForDay(w http.ResponseWriter, r *http.Request) {
	s.GetEvents(w, r, 24*time.Hour)
}

func (s *Server) GetEventsForWeek(w http.ResponseWriter, r *http.Request) {
	s.GetEvents(w, r, 24*7*time.Hour)
}

func (s *Server) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
	s.GetEvents(w, r, 24*31*time.Hour)
}
func (s *Server) GetEvents(w http.ResponseWriter, r *http.Request, d time.Duration) {
	req, err := ParseGetEvents(r)
	if err != nil {
		s.l.Error("Cannot parse GetEvents request", slog.String("err", err.Error()))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	events, err := s.s.GetEvents(r.Context(), req.UserId, service.TimeRange{
		From: req.Date,
		To:   req.Date.Add(d),
	})
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) || errors.Is(err, service.ErrInvalidTimeRange) {
			WriteError(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}

		s.l.Error("Cannot get events", slog.String("err", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	result := make([]Event, len(events))
	for i, v := range events {
		result[i] = EventToJson(v)
	}

	WriteResult(w, result)
}
