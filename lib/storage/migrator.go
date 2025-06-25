package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type Migrator interface {
	Close() error
	Up() error
	Down() error
	StorageVersion() (int64, error)
}

type migratorImpl struct {
	DbURL         string
	MigrationsDir string
	Db            *sql.DB
}

func (migrator *migratorImpl) Close() error {
	if migrator.Db != nil {
		if err := migrator.Db.Close(); err != nil {
			return fmt.Errorf("failed to close database connection: %w", err)
		}
		log.Println("Database connection closed successfully!")
	}

	return nil
}

func (migrator *migratorImpl) Up() error {
	if migrator.Db.Ping() != nil {
		return fmt.Errorf("the DB connection seems to be dead")
	}

	if err := goose.Up(migrator.Db, migrator.MigrationsDir); err != nil {
		return fmt.Errorf("goose up failed: %w", err)
	}

	log.Println("Migrations applied successfully!")
	return nil
}

func (migrator *migratorImpl) Down() error {
	if migrator.Db.Ping() != nil {
		return fmt.Errorf("the DB connection seems to be dead")
	}

	if err := goose.Down(migrator.Db, migrator.MigrationsDir); err != nil {
		return fmt.Errorf("goose down failed: %w", err)
	}

	log.Println("Migration rolled back successfully!")
	return nil
}

func (migrator *migratorImpl) StorageVersion() (int64, error) {
	if migrator.Db.Ping() != nil {
		return -1, fmt.Errorf("the DB connection seems to be dead")
	}

	version, err := goose.GetDBVersion(migrator.Db)
	if err != nil {
		return -1, fmt.Errorf("unable totobtain DB version: %w", err)
	}

	return version, nil
}

func NewMigrator(dbURL, migrationsDir string) (Migrator, error) {
	if dbURL == "" {
		return nil, fmt.Errorf("database URL cannot be empty")
	}
	if migrationsDir == "" {
		return nil, fmt.Errorf("migration directory cannot be emlpty")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	var migrator = &migratorImpl{
		DbURL:         dbURL,
		MigrationsDir: migrationsDir,
		Db:            db,
	}
	return migrator, nil
}
