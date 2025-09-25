package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func openDB() (*sql.DB, error) {

	dsn, err := getRequiredStringEnv("DB_DSN")
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", *dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	maxOpenConns := getOptionalIntEnv("DB_MAX_OPEN_CONNS", 25)
	maxIdleConns := getOptionalIntEnv("DB_MAX_IDLE_CONNS", 10)
	maxIdleTime := getOptionalIntEnv("DB_MAX_LIFE_TIME", 20)
	maxLifeTime := getOptionalIntEnv("DB_MAX_IDLE_TIME", 5)

	db.SetConnMaxIdleTime(time.Duration(maxIdleTime) * time.Second)
	db.SetConnMaxLifetime(time.Duration(maxLifeTime) * time.Minute)

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create database driver: %w", err)
	}

	// Create file source instance
	source, err := (&file.File{}).Open("./migrations")
	if err != nil {
		return fmt.Errorf("could not open migrations directory: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithInstance("file", source, "postgres", driver)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run migrations: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}
