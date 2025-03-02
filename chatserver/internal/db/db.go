package db

import (
	"chatserver/internal/models"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect() (*pgxpool.Pool, error) {
	dsn := "postgres://keycloak:keycloak@postgres:5432/keycloak?sslmode=disable&search_path=chatserver"

	log.Println("Attempting to connect to PostgreSQL...")

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Println("Failed to parse database config:", err)
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	dbPool, err := pgxpool.New(context.Background(), config.ConnString())
	if err != nil {
		log.Println("Failed to create database connection pool:", err)
		return nil, fmt.Errorf("failed to create database connection pool: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL.")

	// Run schema migration from file
	if err := runMigrations(dbPool); err != nil {
		log.Println("Failed to apply database schema:", err)
		return nil, fmt.Errorf("failed to apply database schema: %w", err)
	}

	log.Println("Database schema applied successfully.")

	// Test query to verify connection works
	if err := testQuery(dbPool); err != nil {
		log.Println("Test query failed:", err)
		return nil, fmt.Errorf("test query failed: %w", err)
	}

	log.Println("Database test query successful.")

	DB = dbPool
	return dbPool, nil
}

// runMigrations loads and executes schema.sql
func runMigrations(db *pgxpool.Pool) error {
	sqlFile := "internal/db/schema.sql" // Adjust path if necessary

	// Check if the file exists
	info, err := os.Stat(sqlFile)
	if os.IsNotExist(err) {
		return fmt.Errorf("schema file does not exist at path: %s", sqlFile)
	} else if err != nil {
		return fmt.Errorf("error checking schema file: %w", err)
	}
	log.Printf("Found schema file: %s", info.Name())

	// Read the SQL file
	query, err := os.ReadFile(sqlFile)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// Execute the SQL file
	_, err = db.Exec(context.Background(), string(query))
	if err != nil {
		return fmt.Errorf("failed to execute schema.sql: %w", err)
	}
	log.Println("Database schema applied.")
	return nil
}

// testQuery runs a simple query to ensure the database is responding
func testQuery(db *pgxpool.Pool) error {
	var currentTime time.Time
	err := db.QueryRow(context.Background(), "SELECT NOW();").Scan(&currentTime)
	if err != nil {
		return fmt.Errorf("failed to execute test query: %w", err)
	}
	log.Printf("Database is responding. Current time: %s", currentTime.Format(time.RFC3339))
	return nil
}

func FetchChannels(db *pgxpool.Pool) ([]models.Channel, error) {
	// Update query to select both name and description
	rows, err := db.Query(context.Background(), "SELECT name, description FROM chatserver.channels")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch channels: %w", err)
	}
	defer rows.Close()

	var channels []models.Channel
	for rows.Next() {
		var channel models.Channel
		if err := rows.Scan(&channel.Name, &channel.Description); err != nil {
			return nil, fmt.Errorf("failed to scan channel row: %w", err)
		}
		channels = append(channels, channel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	log.Printf("Loaded %d channels from database", len(channels))
	return channels, nil
}

type ChannelMessageCount struct {
	Channel      string `json:"channel"`
	MessageCount int    `json:"message_count"`
}

func FetchMessageCountByChannel(db *pgxpool.Pool) ([]ChannelMessageCount, error) {
	rows, err := db.Query(context.Background(), `
		SELECT channel, COUNT(*) AS message_count
		FROM chatserver.chat_messages
		GROUP BY channel
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch message counts: %w", err)
	}
	defer rows.Close()

	var counts []ChannelMessageCount
	for rows.Next() {
		var count ChannelMessageCount
		if err := rows.Scan(&count.Channel, &count.MessageCount); err != nil {
			return nil, fmt.Errorf("failed to scan message count row: %w", err)
		}
		counts = append(counts, count)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	log.Println("Loaded message counts")
	return counts, nil
}

func RecordUserSession(db *pgxpool.Pool, userID string, start, end time.Time) error {
	query := `INSERT INTO chatserver.chat_sessions (owner_id, start_time, end_time) VALUES ($1, $2, $3)`
	_, err := db.Exec(context.Background(), query, userID, start, end)
	return err
}

func FetchSessionActivity(db *pgxpool.Pool) ([]models.SessionActivity, error) {
	query := `
		SELECT
			DATE(start_time) AS session_date,
			COUNT(id) AS session_count,
			SUM(end_time - start_time)::TEXT AS total_duration
		FROM chatserver.chat_sessions
		WHERE start_time >= NOW() - INTERVAL '7 days'
		GROUP BY session_date
		ORDER BY session_date;
	`

	rows, err := db.Query(context.Background(), query)
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
