package database

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrateDatabase updates the database to the latest version if it isn't already up-tp-date
//
// Uses github.com/golang-migrate/migrate.
func MigrateDatabase(databaseName, pathToMigrations string, db *sql.DB) error {
	log.Println("begin migration of database")

	// Setup the driver which is used by golang-migrate to perform the migration
	driver, err := postgres.WithInstance(db, &postgres.Config{})

	// Setup the instance on which the migration is executed
	instance, err := migrate.NewWithDatabaseInstance("file:"+pathToMigrations, databaseName, driver)
	if err != nil {
		log.Println("error creating pre-conditions for database migration")
		return err
	}

	// Get the version of the database before the migration. If this is executed on a new database an error is returned
	// that "no migration" exists.
	version, dirty, err := instance.Version()
	if err != nil && err.Error() != "no migration" {
		log.Printf("error determining current database version: %v", err)
		return err
	}

	log.Printf("current database version: %v (dirty: %v)", version, dirty)

	// Run the migration
	if err := instance.Up(); err != nil {
		// If the database is up-to-date this command will return and error with the content "no change". From this
		// applications point of view this is not an error but just means no migration is needed.
		if err.Error() == "no change" {
			log.Println("database is up-to-date")
			return nil
		}
		// Any other error than "no change" indicates an issue while migrating the database
		log.Println("error migrating database")
		return err
	}

	// Get the version of the database after the migration
	version, dirty, err = instance.Version()
	if err != nil {
		log.Println("error determining new database version")
	}

	log.Printf("new database version: %v (dirty: %v)", version, dirty)
	return nil
}
