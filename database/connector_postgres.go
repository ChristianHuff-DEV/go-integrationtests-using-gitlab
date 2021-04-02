package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// Connect initiates the connection to the database.
//
// If no connection can be established before the timeout an error is returned. (timeout is the duration in seconds)
func Connect(host, port, name, user, password string, timeout int) (*sql.DB, error) {
	log.Printf("connection to database")

	// Trying to connect to the database in 1s intervals
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Calculate if the timeout has been reached
	timeoutExceeded := time.After(time.Duration(time.Duration(timeout)) * time.Second)
	for {
		select {
		// Case for when the timeout has been reached and no connection was established
		case <-timeoutExceeded:
			return nil, fmt.Errorf("database connection failed after %v timeout", timeout)

		// Each tick will try to establish a database connection
		case <-ticker.C:
			log.Printf("trying to connect to database %v:%v/%v", host, port, name)

			databaseUrl := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", user, password, host, port, name)
			db, err := sql.Open("postgres", databaseUrl)
			if err == nil {
				// No error means a connection could be established. To ensure the connection is usable it is pinged.
				err = db.Ping()
				if err == nil {
					// If the ping was successful the database is available
					log.Printf("connected to database")
					return db, nil
				} else {
					log.Printf("error trying connecting to database: %v", err)
				}
			} else {
				log.Printf("error trying connecting to database: %v", err)
			}
		}
	}
}
