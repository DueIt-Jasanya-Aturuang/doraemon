package integration

import (
	"github.com/DueIt-Jasanya-Aturuang/doraemon/tests/integration/helper"
	"github.com/ory/dockertest/v3"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	var resources []*dockertest.Resource
	pool := helper.DockerSetup()

	pgResource, dbPg, url := helper.PostgresSetup(pool)
	resources = append(resources, pgResource)

	db = dbPg
	if db == nil {
		panic("db nil")
	}

	helper.MigrationSetup(url, db)

	code := m.Run()

	for _, resource := range resources {
		if err := pool.Purge(resource); err != nil {
			log.Fatal(err)
		}
	}

	os.Exit(code)
}
