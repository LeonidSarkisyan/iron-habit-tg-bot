package postgres

import (
	"HabitsBot/internal/models"
	"context"
	"github.com/jackc/pgx/v5"
)

type ExecutionStorage struct {
	conn *pgx.Conn
}

func NewExecutionStorage(conn *pgx.Conn) *ExecutionStorage {
	return &ExecutionStorage{conn: conn}
}

func (s *ExecutionStorage) Create(execution models.Execution) error {
	stmt := "INSERT INTO executions (habit_id) VALUES ($1)"

	_, err := s.conn.Exec(context.Background(), stmt, execution.HabitID)

	if err != nil {
		return err
	}

	return nil
}
