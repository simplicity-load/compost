package database

import (
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3"
	"log"
	"time"

	"database/sql"
	_ "embed"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed session.sql
var sess_schema_file string

// In minutes
const COOKIE_LIFESPAN = 20

const sess_db_file_name = "sessions.sqlite3"

func CreateSessions() (*session.Store, error) {
	db, err := sql.Open("sqlite3", sess_db_file_name)
	if err != nil {
		return nil, err
	}
	if err = createIfNotExists(db, sess_schema_file); err != nil {
		log.Printf("Failed creating schema of database\nErr: %v", err)
		return nil, err
	}
	db.Close()

	log.Println("Connected with Sessions Database")
	storage := sqlite3.New(sqlite3.Config{
		Database:        sess_db_file_name,
		Table:           "sessions",
		Reset:           false,
		GCInterval:      10 * time.Second,
		MaxOpenConns:    100,
		MaxIdleConns:    100,
		ConnMaxLifetime: 1 * time.Second,
	})

	store := session.New(session.Config{
		Storage:    storage,
		Expiration: COOKIE_LIFESPAN * time.Minute,
		KeyLookup:  "cookie:nts-cookie",
	})

	return store, nil
}
