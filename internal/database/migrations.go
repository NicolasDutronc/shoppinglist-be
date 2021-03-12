package database

import (
	"context"
	"fmt"

	"github.com/NicolasDutronc/shoppinglist-be/pkg/mongomigrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/ssh/terminal"
)

// GetMigrations returns the project migrations
func GetMigrations() []*mongomigrate.Migration {
	return []*mongomigrate.Migration{
		{
			ID:   1,
			Name: "backend_role",
			Migrate: func(ctx context.Context, db *mongo.Database) error {
				return db.RunCommand(
					ctx,
					bson.D{
						{
							Key:   "createRole",
							Value: "backend_role",
						},
						{
							Key: "privileges",
							Value: bson.A{
								bson.D{
									{
										Key: "resource",
										Value: bson.D{
											{
												Key:   "db",
												Value: "shoppinglist",
											},
											{
												Key:   "collection",
												Value: "users",
											},
										}},
									{
										Key:   "actions",
										Value: bson.A{"find", "update", "insert", "remove", "changeStream"},
									},
								},

								bson.D{
									{
										Key: "resource",
										Value: bson.D{
											{
												Key:   "db",
												Value: "shoppinglist",
											},
											{
												Key:   "collection",
												Value: "lists",
											},
										}},
									{
										Key:   "actions",
										Value: bson.A{"find", "update", "insert", "remove", "changeStream"},
									},
								},
							},
						},
						{
							Key:   "roles",
							Value: bson.A{},
						},
					},
				).Err()
			},
			Rollback: func(ctx context.Context, db *mongo.Database) error {
				return db.RunCommand(
					ctx,
					bson.D{
						{
							Key:   "dropRole",
							Value: "backend_role",
						},
					},
				).Err()
			},
		},
		{
			ID:   2,
			Name: "backend_user",
			Migrate: func(ctx context.Context, db *mongo.Database) error {
				fmt.Println("Please enter a password for the backend user")
				pwd, err := terminal.ReadPassword(0)
				if err != nil {
					return err
				}

				return db.RunCommand(
					ctx,
					bson.D{
						{
							Key:   "createUser",
							Value: "backend_user",
						},
						{
							Key:   "pwd",
							Value: string(pwd),
						},
						{
							Key:   "roles",
							Value: bson.A{"backend_role"},
						},
					},
				).Err()
			},
			Rollback: func(ctx context.Context, db *mongo.Database) error {
				return db.RunCommand(
					ctx,
					bson.D{
						{
							Key:   "dropUser",
							Value: "backend_user",
						},
					},
				).Err()
			},
		},
		{
			ID:   3,
			Name: "user_collection",
			Migrate: func(ctx context.Context, db *mongo.Database) error {
				if err := db.RunCommand(
					ctx,
					bson.D{
						{
							Key:   "create",
							Value: "users",
						},
					},
				).Err(); err != nil {
					return err
				}

				if _, err := db.Collection("users").Indexes().CreateOne(
					ctx,
					mongo.IndexModel{
						Keys: bson.M{
							"name": 1,
						},
						Options: options.Index().SetUnique(true).SetName("unique username"),
					},
				); err != nil {
					return err
				}

				return nil
			},
			Rollback: func(ctx context.Context, db *mongo.Database) error {
				if err := db.Collection("users").Drop(ctx); err != nil {
					return err
				}

				return nil
			},
		},
		{
			ID:   4,
			Name: "list_collection",
			Migrate: func(ctx context.Context, db *mongo.Database) error {
				return db.RunCommand(
					ctx,
					bson.D{
						{
							Key:   "create",
							Value: "lists",
						},
					},
				).Err()
			},
			Rollback: func(ctx context.Context, db *mongo.Database) error {
				return db.Collection("lists").Drop(ctx)
			},
		},
	}

}
