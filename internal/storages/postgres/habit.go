package postgres

import (
	"HabitsBot/internal/models"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type HabitStorage struct {
	db *pgx.Conn
}

func New(conn *pgx.Conn) *HabitStorage {
	return &HabitStorage{conn}
}

func (h *HabitStorage) Create(habit models.Habit) (int, error) {

	tx, err := h.db.Begin(context.Background())

	if err != nil {
		return 0, err
	}

	defer tx.Rollback(context.Background())

	stmt := `INSERT INTO habits (title, time_warning, time_completed, user_id) VALUES ($1, $2, $3, $4) RETURNING id;`

	var habitID int

	err = h.db.QueryRow(
		context.Background(), stmt, habit.Title, habit.WarningTime, habit.CompletedTime, habit.UserID,
	).Scan(&habitID)

	if err != nil {
		return 0, err
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
		return 0, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return 0, err
	}

	return habitID, nil
}

func (h *HabitStorage) Name(habitID int, userID int64) (string, error) {
	query := `SELECT h.title FROM habits h WHERE h.id = $1 AND h.user_id = $2`

	var name string

	err := h.db.QueryRow(context.Background(), query, habitID, userID).Scan(&name)

	return name, err
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

func (h *HabitStorage) GetAll(userID int64, offset int) ([]models.Habit, error) {
	var habits []models.Habit

	query := `
	SELECT h.id, h.title, ts.day, ts.time
	FROM habits h 
	LEFT JOIN timestamps ts ON ts.habit_id = h.id
	WHERE h.user_id = $1
	AND h.id IN (
		SELECT id
		FROM habits
		WHERE user_id = $1
		ORDER BY id ASC
		LIMIT 5 OFFSET $2
	)
	ORDER BY h.id ASC;
    `

	habitsMap := make(map[int]*models.Habit)

	rows, err := h.db.Query(context.Background(), query, userID, offset)
	if err != nil {
		log.Error().Err(err).Msg("Error executing query")
		return habits, err
	}
	defer rows.Close()

	for rows.Next() {
		var h models.Habit
		var ts models.Timestamp
		if err := rows.Scan(&h.ID, &h.Title, &ts.Day, &ts.Time); err != nil {
			log.Error().Err(err).Msg("Error scanning row")
			continue
		}

		if _, ok := habitsMap[h.ID]; !ok {
			habitsMap[h.ID] = &models.Habit{
				ID:         h.ID,
				Title:      h.Title,
				Timestamps: make([]models.Timestamp, 0),
			}
		}

		habitsMap[h.ID].Timestamps = append(habitsMap[h.ID].Timestamps, ts)
	}

	for _, habit := range habitsMap {
		habits = append(habits, *habit)
	}

	return habits, nil
}

func (h *HabitStorage) Habits() ([]models.Habit, error) {
	var habits []models.Habit

	query := `
	SELECT h.id, h.title, h.time_warning, h.time_completed, h.user_id, ts.day, ts.time
	FROM habits h 
	LEFT JOIN timestamps ts ON ts.habit_id = h.id
	AND h.id IN (
		SELECT id
		FROM habits
		ORDER BY id ASC
	)
	ORDER BY h.id ASC;
    `

	habitsMap := make(map[int]*models.Habit)

	rows, err := h.db.Query(context.Background(), query)
	if err != nil {
		log.Error().Err(err).Msg("Error executing query")
		return habits, err
	}
	defer rows.Close()

	for rows.Next() {
		var h models.Habit
		var ts models.Timestamp
		if err := rows.Scan(&h.ID, &h.Title, &h.WarningTime, &h.CompletedTime, &h.UserID, &ts.Day, &ts.Time); err != nil {
			log.Error().Err(err).Msg("Error scanning row")
			continue
		}

		if _, ok := habitsMap[h.ID]; !ok {
			habitsMap[h.ID] = &models.Habit{
				ID:            h.ID,
				Title:         h.Title,
				WarningTime:   h.WarningTime,
				CompletedTime: h.CompletedTime,
				UserID:        h.UserID,
				Timestamps:    make([]models.Timestamp, 0),
			}
		}

		habitsMap[h.ID].Timestamps = append(habitsMap[h.ID].Timestamps, ts)
	}

	for _, habit := range habitsMap {
		habits = append(habits, *habit)
	}

	return habits, nil
}
