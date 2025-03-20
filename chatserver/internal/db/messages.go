package db

import (
	"chatserver/internal/models"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func FetchMessages(db *pgxpool.Pool, userID string, channels []string, keyword string, limit, offset int) ([]models.ChatMessage, bool, error) {
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

	searchMessages := []models.ChatMessage{}
	for rows.Next() {
		var msg models.ChatMessage
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

func RemoveMessages(db *pgxpool.Pool, messageIDs []int) (int64, []int, error) {
	if len(messageIDs) == 0 {
		return 0, nil, fmt.Errorf("no message IDs provided")
	}

	ctx := context.Background()

	// Step 1: Fetch cacheIDs before deletion
	queryCacheIDs := `SELECT cache_id FROM chatserver.chat_messages WHERE id = ANY($1)`
	rows, err := db.Query(ctx, queryCacheIDs, messageIDs)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to fetch cacheIDs: %w", err)
	}
	defer rows.Close()

	var cacheIDs []int
	for rows.Next() {
		var cacheID int
		if err := rows.Scan(&cacheID); err != nil {
			return 0, nil, fmt.Errorf("failed to scan cacheID: %w", err)
		}
		cacheIDs = append(cacheIDs, cacheID)
	}

	// Ensure we proceed only if there are valid cacheIDs
	if len(cacheIDs) == 0 {
		return 0, nil, fmt.Errorf("no matching messages found for provided IDs")
	}

	// Step 2: Delete messages from the database
	queryDelete := `DELETE FROM chatserver.chat_messages WHERE id = ANY($1)`
	cmd, err := db.Exec(ctx, queryDelete, messageIDs)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to delete messages: %w", err)
	}

	return cmd.RowsAffected(), cacheIDs, nil
}
