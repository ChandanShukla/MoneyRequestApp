package db

import (
	"database/sql"
	"github.com/gofiber/fiber/v2/log"
)

const create string = `
  CREATE TABLE IF NOT EXISTS client_details (
  id TEXT NOT NULL PRIMARY KEY,
  name TEXT,
  emailAddress TEXT NOT NULL UNIQUE
  );`

const requests string = `CREATE TABLE IF NOT EXISTS request_details ( 
    id TEXT NOT NULL PRIMARY KEY, 
    requestedDate TEXT,
    requesterId TEXT NOT NULL,
    requesterName TEXT,
    requesteeName TEXT,
    requestStatus TEXT,
    expirationDate TEXT,
    amount TEXT,
    requesteeId TEXT NOT NULL,
    invoiceNumber TEXT,
    message TEXT,
    CONSTRAINT fk_clients
    FOREIGN KEY (requesteeId)
    REFERENCES client_details (id)
    );`

const foreign_key string = `PRAGMA foreign_keys = ON`

func SetupDB(db *sql.DB) error {
	queries := []string{foreign_key, create, requests}
	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}
