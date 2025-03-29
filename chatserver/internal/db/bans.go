package db

import (
	"chatserver/internal/models"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

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

func PardonUser(db *pgxpool.Pool, banishedID int) error {
	query := `
		UPDATE chatserver.bans
		SET pardoned = TRUE
		WHERE id = $1;
	`
	_, err := db.Exec(context.Background(), query, banishedID)
	return err
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
