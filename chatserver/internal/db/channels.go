package db

import (
	"chatserver/internal/messages/api"
	"chatserver/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateChannel(db *pgxpool.Pool, name string, description string) error {
	placeholderOwner := "00000000-0000-0000-0000-000000000000"

	// Get current max sort_order
	var maxOrder int
	err := db.QueryRow(context.Background(),
		`SELECT COALESCE(MAX(sort_order), 0) FROM chatserver.channels`,
	).Scan(&maxOrder)
	if err != nil {
		return fmt.Errorf("failed to get max sort_order: %w", err)
	}

	_, err = db.Exec(context.Background(), `
		INSERT INTO chatserver.channels (name, description, owner_id, sort_order)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (name) DO NOTHING
	`, name, description, placeholderOwner, maxOrder+1)

	return err
}

// FetchChannels retrieves all channels from the database's `channels` table.
// For each channel, if the `description` column is NULL, the Channel.Description
// field will be set to nil in the returned slice. Otherwise, it contains the
// description as a pointer to a string.
// It returns:
//  1. A slice of Channel models
//  2. An error, if any
func FetchChannels(db *pgxpool.Pool) ([]models.Channel, error) {
	rows, err := db.Query(context.Background(), "SELECT id, name, description FROM chatserver.channels ORDER BY sort_order, id")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch channels: %w", err)
	}
	defer rows.Close()

	var channels []models.Channel
	for rows.Next() {
		var channel models.Channel
		var description sql.NullString

		if err := rows.Scan(&channel.ID, &channel.Name, &description); err != nil {
			return nil, fmt.Errorf("failed to scan channel row: %w", err)
		}

		// Convert sql.NullString to *string
		if description.Valid {
			channel.Description = &description.String
		} else {
			channel.Description = nil // Keep as nil if it was NULL in DB
		}

		channels = append(channels, channel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	log.Printf("Loaded %d channels from database", len(channels))
	return channels, nil
}

func UpdateChannel(db *pgxpool.Pool, ID int, name *string, description *string) error {
	clauses := []string{}
	params := []interface{}{}
	paramIndex := 1

	if name != nil {
		clauses = append(clauses, "name = $"+strconv.Itoa(paramIndex))
		params = append(params, *name)
		paramIndex++
	}

	if description != nil {
		clauses = append(clauses, "description = $"+strconv.Itoa(paramIndex))
		params = append(params, *description)
		paramIndex++
	}

	if len(clauses) == 0 {
		return errors.New("no fields provided to update")
	}

	query := `
        UPDATE chatserver.channels
        SET ` + strings.Join(clauses, ", ") + ` 
        WHERE id = $` + strconv.Itoa(paramIndex)

	params = append(params, ID)

	_, err := db.Exec(context.Background(), query, params...)
	return err
}

// FetchMessageCountByChannel returns the total number of chat messages per channel.
// It queries the `chat_messages` table, grouping by the `channel` column to produce
// a list of channels and their respective message counts.
// It returns:
//  1. A slice of ChannelMessageCount, where each item contains a channel name and a message count
//  2. An error, if the query or row scanning fails
func FetchMessageCountByChannel(db *pgxpool.Pool) ([]api.ChannelMessageCount, error) {
	rows, err := db.Query(context.Background(), `
		SELECT channel, COUNT(*) AS message_count
		FROM chatserver.chat_messages
		GROUP BY channel
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch message counts: %w", err)
	}
	defer rows.Close()

	var counts []api.ChannelMessageCount
	for rows.Next() {
		var count api.ChannelMessageCount
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

// MoveChannelBefore reorders a channel to appear before another channel.
// If beforeID is nil, it moves the channel to the end.
func MoveChannelBefore(db *pgxpool.Pool, movedID int, beforeID *int) error {
	ctx := context.Background()

	if beforeID != nil && *beforeID == movedID {
		return nil
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// Step 1: Fetch current ordered list of channel IDs
	rows, err := tx.Query(ctx, `SELECT id FROM chatserver.channels ORDER BY sort_order, id`)
	if err != nil {
		return fmt.Errorf("failed to fetch channel order: %w", err)
	}
	defer rows.Close()

	var ordered []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("failed to scan channel ID: %w", err)
		}
		ordered = append(ordered, id)
	}

	// Step 2: Remove movedID from the list
	newOrder := make([]int, 0, len(ordered))
	for _, id := range ordered {
		if id != movedID {
			newOrder = append(newOrder, id)
		}
	}

	// Step 3: Insert movedID before target or at end
	inserted := false
	if beforeID != nil {
		for i, id := range newOrder {
			if id == *beforeID {
				// insert movedID before this index
				newOrder = append(newOrder[:i], append([]int{movedID}, newOrder[i:]...)...)
				inserted = true
				break
			}
		}
	}
	if !inserted {
		newOrder = append(newOrder, movedID)
	}

	// Step 4: Reassign sort_order based on new list
	for index, id := range newOrder {
		_, err := tx.Exec(ctx, `
			UPDATE chatserver.channels
			SET sort_order = $1
			WHERE id = $2
		`, index+1, id)
		if err != nil {
			return fmt.Errorf("failed to update sort_order for id %d: %w", id, err)
		}
	}

	return tx.Commit(ctx)
}

func RemoveChannelByID(db *pgxpool.Pool, channelID int, purgeMessages bool) error {
	ctx := context.Background()

	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// Step 1: Get the channel name (needed for message deletion)
	var channelName string
	err = tx.QueryRow(ctx, `
		SELECT name FROM chatserver.channels WHERE id = $1
	`, channelID).Scan(&channelName)
	if err != nil {
		return fmt.Errorf("failed to fetch channel name: %w", err)
	}

	// Step 2: Delete the channel row
	_, err = tx.Exec(ctx, `
		DELETE FROM chatserver.channels
		WHERE id = $1
	`, channelID)
	if err != nil {
		return fmt.Errorf("failed to delete channel: %w", err)
	}

	// Step 3: Optionally purge messages
	if purgeMessages {
		_, err := removeMessagesByChannelTx(tx, channelName)
		if err != nil {
			return fmt.Errorf("failed to purge messages: %w", err)
		}
	}

	// Step 4: Renumber sort_order
	rows, err := tx.Query(ctx, `
		SELECT id FROM chatserver.channels ORDER BY sort_order, id
	`)
	if err != nil {
		return fmt.Errorf("failed to fetch remaining channels: %w", err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("failed to scan remaining channel ID: %w", err)
		}
		ids = append(ids, id)
	}

	for i, id := range ids {
		_, err := tx.Exec(ctx, `
			UPDATE chatserver.channels
			SET sort_order = $1
			WHERE id = $2
		`, i+1, id)
		if err != nil {
			return fmt.Errorf("failed to renumber sort_order: %w", err)
		}
	}

	return tx.Commit(ctx)
}

// Called inside RemoveChannelByID, already in a tx
func removeMessagesByChannelTx(tx pgx.Tx, channelName string) (int64, error) {
	cmd, err := tx.Exec(context.Background(), `
		DELETE FROM chatserver.chat_messages
		WHERE channel = $1
	`, channelName)
	if err != nil {
		return 0, fmt.Errorf("failed to delete messages for channel '%s': %w", channelName, err)
	}
	return cmd.RowsAffected(), nil
}
