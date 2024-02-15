package postgres

import (
	"HabitsBot/internal/models"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
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

func (t *TimestampStorage) GetByHabitID(habitID int, offset int) ([]models.Timestamp, bool, error) {
	var timestamps []models.Timestamp

	stmt := `
	SELECT id, day, time,
		EXISTS (
		   SELECT 1
		   FROM timestamps ts
		   WHERE ts.habit_id = $1
		   OFFSET $2 + 5
		   LIMIT 1
	   	) AS more_records_exist
	FROM timestamps
	WHERE habit_id = $1
	ORDER BY id ASC
	LIMIT 5 OFFSET $2;
	`

	rows, err := t.db.Query(context.Background(), stmt, habitID, offset)

	if err != nil {
		log.Error().Err(err).Send()
		return nil, false, err
	}

	var existsMore bool

	for rows.Next() {
		var timestamp models.Timestamp
		err = rows.Scan(&timestamp.ID, &timestamp.Day, &timestamp.Time, &existsMore)
		if err != nil {
			log.Error().Err(err).Send()
			return nil, false, err
		}
		timestamps = append(timestamps, timestamp)
	}

	return timestamps, existsMore, nil
}
