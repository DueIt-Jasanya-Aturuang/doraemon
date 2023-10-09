package infra

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func NewPgConn() *sql.DB {
	dbHost := PgHost
	dbPort := PgPort
	dbUser := PgUser
	dbPass := PgPass
	dbName := PgName
	dbSSL := PgSSL
	dbSchema := PgSchema

	fDB := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s search_path=%s sslmode=%s", dbHost, dbPort, dbUser, dbPass, dbName, dbSchema, dbSSL)

	db, err := sql.Open("postgres", fDB)
	if err != nil {
		log.Err(err).Msg("cannot open db")
	}

	ctx, cancel := context.WithTimeout(context.Background(), pgPingTimeOut)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		log.Err(err).Msg("cannot ping db")
	}

	db.SetMaxIdleConns(setMaxIdleConnDB)
	db.SetMaxOpenConns(setMaxOpenConnDB)
	db.SetConnMaxIdleTime(SetConnMaxIdleTimeDB)
	db.SetConnMaxLifetime(setConnMaxLifetimeDB)

	log.Info().Msgf("connection postgres successfully : %s", PgName)
	return db
}

const (
	setMaxIdleConnDB     = 5
	setMaxOpenConnDB     = 100
	SetConnMaxIdleTimeDB = 5 * time.Minute
	setConnMaxLifetimeDB = 60 * time.Minute
	pgPingTimeOut        = 5 * time.Second
)
