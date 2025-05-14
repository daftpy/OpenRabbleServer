package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"onrabble.com/chatserver/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RecordUserSession(db *pgxpool.Pool, userID string, start, end time.Time) error {
	query := `INSERT INTO chatserver.chat_sessions (owner_id, start_time, end_time) VALUES ($1, $2, $3)`
	_, err := db.Exec(context.Background(), query, userID, start, end)
	return err
}

func FetchSessionActivity(db *pgxpool.Pool, userID string) ([]models.SessionActivity, error) {
	query := `
		SELECT
			DATE(start_time) AS session_date,
			COUNT(id) AS session_count,
			SUM(end_time - start_time)::TEXT AS total_duration
		FROM chatserver.chat_sessions
		WHERE start_time >= NOW() - INTERVAL '7 days'
	`

	var args []interface{}

	if userID != "" {
		query += " AND owner_id = $1"
		args = append(args, userID)
	}

	query += `
		GROUP BY session_date
		ORDER BY session_date;
	`

	rows, err := db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch session activity: %w", err)
	}
	defer rows.Close()

	var activity []models.SessionActivity
	for rows.Next() {
		var sa models.SessionActivity
		var sessionDate time.Time // Use time.Time to scan DATE properly
		var totalDuration string  // Keep as string since SUM(interval) returns TEXT

		err := rows.Scan(&sessionDate, &sa.SessionCount, &totalDuration)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session activity row: %w", err)
		}

		// Convert sessionDate to a string in YYYY-MM-DD format
		sa.SessionDate = sessionDate.Format("2006-01-02")
		sa.TotalDuration = totalDuration

		activity = append(activity, sa)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating session activity rows: %w", err)
	}

	log.Printf("Loaded %d session activity records", len(activity))
	return activity, nil
}
