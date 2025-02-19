package main

import (
	"log"

	"github.com/backend-production-go-1/internal/db"
	"github.com/backend-production-go-1/internal/env"
	"github.com/backend-production-go-1/internal/store"
)

func main() {
	addr := env.GetString("DB_MIGRATOR_ADDR", "")
	// conn, err := db.New(addr, 3, 3, "15m")
	conn, err := db.New(
		addr,
		3,
		3,
		"15m",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	log.Println("database connection pool established")
	store := store.NewPostgresStorage(conn)
	db.Seed(store, conn)
}
