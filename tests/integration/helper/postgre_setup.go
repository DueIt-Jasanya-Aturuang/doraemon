package helper

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/rs/zerolog/log"
)

func PostgresSetup(pool *dockertest.Pool) (*dockertest.Resource, *sql.DB, string) {
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_PASSWORD=duiet_password",
			"POSTGRES_USER=dueit_user",
			"POSTGRES_DB=dueit_db",
			"listen_addresses=*",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Err(err).Msg("could not start postgres")
		os.Exit(1)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://dueit_user:duiet_password@%s/dueit_db?sslmode=disable", hostAndPort)

	log.Info().Msgf("Connecting to database on url: %s", databaseUrl)

	err = resource.Expire(120)
	if err != nil {
		log.Err(err).Msg("cannot set expired resource")
		os.Exit(1)
	}

	pool.MaxWait = 120 * time.Second

	var db *sql.DB

	if err := pool.Retry(func() error {
		db, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Err(err).Msgf("Could not connect to docker: %s", err)
		os.Exit(1)
	}

	return resource, db, databaseUrl
}
