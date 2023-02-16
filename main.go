package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/yagoernandes/simplebank/api"
	db "github.com/yagoernandes/simplebank/db/sqlc"
	"github.com/yagoernandes/simplebank/util"
)

var dbConn *sql.DB

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	dbConn, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(dbConn)
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
