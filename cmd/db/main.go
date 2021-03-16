package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/NicolasDutronc/shoppinglist-be/internal/database"
	"github.com/NicolasDutronc/shoppinglist-be/pkg/mongomigrate"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	ctx := context.Background()

	path, err := filepath.Abs("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	config, err := mongomigrate.NewConfig(path)
	if err != nil {
		log.Fatal(err)
	}

	// database client
	ctxMongo, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctxMongo, options.Client().ApplyURI(config.BuildMongoDBConnexionString()))
	defer client.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// database connection test
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	// get database and collections
	db := client.Database(config.Migrations.DB)

	m := mongomigrate.New(db, config.Migrations.Collection, database.GetMigrations(), database.GetSeeders())
	app := mongomigrate.GetCLI(m)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
