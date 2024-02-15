package sqlite

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"github.com/swmh/wbl2/develop/dev11/internal/service"
	"time"
)

var schema = `
DROP TABLE IF EXISTS users;
CREATE TABLE users
(
    id INTEGER PRIMARY KEY AUTOINCREMENT
);

INSERT INTO users (id)
VALUES (1), (2), (3);

DROP TABLE IF EXISTS events;
CREATE TABLE events
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    date        TEXT,
    description TEXT,
    user_id REFERENCES users (id) ON DELETE CASCADE
);
`

type Event struct {
	Id          int    `db:"id"`
	Date        string `db:"date"`
	Description string `db:"description"`
	UserId      int    `db:"user_id"`
}
type SqliteRepo struct {
	db *sqlx.DB
}

var ErrEventNotFound = errors.New("event not found")
var ErrUserNotFound = errors.New("user not found")

func New() (*SqliteRepo, error) {
	//db, err := sqlx.Connect("sqlite3", ":memory:")
	db, err := sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return &SqliteRepo{db: db}, err
}

func (s *SqliteRepo) CreateEvent(ctx context.Context, event service.Event) (id int, err error) {
	err = s.db.QueryRowxContext(ctx,
		`INSERT INTO events (date, description, user_id) VALUES (?, ?, ?) RETURNING id`,
		event.Date.Format("2006-01-02"), event.Description, event.UserId).Scan(&id)
	if err != nil {
		if errors.Is(err, sqlite3.ErrConstraintForeignKey) {
			return 0, ErrUserNotFound
		}

		return 0, err
	}

	return
}

func (s SqliteRepo) UpdateEvent(ctx context.Context, id int, description string) error {
	r, err := s.db.ExecContext(ctx, `UPDATE events SET description = ? WHERE id = ?`, description, id)
	if err != nil {
		return err
	}

	a, err := r.RowsAffected()
	if err != nil {
		return err
	}

	if a != 1 {
		return ErrEventNotFound
	}

	return nil
}

func (s SqliteRepo) DeleteEvent(ctx context.Context, id int) error {
	r, err := s.db.ExecContext(ctx, `DELETE FROM events WHERE id = ?`, id)
	if err != nil {
		return err
	}

	a, err := r.RowsAffected()
	if err != nil {
		return err
	}

	if a != 1 {
		return ErrEventNotFound
	}

	return err
}

func (s SqliteRepo) GetEvents(ctx context.Context, userId int, t service.TimeRange) ([]service.Event, error) {
	var events []Event
	err := s.db.SelectContext(ctx, &events,
		`SELECT id, date, description, user_id FROM events WHERE user_id = ? AND date BETWEEN ? AND ?`,
		userId, t.From.Format("2006-01-02"), t.To.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}

	result := make([]service.Event, len(events))

	for i, v := range events {
		t, err := time.Parse("2006-01-02", v.Date)
		if err != nil {
			return nil, err
		}

		result[i] = service.Event{
			Id:          v.Id,
			UserId:      v.UserId,
			Description: v.Description,
			Date:        t,
		}
	}

	return result, nil
}

func (s SqliteRepo) UserNotFound(err error) bool {
	return errors.Is(err, ErrUserNotFound)
}

func (s SqliteRepo) EventNotFound(err error) bool {
	return errors.Is(err, ErrEventNotFound)
}
