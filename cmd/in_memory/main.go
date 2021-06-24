package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/NicolasDutronc/autokey"
	"github.com/NicolasDutronc/shoppinglist-be/internal/api"
	"github.com/NicolasDutronc/shoppinglist-be/internal/config"
	"github.com/NicolasDutronc/shoppinglist-be/internal/list"
	"github.com/NicolasDutronc/shoppinglist-be/internal/user"
	"github.com/NicolasDutronc/shoppinglist-be/pkg/hub"
)

func main() {
	ctx := context.Background()

	// Create an interruption channel
	quit := make(chan struct{}, 1)

	// build config from the environment
	conf, err := config.NewConfig("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	// create data repositories
	listRepository := list.NewInMemoryRepository()
	userRepository := user.NewInMemoryRepository()

	// create and start hub
	// get the current lists to create topics
	currentLists, err := listRepository.FindAllLists(ctx)
	if err != nil {
		log.Fatalf("Error getting the current lists : %v", err.Error())
	}
	topics := make([]hub.Topic, len(currentLists))
	for i, list := range currentLists {
		topics[i] = hub.TopicFromString(list.ID.Hex())
	}
	storage := hub.NewStorage()
	h, err := hub.NewChannelHub(ctx, storage, topics...)
	if err != nil {
		log.Fatal(err)
	}
	go h.Run(ctx, quit)
	defer h.Close(ctx)

	// create and start the key manager
	manager := autokey.NewManager(conf, conf.KeyConfig.Size, conf.KeyConfig.ValidDuration)
	go manager.Start(quit)
	defer manager.Stop()

	// create services
	listSrv := list.NewService(listRepository, h)
	userSrv := user.NewService(userRepository, conf)

	// create admin user
	_, err = userSrv.Store(ctx, "admin", "password", &user.Permission{
		ResourceID: "*",
		Action:     "*",
	})
	if err != nil {
		log.Fatalf("Error creating admin user : %v", err.Error())
	}

	// setup routes
	r := api.SetupRoutes(userSrv, listSrv, h)
	r.StaticFile("/", "./public/index.html")

	// setup server
	server := &http.Server{
		Addr:    conf.BuildServerAdress(),
		Handler: r,
	}

	// start server
	go func() {
		log.Printf("Listening on %s:%s", conf.Server.Hostname, conf.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error : %s", err)
		}
	}()

	// block until a signal is received
	<-quit

	log.Println("Shuting down server...")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown before timeout : %s", err)
	}

	log.Println("Server exiting")
}
