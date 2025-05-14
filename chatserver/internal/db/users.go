package db

import (
	"context"
	"fmt"

	"onrabble.com/chatserver/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func FetchUsers(db *pgxpool.Pool, username string) ([]models.User, error) {
	ctx := context.Background()

	query := `
		SELECT DISTINCT ON (u.id) 
			u.id, 
			u.username,
			CASE 
				WHEN b.id IS NOT NULL THEN TRUE 
				ELSE FALSE 
			END AS banned
		FROM keycloak.public.user_entity u
		LEFT JOIN chatserver.bans b 
			ON u.id = b.banished_id 
			AND b.pardoned = FALSE
			AND (
				b.end_time IS NULL OR b.end_time > NOW()
			)
	`

	var args []interface{}
	if username != "" {
		query += " WHERE u.username = $1"
		args = append(args, username)
	}

	query += " ORDER BY u.id, b.start_time DESC"

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
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
