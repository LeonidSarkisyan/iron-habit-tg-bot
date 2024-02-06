package postgres

import (
	"HabitsBot/internal/models"
	"context"
	"github.com/jackc/pgx/v5"
)

type RejectionStorage struct {
	db *pgx.Conn
}

func NewRejection(conn *pgx.Conn) *RejectionStorage {
	return &RejectionStorage{conn}
}

func (h *RejectionStorage) Create(rejection models.Rejection) error {
	stmt := "INSERT INTO rejections (text, habit_id) VALUES ($1, $2)"

	_, err := h.db.Exec(context.Background(), stmt, rejection.Text, rejection.HabitID)

	return err
}
