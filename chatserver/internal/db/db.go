package db

import (
	"chatserver/internal/messages"
	"chatserver/internal/models"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
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

func FetchMessageCountByChannel(db *pgxpool.Pool) ([]messages.ChannelMessageCount, error) {
	rows, err := db.Query(context.Background(), `
		SELECT channel, COUNT(*) AS message_count
		FROM chatserver.chat_messages
		GROUP BY channel
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch message counts: %w", err)
	}
	defer rows.Close()

	var counts []messages.ChannelMessageCount
	for rows.Next() {
		var count messages.ChannelMessageCount
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

func FetchMessages(db *pgxpool.Pool, userID string, channels []string, keyword string, limit, offset int) ([]messages.ChatMessagePayload, bool, error) {
	var query string
	var args []interface{}
	var conditions []string

	query = `
		SELECT 
			m.id, 
			m.owner_id, 
			COALESCE(u.username, '[Unknown]') AS username, 
			m.channel, 
			m.message, 
			m.authored_at
		FROM chatserver.chat_messages m
		LEFT JOIN keycloak.public.user_entity u ON m.owner_id::TEXT = u.id
	`

	argIndex := 1

	// Add user filter if provided
	if userID != "" {
		conditions = append(conditions, fmt.Sprintf("m.owner_id = $%d", argIndex))
		args = append(args, userID)
		argIndex++
	}

	// Add channel filter if provided
	if len(channels) > 0 {
		placeholders := []string{}
		for _, channel := range channels {
			placeholders = append(placeholders, fmt.Sprintf("$%d", argIndex))
			args = append(args, channel)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("m.channel IN (%s)", strings.Join(placeholders, ", ")))
	}

	// Add keyword search if provided
	if keyword != "" {
		conditions = append(conditions, fmt.Sprintf("m.search_vector @@ plainto_tsquery('english', $%d)", argIndex))
		args = append(args, keyword)
		argIndex++
	}

	// Combine WHERE conditions
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Add ORDER BY, LIMIT, OFFSET
	query += fmt.Sprintf(`
		ORDER BY m.authored_at DESC
		LIMIT $%d OFFSET $%d;
	`, argIndex, argIndex+1)

	args = append(args, limit+1, offset) // Request one extra row

	rows, err := db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, false, fmt.Errorf("failed to fetch chat messages: %w", err)
	}
	defer rows.Close()

	searchMessages := []messages.ChatMessagePayload{}
	for rows.Next() {
		var msg messages.ChatMessagePayload
		if err := rows.Scan(&msg.ID, &msg.OwnerID, &msg.Username, &msg.Channel, &msg.Message, &msg.Sent); err != nil {
			return nil, false, fmt.Errorf("failed to scan chat message row: %w", err)
		}
		searchMessages = append(searchMessages, msg)
	}

	// Determine if there are more messages beyond this 'page'
	hasMore := len(searchMessages) > limit
	if hasMore {
		searchMessages = searchMessages[:limit]
	}

	log.Printf("Fetched %d messages from database (user: %s, channels: %v, keyword: %s)", len(searchMessages), userID, channels, keyword)
	return searchMessages, hasMore, nil
}

func RemoveMessage(db *pgxpool.Pool, messageID int) (bool, error) {
	query := `DELETE FROM chatserver.chat_messages WHERE id = $1`

	cmd, err := db.Exec(context.Background(), query, messageID)
	if err != nil {
		return false, fmt.Errorf("failed to delete message: %w", err)
	}

	// Check if a row was actually deleted
	if cmd.RowsAffected() == 0 {
		// No message found with that ID
		return false, nil
	}

	// Message was successfully deleted
	return true, nil
}

func RemoveMessages(db *pgxpool.Pool, messageIDs []int) (int64, error) {
	if len(messageIDs) == 0 {
		return 0, fmt.Errorf("no message IDs provided")
	}

	query := `DELETE FROM chatserver.chat_messages WHERE id = ANY($1)`
	cmd, err := db.Exec(context.Background(), query, messageIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to delete messages: %w", err)
	}

	return cmd.RowsAffected(), nil
}

func FetchUsers(db *pgxpool.Pool, username string) ([]messages.UserSearchResult, error) {
	ctx := context.Background()

	query := `
		SELECT DISTINCT ON (u.id) 
			u.id, 
			u.username,
			CASE 
				WHEN b.banished_id IS NOT NULL THEN TRUE 
				ELSE FALSE 
			END AS banned
		FROM keycloak.public.user_entity u
		LEFT JOIN chatserver.bans b 
			ON u.id = b.banished_id 
			AND (b.end_time IS NULL OR b.end_time > NOW())
	`

	var args []interface{}
	if username != "" {
		query += " WHERE u.username = $1"
		args = append(args, username)
	}

	// Ensure ORDER BY always comes after WHERE
	query += " ORDER BY u.id, b.start_time DESC"

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	defer rows.Close()

	var users []messages.UserSearchResult
	for rows.Next() {
		var user messages.UserSearchResult
		if err := rows.Scan(&user.ID, &user.Username, &user.Banned); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over user rows: %w", err)
	}

	return users, nil
}

func BanUser(db *pgxpool.Pool, ownerID, banishedID, reason string, duration int) error {
	ctx := context.Background()

	// Determine endTime based on duration
	var endTime *time.Time // Use a pointer to allow NULL values
	if duration > 0 {
		tempEndTime := time.Now().Add(time.Duration(duration) * time.Hour)
		endTime = &tempEndTime
	}

	// Convert empty reason to NULL
	var reasonSQL interface{}
	if reason == "" {
		reasonSQL = nil
	} else {
		reasonSQL = reason
	}

	query := `
		INSERT INTO chatserver.bans (owner_id, banished_id, reason, end_time)
		VALUES ($1, $2, $3, $4)
	`

	_, err := db.Exec(ctx, query, ownerID, banishedID, reasonSQL, endTime)
	if err != nil {
		log.Printf("Failed to insert ban record: %v", err)
		return err
	}

	log.Printf("User %s banned by %s for %d hours. Reason: %s", banishedID, ownerID, duration, reason)
	return nil
}

func IsUserBanned(db *pgxpool.Pool, banishedID string) (bool, error) {
	ctx := context.Background()

	// Check if the user is banned or pardoned
	query := `
		SELECT COUNT(*)
		FROM chatserver.bans
		WHERE banished_id = $1
		AND (end_time is NULL OR end_time > NOW())
		AND (pardoned IS NULL OR pardoned = FALSE)
	`

	var count int
	err := db.QueryRow(ctx, query, banishedID).Scan(&count)
	if err != nil {
		log.Printf("Failed to check ban status for user %s: %v", banishedID, err)
		return false, err
	}

	return count > 0, nil
}

func FetchBanRecords(db *pgxpool.Pool, limit, offset int) ([]models.BanRecord, bool, error) {
	ctx := context.Background()

	query := `
		SELECT 
			b.id, 
			b.owner_id, 
			b.banished_id, 
			COALESCE(u.username, '[Unknown]') AS banished_username, 
			b.start_time, 
			b.reason, 
			b.end_time, 
			b.duration::TEXT, 
			b.pardoned
		FROM chatserver.bans b
		LEFT JOIN keycloak.public.user_entity u 
			ON b.banished_id = u.id
		ORDER BY b.start_time DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := db.Query(ctx, query, limit+1, offset) // Fetch one extra row to check for more results
	if err != nil {
		return nil, false, fmt.Errorf("failed to fetch ban records: %w", err)
	}
	defer rows.Close()

	var bans []models.BanRecord
	for rows.Next() {
		var ban models.BanRecord
		var reason sql.NullString
		var duration sql.NullString
		var banishedUsername sql.NullString

		err := rows.Scan(
			&ban.ID,
			&ban.OwnerID,
			&ban.BanishedID,
			&banishedUsername,
			&ban.Start,
			&reason,
			&ban.End,
			&duration,
			&ban.Pardoned,
		)
		if err != nil {
			return nil, false, fmt.Errorf("failed to scan ban record: %w", err)
		}

		if reason.Valid {
			ban.Reason = &reason.String
		} else {
			ban.Reason = nil
		}

		if duration.Valid {
			ban.Duration = &duration.String
		} else {
			ban.Duration = nil
		}

		if banishedUsername.Valid {
			ban.BanishedUsername = banishedUsername.String
		} else {
			ban.BanishedUsername = "[Unknown]"
		}

		bans = append(bans, ban)
	}

	// Check if there are more results beyond this page
	hasMore := len(bans) > limit
	if hasMore {
		bans = bans[:limit] // Remove the extra record used for checking
	}

	log.Printf("Fetched %d ban records (limit: %d, offset: %d)", len(bans), limit, offset)
	return bans, hasMore, nil
}
