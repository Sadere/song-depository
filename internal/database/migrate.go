package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/Sadere/song-depository/migrations"

	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Performs DB migrations
func MigrateUp(DSN string) error {
	db, err := sql.Open("pgx", DSN)
	if err != nil {
		return err
	}
	defer db.Close()

	// Ping
	err = db.Ping()
	if err != nil {
		return err
	}

	goose.SetBaseFS(migrations.Migrations)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = goose.RunContext(ctx, "up", db, ".")
	if err != nil {
		return err
	}

	return nil
}
