package main

import (
	"context"
	"fmt"
	"log"

	"github.com/NicolasDutronc/shoppinglist-be/internal/list"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:password@localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("shoplist")
	collection := db.Collection("lists")

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(count)

	repo := &list.MongoDBRepository{
		ShoppinglistsCollection: collection,
	}

	list, err := repo.StoreList(ctx, "courses")
	if err != nil {
		log.Fatal(err)
	}

	chocolat, err := repo.AddItem(ctx, list.ID.Hex(), "chocolat", "800g")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(chocolat)

	result, err := repo.UpdateItem(ctx, list.ID.Hex(), "chocolat", "800g", "chocolat", "700g")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)

	result, err = repo.ToggleItem(ctx, list.ID.Hex(), "chocolat", "700g", true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)

	// result, err = repo.RemoveItem(ctx, list.ID.Hex(), "chocolat", "700g")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	result, err = repo.RemoveAllItems(ctx, list.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)

}
