package api

import (
	"encoding/json"
	"fmt"
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
	json.NewEncoder(writer).Encode(respBody)
}

func respondWithJSON(writer http.ResponseWriter, data any, statusCode int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	json.NewEncoder(writer).Encode(data)
}

func respondWithCode(writer http.ResponseWriter, statusCode int) {
	writer.WriteHeader(statusCode)
}
