package app

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeJSON(code int, writer http.ResponseWriter, data interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	_, err = writer.Write(response)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func handleRequest(writer http.ResponseWriter, request *http.Request, data interface{}) bool {
	if request.Header.Get("Content-Type") != "application/json" {
		writer.WriteHeader(http.StatusUnsupportedMediaType)
		return false
	}

	if err := json.NewDecoder(request.Body).Decode(data); err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return false
	}
	if err := request.Body.Close(); err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	return true
}
