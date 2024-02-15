package http

import (
	"encoding/json"
	"github.com/swmh/wbl2/develop/dev11/internal/service"
	"net/http"
	"time"
)

type Event struct {
	Id          int       `json:"id"`
	UserId      int       `json:"userId"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

func EventToJson(event service.Event) Event {
	return Event{
		Id:          event.Id,
		UserId:      event.UserId,
		Description: event.Description,
		Date:        event.Date,
	}
}

func WriteResponse(w http.ResponseWriter, v any, code int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

type Response struct {
	Result any `json:"result"`
}

func WriteResult(w http.ResponseWriter, result any) error {
	return WriteResponse(w, Response{Result: result}, http.StatusOK)
}

type Error struct {
	Error string `json:"error"`
}

func WriteError(w http.ResponseWriter, err string, code int) error {
	return WriteResponse(w, Error{Error: err}, code)
}
