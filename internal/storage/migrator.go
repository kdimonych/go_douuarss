package storage

import (
	"context"
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
	MigrationsDir string
	Db            *sql.DB
}

func (migrator *migratorImpl) Close() error {
	if migrator.Db != nil {
		if err := migrator.Db.Close(); err != nil {
			return fmt.Errorf("failed to close database connection: %w", err)
		}
		log.Println("[Info] Database connection closed successfully!")
	}

	return nil
}

func (migrator *migratorImpl) Up() error {
	if migrator.Db.PingContext(context.Background()) != nil {
		return fmt.Errorf("the DB connection seems to be dead")
	}

	if err := goose.Up(migrator.Db, migrator.MigrationsDir); err != nil {
		return fmt.Errorf("goose up failed: %w", err)
	}

	log.Println("[Info] Migrations applied successfully!")
	return nil
}

func (migrator *migratorImpl) Down() error {
	if migrator.Db.PingContext(context.Background()) != nil {
		return fmt.Errorf("the DB connection seems to be dead")
	}

	if err := goose.Down(migrator.Db, migrator.MigrationsDir); err != nil {
		return fmt.Errorf("goose down failed: %w", err)
	}

	log.Println("[Info] Migration rolled back successfully!")
	return nil
}

func (migrator *migratorImpl) StorageVersion() (int64, error) {
	if migrator.Db.PingContext(context.Background()) != nil {
		return -1, fmt.Errorf("the DB connection seems to be dead")
	}

	version, err := goose.GetDBVersion(migrator.Db)
	if err != nil {
		return -1, fmt.Errorf("unable totobtain DB version: %w", err)
	}

	return version, nil
}

type MigratorBuilder interface {
	WithDbConnectionFabric(dbConnectionFabric DbConnectionFabric) MigratorBuilder
	Build(ctx context.Context, dbURL, migrationsDir string) (Migrator, error)
}

type migratorBuilderImpl struct {
	dbConnectionFabric DbConnectionFabric
}

func NewMigratorBuilder() MigratorBuilder {
	return &migratorBuilderImpl{
		dbConnectionFabric: NewDbConnectionFabric(),
	}
}

func (builder *migratorBuilderImpl) WithDbConnectionFabric(dbConnectionFabric DbConnectionFabric) MigratorBuilder {
	if dbConnectionFabric == nil {
		log.Println("[Warning] Nil DbConnectionFabric provided. Using default DbConnectionFabric")
		dbConnectionFabric = NewDbConnectionFabric()
	}
	builder.dbConnectionFabric = dbConnectionFabric
	return builder
}

func (builder *migratorBuilderImpl) Build(ctx context.Context, dbURL, migrationsDir string) (Migrator, error) {
	if dbURL == "" {
		return nil, fmt.Errorf("database URL cannot be empty")
	}

	if migrationsDir == "" {
		return nil, fmt.Errorf("migration directory cannot be empty")
	}

	db, err := builder.dbConnectionFabric.CreateDbConnection(ctx, "postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	var migrator = &migratorImpl{
		MigrationsDir: migrationsDir,
		Db:            db,
	}
	return migrator, nil
}
