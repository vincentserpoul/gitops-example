package postgres

import (
	"archiver/pkg/configuration"
	"archiver/pkg/internal/db"
	"context"
	"database/sql"
	"fmt"
)

func New(ctx context.Context, dbConf *configuration.Database) (*sql.DB, *db.Queries, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbConf.Host, dbConf.Port, dbConf.Username, dbConf.Password, dbConf.DatabaseName, dbConf.SSLMode,
	)

	dbConn, err := sql.Open(
		"postgres",
		psqlInfo,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("opening conn: %w", err)
	}

	if err := dbConn.Ping(); err != nil {
		return nil, nil, fmt.Errorf("ping: %w", err)
	}

	q, err := db.Prepare(ctx, dbConn)
	if err != nil {
		return nil, nil, fmt.Errorf("prepare: %w", err)
	}

	return dbConn, q, nil
}
