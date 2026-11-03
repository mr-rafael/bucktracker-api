package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type validateResponseErrorParams struct {
	Error string `json:"error"`
}

func respondWithErrorCode(writer http.ResponseWriter, logMessage string, statusCode int) {
	fmt.Printf("[Error]: %v\n", logMessage)
	writer.WriteHeader(statusCode)
}

func respondWithError(writer http.ResponseWriter, logMessage string, apiErrorMessage string, statusCode int) {
	fmt.Printf("[Error]: %v\n", logMessage)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	respBody := validateResponseErrorParams{
		Error: apiErrorMessage,
	}
	if err := json.NewEncoder(writer).Encode(respBody); err != nil {
		log.Printf("[Error]: failed to encode error response: %v", err)
	}
}

func respondWithJSON(writer http.ResponseWriter, data any, statusCode int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	if err := json.NewEncoder(writer).Encode(data); err != nil {
		log.Printf("[Error]: failed to encode JSON response: %v", err)
	}
}

func respondWithCode(writer http.ResponseWriter, statusCode int) {
	writer.WriteHeader(statusCode)
}
