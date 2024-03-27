package sql

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/evmos/validator-status/pkg/config"
)

//go:embed schema.sql
var ddl string

func InitDatabase(ctx context.Context, config *config.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", config.DatabaseFile)
	if err != nil {
		return nil, err
	}

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return nil, err
	}

	return db, nil
}
