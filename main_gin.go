package main

import (
	"database/sql"
	"github.com/aalug/blog-go/api"
	db "github.com/aalug/blog-go/db/sqlc"
	"github.com/aalug/blog-go/utils"
	"log"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load env file: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to the db: ", err)
	}

	store := db.NewStore(conn)

	runGinServer(config, store)
}

func runGinServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start the server:", err)
	}
}
