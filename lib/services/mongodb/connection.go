package mongodb

import (
	"binadesa2020-backend/lib/clog"
	vars "binadesa2020-backend/lib/variable"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Client mongoDB
var Client *mongo.Client

// DB vokasi_binadesa
var DB *mongo.Database

// Connection to mongoDB
func Connection(ctx context.Context) {
	var err error
	config := vars.MongoConfig

	// URI := "mongodb://" + user + ":" + pass + "@" + host + "/" + database
	URI := "mongodb://" + config.User + ":" + config.Password + "@" + config.Host

	Client, err = mongo.Connect(ctx, options.Client().ApplyURI(URI))
	if err != nil {
		clog.Fatal(err, "connecting database")
	}

	err = Client.Ping(ctx, readpref.Primary())
	if err != nil {
		clog.Fatal(err, "ping to database")
	}

	DB = Client.Database(config.Database)
}
