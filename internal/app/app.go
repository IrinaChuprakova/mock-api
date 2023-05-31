package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/IrinaChuprakova/mock-api/docs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ping(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		return err
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		return err
	}

	return nil
}

// @title           Swagger UI
// @version         1.0

func Run() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		log.Println(err)
		return
	}

	if err = ping(client); err != nil {
		log.Println(err)
		return
	}

	db := client.Database("cards")

	server := &http.Server{
		Addr:    ":8080",
		Handler: newRouter(db),
	}

	go func() {
		if listenErr := server.ListenAndServe(); listenErr != nil && !errors.Is(listenErr, http.ErrServerClosed) {
			log.Println(listenErr)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err = server.Shutdown(context.Background()); err != nil {
		log.Println(err)
	}
}
