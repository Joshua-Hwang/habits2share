package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	log.Printf("Starting migrations")
	pgs_url := os.Getenv("DATABASE_URL")

	m, err := migrate.New("file://db/migrations", pgs_url)
	if err != nil {
		log.Fatalf("Failed to find migrations: %s\n", err)
	}

	//err = m.Migrate(17)
	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Printf("No migrations to run: %s\n", err)
		} else {
			log.Fatalf("Failed to migrate: %s\n", err)
		}
	}
	log.Printf("Migrations succeeded")
}
