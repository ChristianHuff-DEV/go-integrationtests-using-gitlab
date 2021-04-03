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

	// Read the values from the config. (Environmental variables will override what is defined inside the "config.env" file
	configuration := config.InitializeConfig("..")

	// Setup and start the Docker container for the PostgreSQL database (The environmental variables
	//  // are the same as if we would start the container on the command line.)
	postgresContainer, err := pool.Run("postgres", "13", []string{
		"POSTGRES_PASSWORD=" + configuration.DatabasePassword,
		"POSTGRES_USER=" + configuration.DatabaseUser,
		"POSTGRES_DB=" + configuration.DatabaseName,
	})
	if err != nil {
		log.Fatalf("error starting postgres docker container: %s", err)
	}

	// The port mapping for the Docker container is randomly assigned. Here we ask the container on which port
	// the database will be available.
	port := postgresContainer.GetPort("5432/tcp")

	// Establish the connection to the database
	DB, err = Connect(configuration.DatabaseHost, port, configuration.DatabaseName, configuration.DatabaseUser, configuration.DatabasePassword, configuration.DatabaseOpenTimeout)
	if err != nil {
		log.Fatalf("error trying to connect to database: %v", err)
	}

	// Execute the tests
	code := m.Run()

	// Ensure all containers created are deleted againg
	if err := pool.Purge(postgresContainer); err != nil {
		// Even if this fails there is nothing more we can do than logging it. The test execution will be finished after this anyway.
		log.Printf("error purging resources of integration tests: %v", err)
	}

	os.Exit(code)
}
