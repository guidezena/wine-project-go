package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func SendError(w http.ResponseWriter, message string, statusCode int) {
	log.Printf(message)

	data := map[string]string{
		"message": message,
	}

	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
	log.Printf(message)
}
