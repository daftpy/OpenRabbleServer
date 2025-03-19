package db

import (
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
