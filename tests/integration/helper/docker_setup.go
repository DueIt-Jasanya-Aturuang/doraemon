package helper

import (
	"github.com/ory/dockertest/v3"
	"github.com/rs/zerolog/log"
	"os"
)

func DockerSetup() *dockertest.Pool {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal().Err(err)
		os.Exit(1)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatal().Err(err)
		os.Exit(1)
	}

	return pool
}
