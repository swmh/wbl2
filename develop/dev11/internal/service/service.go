package service

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Event struct {
	Id          int
	UserId      int
	Description string
	Date        time.Time
}

type TimeRange struct {
	From time.Time
	To   time.Time
}

func (t *TimeRange) Validate() bool {
	return !t.From.After(t.To)
}

type Repo interface {
	CreateEvent(ctx context.Context, event Event) (id int, err error)
	UpdateEvent(ctx context.Context, id int, description string) error
	DeleteEvent(ctx context.Context, id int) error
	GetEvents(ctx context.Context, userId int, t TimeRange) ([]Event, error)

	UserNotFound(err error) bool
	EventNotFound(err error) bool
}

type Service struct {
	repo Repo
}

func New(repo Repo) *Service {
	return &Service{
		repo: repo,
	}

}

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrEventNotFound    = errors.New("event not found")
	ErrInvalidTimeRange = errors.New("invalid time range")
)

type CreateEvent struct {
}

func (s *Service) CreateEvent(ctx context.Context, event Event) (int, error) {
	id, err := s.repo.CreateEvent(ctx, event)
	if err != nil {
		if s.repo.UserNotFound(err) {
			return 0, errors.Join(ErrUserNotFound, err)
		}

		return 0, fmt.Errorf("cannot create event in repo: %w", err)
	}

	return id, nil
}

func (s *Service) UpdateEvent(ctx context.Context, id int, description string) error {
	err := s.repo.UpdateEvent(ctx, id, description)
	if err != nil {
		if s.repo.EventNotFound(err) {
			return errors.Join(ErrEventNotFound, err)
		}
		return fmt.Errorf("cannot update event in repo: %w", err)
	}

	return nil
}

func (s *Service) DeleteEvent(ctx context.Context, id int) error {
	err := s.repo.DeleteEvent(ctx, id)
	if err != nil {
		if s.repo.EventNotFound(err) {
			return errors.Join(ErrEventNotFound, err)
		}

		return fmt.Errorf("cannot delete event in repo: %w", err)
	}

	return nil
}

func (s *Service) GetEvents(ctx context.Context, id int, t TimeRange) ([]Event, error) {
	if !t.Validate() {
		return nil, ErrInvalidTimeRange
	}

	events, err := s.repo.GetEvents(ctx, id, t)
	if err != nil {
		return nil, fmt.Errorf("cannot get events from repo: %w", err)
	}

	return events, nil
}
