package storage

import (
	"context"
	"testing"
)

func TestDown(t *testing.T) {
	t.Skip("This test requires a live database connection and is skipped by default")
	migrationsDir := testMigrationsDir
	builder := NewMigratorBuilder()
	migrator, err := builder.Build(context.Background(), testDBURL, migrationsDir)
	if err != nil {
		t.Fatalf("failed to create migrator: %v", err)
	}
	defer migrator.Close()

	if err := migrator.Down(); err != nil {
		t.Fatalf("goose down failed: %v", err)
	}
	t.Log("Migrations applied successfully!")
}

func TestUp(t *testing.T) {
	t.Skip("This test requires a live database connection and is skipped by default")
	migrationsDir := testMigrationsDir
	builder := NewMigratorBuilder()
	migrator, err := builder.Build(context.Background(), testDBURL, migrationsDir)
	if err != nil {
		t.Fatalf("failed to create migrator: %v", err)
	}
	defer migrator.Close()

	if err := migrator.Up(); err != nil {
		t.Fatalf("goose up failed: %v", err)
	}
	t.Log("Migrations applied successfully!")
}
