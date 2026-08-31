// Package database manages database connectivity, configuration, and migrations.
package database

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

// getEnv retrieves an environment variable or returns a fallback value if unset.
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// InitDB initializes PostgreSQL connection, runs migrations, and returns SQL and GORM DB handles.
func InitDB() (*sql.DB, *gorm.DB, error) {
	host := getEnv("DB_HOST", "postgres")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "myuser")
	password := getEnv("DB_PASSWORD", "mypassword")
	dbname := getEnv("DB_NAME", "mydatabase")
	sslmode := getEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, nil, err
	}

	if err := runMigrations(db); err != nil {
		return nil, nil, fmt.Errorf("migration error: %w", err)
	}

	gormDB, err := gorm.Open(gormpostgres.New(gormpostgres.Config{
		Conn: db,
	}), &gorm.Config{})

	return db, gormDB, nil
}

// runMigrations applies pending database schema migrations using embedded SQL files.
func runMigrations(db *sql.DB) error {
	// Read migration files from the embedded filesystem
	sourceDriver, err := iofs.New(migrationFS, "migrations")
	if err != nil {
		return fmt.Errorf("migrasyon dosyalari okunamadi: %w", err)
	}

	// Prepare the PostgreSQL driver instance
	dbDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("postgres surucu hatasi: %w", err)
	}

	// Initialize the migrate engine instance
	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", dbDriver)
	if err != nil {
		return fmt.Errorf("migrate ornegi olusturulamadi: %w", err)
	}

	// Apply pending migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	if err == migrate.ErrNoChange {
		log.Println("[DB] Schema has no changes, database is up to date.")
	} else {
		log.Println("[DB] Success! Tables and migrations applied.")
	}

	return nil
}
