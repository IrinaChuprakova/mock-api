package app

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	cardsCollectionName     = "cards"
	cartCollectionName      = "cart"
	favoritesCollectionName = "favorites"
)

type Card struct {
	ID    string  `json:"id" bson:"_id"`
	Name  string  `json:"name" bson:"name"`
	Price float64 `json:"price" bson:"price"`
	Img   string  `json:"img" bson:"img"`
}

type Order struct {
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	Cards     []Card    `json:"cards" bson:"cards"`
}

// AllCards godoc
// @Summary      Получить массив карточек
// @Tags         cards
// @Produce      json
// @Content-Type application/json
// @Success      200 {object} []Card
// @Router       /api/cards [get]
func AllCards(db *mongo.Database) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		cursor, err := db.Collection(cardsCollectionName).Find(request.Context(), bson.D{})
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		var data []Card
		if err = cursor.All(request.Context(), &data); err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writeJSON(http.StatusOK, writer, data)
	}
}

type CardRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Img   string  `json:"img"`
}

// PostCard godoc
// @Summary      Создать карточку
// @Tags         cards
// @Accept       json
// @Produce      json
// @Content-Type application/json
// @param        request body CardRequest true "body"
// @Success      200 {object} Card
// @Router       /api/cards [post]
func PostCard(db *mongo.Database) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var body CardRequest
		if !handleRequest(writer, request, &body) {
			return
		}

		id := uuid.New().String()
		card := Card{
			ID:    id,
			Name:  body.Name,
			Price: body.Price,
			Img:   body.Img,
		}

		_, err := db.Collection(cardsCollectionName).InsertOne(request.Context(), card)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writeJSON(http.StatusCreated, writer, card)
	}
}

// PostFavorite godoc
// @Summary      добавить карточку в избранное
// @Tags         favorite
// @Accept       json
// @Produce      json
// @Content-Type application/json
// @param        request body Card true "body"
// @Success      200 {object} Card
// @Router       /api/cards/favorite [post]
func PostFavorite(db *mongo.Database) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var body Card
		if !handleRequest(writer, request, &body) {
			return
		}

		_, err := db.Collection(favoritesCollectionName).InsertOne(request.Context(), body)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writeJSON(http.StatusCreated, writer, body)
	}
}

// GetFavorites godoc
// @Summary      Получить массив карточек из избранного
// @Tags         favorite
// @Produce      json
// @Content-Type application/json
// @Success      200 {object} []Card
// @Router       /api/cards/favorite [get]
func GetFavorites(db *mongo.Database) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		cursor, err := db.Collection(favoritesCollectionName).Find(request.Context(), bson.D{})
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		var data []Card
		if err = cursor.All(request.Context(), &data); err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writeJSON(http.StatusOK, writer, data)
	}
}

// DeleteFavorite godoc
// @Summary      удалить карточку из избранного
// @Tags         favorite
// @param        id path string true "id"
// @Success      204
// @Router       /api/cards/favorite/{id} [delete]
func DeleteFavorite(db *mongo.Database) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		filter := bson.D{{Key: "_id", Value: chi.URLParam(request, "id")}}
		result, err := db.Collection(favoritesCollectionName).DeleteOne(request.Context(), filter)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if result.DeletedCount == 0 {
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		writer.WriteHeader(http.StatusNoContent)
	}
}

// PostCart godoc
// @Summary      добавить карточку в корзину
// @Tags         cart
// @Accept       json
// @Produce      json
// @Content-Type application/json
// @param        request body Card true "body"
// @Success      200 {object} Card
// @Router       /api/cards/cart [post]
func PostCart(db *mongo.Database) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var body Card
		if !handleRequest(writer, request, &body) {
			return
		}

		_, err := db.Collection(cartCollectionName).InsertOne(request.Context(), body)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writeJSON(http.StatusCreated, writer, body)
	}
}

// GetCart godoc
// @Summary      Получить массив карточек из корзины
// @Tags         cart
// @Produce      json
// @Content-Type application/json
// @Success      200 {object} []Card
// @Router       /api/cards/cart [get]
func GetCart(db *mongo.Database) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		cursor, err := db.Collection(cartCollectionName).Find(request.Context(), bson.D{})
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		var data []Card
		if err = cursor.All(request.Context(), &data); err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writeJSON(http.StatusOK, writer, data)
	}
}

// DeleteCart godoc
// @Summary      удалить карточку из корзины
// @Tags         cart
// @param        id path string true "id"
// @Success      204
// @Router       /api/cards/cart/{id} [delete]
func DeleteCart(db *mongo.Database) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		filter := bson.D{{Key: "_id", Value: chi.URLParam(request, "id")}}
		result, err := db.Collection(cartCollectionName).DeleteOne(request.Context(), filter)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if result.DeletedCount == 0 {
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		writer.WriteHeader(http.StatusNoContent)
	}
}

type orderResponse struct {
	CreatedAt string `json:"created_at"`
	Cards     []Card `json:"cards"`
}

// GetOrders godoc
// @Summary      Получить массив карточек заказов
// @Tags         order
// @Produce      json
// @Content-Type application/json
// @Success      200 {object} []orderResponse
// @Router       /api/cards/order [get]
func GetOrders(db *mongo.Database) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		cursor, err := db.Collection("orders").Find(request.Context(), bson.D{})
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		var data []Order
		if err = cursor.All(request.Context(), &data); err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		result := make(map[string][]Card)
		for _, item := range data {
			key := item.CreatedAt.Format("02.01.2006")
			result[key] = append(result[key], item.Cards...)
		}

		var response []orderResponse
		for key, cards := range result {
			response = append(response, orderResponse{
				CreatedAt: key,
				Cards:     cards,
			})
		}

		sort.Slice(response, func(i, j int) bool {
			return response[i].CreatedAt > response[j].CreatedAt
		})

		writeJSON(http.StatusOK, writer, response)
	}
}

type orderRequest struct {
	Cards []Card `json:"cards"`
}

// PostOrder godoc
// @Summary      добавить карточку в список заказов
// @Tags         order
// @Accept       json
// @Produce      json
// @Content-Type application/json
// @param        request body orderRequest true "body"
// @Success      201
// @Router       /api/cards/order [post]
func PostOrder(db *mongo.Database) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var body orderRequest
		if !handleRequest(writer, request, &body) {
			return
		}

		order := Order{
			CreatedAt: time.Now(),
			Cards:     body.Cards,
		}

		_, err := db.Collection("orders").InsertOne(request.Context(), order)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusCreated)
	}
}

type imageResponse struct {
	URL string `json:"url"`
}

// UploadImage godoc
// @Summary      Загрузить картинку
// @Tags         storage
// @Accept       multipart/form-data
// @Produce      json
// @param        file formData file true "file"
// @Success      201 {object} imageResponse
// @Router       /api/storage [post]
func UploadImage() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		file, header, err := request.FormFile("file")
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		data := make([]byte, header.Size)

		if _, err = file.Read(data); err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = file.Close(); err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = request.Body.Close(); err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		savedFile, err := os.OpenFile("./storage/"+header.Filename, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer func() {
			if err := savedFile.Close(); err != nil {
				log.Println(err)
			}
		}()

		if _, err = savedFile.Write(data); err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)

		response := imageResponse{
			URL: "http://localhost:8080/api/storage/" + header.Filename,
		}

		bytes, err := json.Marshal(response)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err = writer.Write(bytes); err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func GetImage() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		filename := chi.URLParam(request, "id")
		file, err := os.OpenFile("./storage/"+filename, os.O_RDONLY, os.ModePerm)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				writer.WriteHeader(http.StatusNotFound)
				return
			}
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer func() {
			if err := file.Close(); err != nil {
				log.Println(err)
			}
		}()

		stats, err := file.Stat()
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		data := make([]byte, stats.Size())
		if _, err = file.Read(data); err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		contentType := ""
		switch {
		case strings.Contains(filename, "jpg"):
			contentType = "image/jpeg"
		case strings.Contains(filename, "png"):
			contentType = "image/png"
		}

		writer.Header().Set("Content-Type", contentType)
		writer.Header().Set("Content-Length", strconv.Itoa(int(stats.Size())))
		writer.WriteHeader(http.StatusOK)
		if _, err = writer.Write(data); err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
		}
	}
}
