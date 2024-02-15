package http

import (
	"errors"
	"fmt"
	"github.com/swmh/wbl2/develop/dev11/internal/service"
	"log/slog"
	"net/http"
	"strconv"
)

type DeleteEventRequest struct {
	Id int
}

func ParseDeleteEvent(r *http.Request) (DeleteEventRequest, error) {
	err := r.ParseForm()
	if err != nil {
		return DeleteEventRequest{}, err
	}

	var c DeleteEventRequest

	c.Id, err = strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		return DeleteEventRequest{}, fmt.Errorf("id is invalid: %w", err)
	}

	return c, nil
}

func (s *Server) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	req, err := ParseDeleteEvent(r)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = s.s.DeleteEvent(r.Context(), req.Id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			WriteError(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		s.l.Error("Cannot create req", slog.String("err", err.Error()))
		return
	}

	WriteResult(w, "ok")
}
