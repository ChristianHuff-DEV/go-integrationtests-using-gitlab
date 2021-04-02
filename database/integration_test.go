package database

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/ChristianHuff-DEV/go-integrationtests-using-gitlab/config"

	"github.com/ory/dockertest/v3"
)

var DB *sql.DB

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("error connection to docker. error: %v", err)
	}

	configuration := config.InitializeConfig("..")

	log.Printf("DatabaseHost: %v", configuration.DatabaseHost)

	postgresContainer, err := pool.Run("postgres", "13", []string{
		"POSTGRES_PASSWORD=" + configuration.DatabasePassword,
		"POSTGRES_USER=" + configuration.DatabaseUser,
		"POSTGRES_DB=" + configuration.DatabaseName,
	})
	if err != nil {
		log.Fatalf("error starting postgres docker container: %s", err)
	}

	port := postgresContainer.GetPort("5432/tcp")

	// FIXME: Should this be wrapped in "pool.Retry(...)"?
	DB, err = Connect(configuration.DatabaseHost, port, configuration.DatabaseName, configuration.DatabaseUser, configuration.DatabasePassword, configuration.DatabaseOpenTimeout)
	if err != nil {
		log.Fatalf("error trying to connect to database: %v", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(postgresContainer); err != nil {
		log.Fatalf("error purging resources of integration tests")
	}

	os.Exit(code)
}
