package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"

	httpSwagger "github.com/swaggo/http-swagger"
)

func newRouter(db *mongo.Database) http.Handler {
	router := chi.NewRouter()

	router.Use(Cors)

	router.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("doc.json")))

	router.Get("/api/storage/{id}", GetImage())
	router.Post("/api/storage", UploadImage())

	router.Get("/api/cards", AllCards(db))
	router.Post("/api/cards", PostCard(db))

	router.Get("/api/cards/favorite", GetFavorites(db))
	router.Post("/api/cards/favorite", PostFavorite(db))
	router.Delete("/api/cards/favorite/{id}", DeleteFavorite(db))

	router.Get("/api/cards/cart", GetCart(db))
	router.Post("/api/cards/cart", PostCart(db))
	router.Delete("/api/cards/cart/{id}", DeleteCart(db))

	router.Get("/api/cards/order", GetOrders(db))
	router.Post("/api/cards/order", PostOrder(db))

	return router
}
