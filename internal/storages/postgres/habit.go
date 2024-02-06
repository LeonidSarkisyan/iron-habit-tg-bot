package postgres

import (
	"HabitsBot/internal/models"
	"context"
	"github.com/jackc/pgx/v5"
)

type HabitStorage struct {
	db *pgx.Conn
}

func New(conn *pgx.Conn) *HabitStorage {
	return &HabitStorage{conn}
}

func (h *HabitStorage) Create(habit models.Habit) error {

	tx, err := h.db.Begin(context.Background())

	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	stmt := `INSERT INTO habits (title, time_warning, time_completed, user_id) VALUES ($1, $2, $3, $4) RETURNING id;`

	var habitID int

	err = h.db.QueryRow(
		context.Background(), stmt, habit.Title, habit.WarningTime, habit.CompletedTime, habit.UserID,
	).Scan(&habitID)

	if err != nil {
		return err
	}

	var rows [][]any
	for _, ts := range habit.Timestamps {
		rows = append(rows, []any{ts.Day, ts.Time, habitID})
	}

	_, err = h.db.CopyFrom(
		context.Background(),
		pgx.Identifier{"timestamps"},
		[]string{"day", "time", "habit_id"},
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (h *HabitStorage) Get(userID int64) (models.Habit, error) {
	query := `
	SELECT h.id, h.title, h.time_warning, h.time_completed, t.day, t.time FROM habits h 
	JOIN timestamps t ON t.habit_id = h.id
	WHERE h.user_id = $1
	`

	var habit models.Habit

	rows, err := h.db.Query(context.Background(), query, userID)

	defer rows.Close()

	if err != nil {
		return habit, err
	}

	for rows.Next() {
		var ts models.Timestamp
		err = rows.Scan(&habit.ID, &habit.Title, &habit.WarningTime, &habit.CompletedTime, &ts.Day, &ts.Time)
		if err != nil {
			return habit, err
		}
		habit.Timestamps = append(habit.Timestamps, ts)
	}

	return habit, nil
}
