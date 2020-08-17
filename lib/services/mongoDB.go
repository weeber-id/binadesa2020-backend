package services

import (
	"log"
	"os"

	"github.com/go-bongo/bongo"
)

var MDB *bongo.Connection

func MongodbConnection() {
	user := os.Getenv("MONGODB_USER")
	pass := os.Getenv("MONGODB_PASS")
	host := os.Getenv("MONGODB_HOST")
	database := os.Getenv("MONGODB_DATABASE")

	config := &bongo.Config{
		ConnectionString: user + ":" + pass + "@" + host + "/" + database,
		Database:         database,
	}

	conn, err := bongo.Connect(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	MDB = conn
}
