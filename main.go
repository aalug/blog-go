package main

import (
	"database/sql"
	"github.com/aalug/blog-go/api"
	db "github.com/aalug/blog-go/db/sqlc"
	"log"
)

func main() {
	conn, err := sql.Open("postgres", "postgresql://devuser:admin@db:5432/blog_go_db?sslmode=disable")
	if err != nil {
		log.Fatal("cannot connect to the db: ", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start("0.0.0.0:8080")
	if err != nil {
		log.Fatal("cannot start the server:", err)
	}
}
