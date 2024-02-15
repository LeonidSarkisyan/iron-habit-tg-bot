package postgres

import (
	"HabitsBot/internal/models"
	"context"
	"github.com/jackc/pgx/v5"
)

type TimestampStorage struct {
	db *pgx.Conn
}

func NewTimestamp(db *pgx.Conn) *TimestampStorage {
	return &TimestampStorage{db}
}

func (t *TimestampStorage) Create(ts models.Timestamp) error {
	stmt := `INSERT INTO timestamps (day, time, habit_id) VALUES ($1, $2, $3)`

	_, err := t.db.Exec(
		context.Background(),
		stmt,
		ts.Day,
		ts.Time,
		ts.HabitID,
	)

	return err
}
