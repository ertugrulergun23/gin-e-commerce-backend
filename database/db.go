package database

import (
	"database/sql"
	"embed"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

func InitDB() (*sql.DB, *gorm.DB, error) {
	dsn := "host=postgres port=5432 user=myuser password=mypassword dbname=mydatabase sslmode=disable"
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, nil, err
	}

	if err := runMigrations(db); err != nil {
		return nil, nil, fmt.Errorf("migrasyon hatasi: %w", err)
	}

	gormDB, err := gorm.Open(gormpostgres.New(gormpostgres.Config{
		Conn: db,
	}), &gorm.Config{})

	return db, gormDB, nil
}

func runMigrations(db *sql.DB) error {
	// Gömülü dosya sisteminden okuyoruz
	sourceDriver, err := iofs.New(migrationFS, "migrations")
	if err != nil {
		return fmt.Errorf("migrasyon dosyalari okunamadi: %w", err)
	}

	// PostgreSQL sürücüsünü hazırlıyoruz
	dbDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("postgres surucu hatasi: %w", err)
	}

	// Migrate motorunu başlatıyoruz
	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", dbDriver)
	if err != nil {
		return fmt.Errorf("migrate ornegi olusturulamadi: %w", err)
	}

	// Eksik olan migrasyonları uyguluyoruz
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	if err == migrate.ErrNoChange {
		log.Println("[DB] Semada değişiklik yok, veritabanı güncel.")
	} else {
		log.Println("[DB] Başarılı! Tablolar ve migrasyonlar uygulandı.")
	}

	return nil
}
