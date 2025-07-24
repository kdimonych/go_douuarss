package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type DbConectionFabric interface {
	CreateDbConnection(ctx context.Context, driver string, dbURL string) (*sql.DB, error)
}

type dbConectionFabricImpl struct {
}

func NewDbConnectionFabric() DbConectionFabric {
	return &dbConectionFabricImpl{}
}

func (*dbConectionFabricImpl) CreateDbConnection(ctx context.Context, driver, dbURL string) (*sql.DB, error) {
	if driver == "" || dbURL == "" {
		return nil, errors.New("driver and dbURL must be provided")
	}

	db, err := sql.Open(driver, dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("[Info] Database connection established successfully!")
	return db, nil
}
