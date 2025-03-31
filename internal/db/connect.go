package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"cloud-sprint/config"
	sqlc "cloud-sprint/internal/db/sqlc"
)

func Connect(dbConfig config.DBConfig, log *zap.Logger) (*sql.DB, sqlc.Querier, error) {
	conn, err := sql.Open(dbConfig.Driver, dbConfig.Source)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	err = conn.Ping()
	if err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("database connected successfully")

	queries := sqlc.New(conn)

	return conn, queries, nil
}
