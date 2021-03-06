package database

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/NicolasDutronc/shoppinglist-be/internal/common"
	"github.com/NicolasDutronc/shoppinglist-be/internal/list"
	"github.com/NicolasDutronc/shoppinglist-be/internal/user"
	"github.com/NicolasDutronc/shoppinglist-be/pkg/mongomigrate"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
)

// GetSeeders returns the seeders of the project
func GetSeeders() []*mongomigrate.Seeder {
	return []*mongomigrate.Seeder{
		{
			Name: "admin_user",
			Seed: func(ctx context.Context, db *mongo.Database) error {
				admin := &user.User{
					BaseModel: common.BaseModel{
						ID:        primitive.NewObjectID(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					Name: "admin",
					Permissions: []*user.Permission{
						{
							Action:     "*",
							ResourceID: "*",
						},
					},
				}

				fmt.Println("Please enter a password for admin")
				pwd, err := terminal.ReadPassword(0)
				if err != nil {
					return err
				}

				encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
				if err != nil {
					return err
				}

				admin.Password = string(encryptedPassword)
				if _, err := db.Collection("users").InsertOne(ctx, admin); err != nil {
					return err
				}

				return nil
			},
		},
		{
			Name: "dummy_users",
			Seed: func(ctx context.Context, db *mongo.Database) error {
				users := []*user.User{
					{
						BaseModel: common.BaseModel{
							ID:        primitive.NewObjectID(),
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
						Name:     "toto",
						Password: "1234",
					},
					{
						BaseModel: common.BaseModel{
							ID:        primitive.NewObjectID(),
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
						Name:     "titi",
						Password: "5678",
					},
				}

				for _, user := range users {
					encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
					if err != nil {
						return err
					}
					user.Password = string(encryptedPassword)
					if _, err := db.Collection("users").InsertOne(ctx, user); err != nil {
						return err
					}
				}

				return nil
			},
		},
		{
			Name: "dummy_lists",
			Seed: func(ctx context.Context, db *mongo.Database) error {
				lists := []*list.Shoppinglist{
					{
						BaseModel: common.BaseModel{
							ID:        primitive.NewObjectID(),
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
						Name: "Bonnes choses",
						Items: []*list.Item{
							{
								Name:     "chocolat",
								Quantity: "500g",
								Done:     false,
							},
							{
								Name:     "baguettes",
								Quantity: "12",
								Done:     true,
							},
						},
					},
					{
						BaseModel: common.BaseModel{
							ID:        primitive.NewObjectID(),
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
						Name: "Le reste...",
						Items: []*list.Item{
							{
								Name:     "légumes",
								Quantity: "500g",
								Done:     false,
							},
							{
								Name:     "salade",
								Quantity: "1",
								Done:     true,
							},
						},
					},
				}

				for _, l := range lists {
					if _, err := db.Collection("lists").InsertOne(ctx, l); err != nil {
						return err
					}
				}

				return nil
			},
		},
		{
			Name: "new_user",
			Seed: func(ctx context.Context, db *mongo.Database) error {
				newUser := &user.User{
					BaseModel: common.BaseModel{
						ID:        primitive.NewObjectID(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					Permissions: []*user.Permission{
						{
							Action:     "*",
							ResourceID: "*",
						},
					},
				}

				fmt.Println("Please enter a name for the new user")
				reader := bufio.NewReader(os.Stdin)
				userName, err := reader.ReadString('\n')
				if err != nil {
					return err
				}

				// remove delimiter
				newUser.Name = strings.Replace(userName, "\n", "", -1)

				fmt.Printf("Please enter a password for %s", userName)
				pwd, err := terminal.ReadPassword(0)
				if err != nil {
					return err
				}

				encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
				if err != nil {
					return err
				}

				newUser.Password = string(encryptedPassword)
				if _, err := db.Collection("users").InsertOne(ctx, newUser); err != nil {
					return err
				}

				return nil
			},
		},
	}
}
