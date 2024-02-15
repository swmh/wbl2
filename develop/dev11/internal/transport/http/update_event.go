package http

import (
	"errors"
	"fmt"
	"github.com/swmh/wbl2/develop/dev11/internal/service"
	"log/slog"
	"net/http"
	"strconv"
)

type UpdateEventRequest struct {
	Id          int
	Description string
}

func ParseUpdateEventRequest(r *http.Request) (UpdateEventRequest, error) {
	err := r.ParseForm()
	if err != nil {
		return UpdateEventRequest{}, err
	}

	var c UpdateEventRequest

	c.Id, err = strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		return UpdateEventRequest{}, fmt.Errorf("invalid id: %w", err)
	}

	c.Description = r.Form.Get("description")
	if c.Description == "" {
		return UpdateEventRequest{}, errors.New("empty description")

	}

	return c, nil
}

func (s *Server) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	req, err := ParseUpdateEventRequest(r)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = s.s.UpdateEvent(r.Context(), req.Id, req.Description)
	if err != nil {
		if errors.Is(err, service.ErrEventNotFound) {
			WriteError(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}

		s.l.Error("Cannot update event", slog.String("err", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	WriteResult(w, "ok")
}
