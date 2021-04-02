package database

import (
	"database/sql"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// TestMigrateDatabase tests the SQL migration scripts by checking the migration version before and after the
// migration is executed.
func TestMigrateDatabase(t *testing.T) {

	// Before the migration runs we expect the version information not to be present
	// an error is therefore expected.
	_, _, err := extractMigrationVersionFromDatabase(DB)
	if err == nil {
		t.Fatalf("error determining migration version from database: %v", err)
	}

	// Run the migration
	err = MigrateDatabase(
		"integration-test",
		// Usually hard coding isn't great. But in this case the test would fail if the migrations are moved. That is
		// a good reason for this test to fail.
		"//migrations",
		DB,
	)
	if err != nil {
		t.Fatalf("error migrating database: %v", err)
	}

	// Based on the migration scripts the expected version is determined
	versionScript, err := extractMigrationVersionFromScripts("migrations")
	if err != nil {
		t.Fatalf("error determining the migration version: %v", err)
	}

	// After the migration we again get the version information from the database
	versionDatabase, dirty, err := extractMigrationVersionFromDatabase(DB)
	if err != nil {
		t.Fatalf("error extracting version from database after migration: %v", err)
	}

	// Check that the database is clean
	if dirty {
		t.Fatalf("the database version after migration was 'dirty'. it was expected to not be 'dirty'.")
	}

	// Check that the version in the database matches the version based on the migration scripts
	if versionScript != versionDatabase {
		t.Fatalf("expected migration version in database and migration scripts to match. version database: %v version scripts: %v", versionDatabase, versionScript)
	}

}

// Extracts the migration version from the database (The version is automatically handled by golang-migrate.)
func extractMigrationVersionFromDatabase(db *sql.DB) (version int64, dirty bool, err error) {
	sqlStatement := "SELECT version, dirty FROM schema_migrations"

	row := db.QueryRow(sqlStatement)

	err = row.Scan(&version, &dirty)
	if err != nil {
		return 0, false, err
	}

	return version, dirty, nil
}

// Extracts the latest migration version based on the migration SQL scripts.
//
// Since all migration scripts are numbered we only have to find the script with the greatest number in its prefix.
func extractMigrationVersionFromScripts(folder string) (version int64, err error) {

	err = filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ".sql" {
			// Extract the version number. i.e. extracts "000001" from "migrations/000001_first-migration.up.sql"
			fileName := strings.Split(path, folder+string(os.PathSeparator))[1]
			versionString := strings.Split(fileName, "_")[0]
			currentVersion, err := strconv.ParseInt(versionString, 10, 64)
			if err != nil {
				return err
			}
			// Check if the version of the current file is greater than the one already cached
			if currentVersion > version {
				version = currentVersion
			}
			return nil
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return version, nil
}
