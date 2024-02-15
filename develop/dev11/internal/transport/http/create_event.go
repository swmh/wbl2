package http

import (
	"errors"
	"fmt"
	"github.com/swmh/wbl2/develop/dev11/internal/service"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type CreateEventRequest struct {
	UserId      int
	Description string
	Date        time.Time
}

func ParseCreateEvent(r *http.Request) (CreateEventRequest, error) {
	err := r.ParseForm()
	if err != nil {
		return CreateEventRequest{}, err
	}

	var c CreateEventRequest

	c.UserId, err = strconv.Atoi(r.Form.Get("user_id"))
	if err != nil {
		return CreateEventRequest{}, fmt.Errorf("user_id is invalid: %w", err)
	}

	c.Description = r.Form.Get("description")
	if c.Description == "" {
		return CreateEventRequest{}, errors.New("description is empty")
	}

	c.Date, err = time.Parse("2006-01-02", r.Form.Get("date"))
	if err != nil {
		return CreateEventRequest{}, fmt.Errorf("date is invalid: %w", err)
	}

	return c, nil
}

func (s *Server) CreateEvent(w http.ResponseWriter, r *http.Request) {
	event, err := ParseCreateEvent(r)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	id, err := s.s.CreateEvent(r.Context(), service.Event{
		UserId:      event.UserId,
		Description: event.Description,
		Date:        event.Date,
	})

	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			WriteError(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		s.l.Error("Cannot create event", slog.String("err", err.Error()))
		return
	}

	WriteResult(w, strconv.Itoa(id))
}
