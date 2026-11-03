package api

import (
	"log"
	"net/http"
)

type AdminHandler struct {
}

func NewAdminHandler() AdminHandler {
	return AdminHandler{}
}

func (admin *AdminHandler) HandlerHealthZ(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := writer.Write([]byte("OK")); err != nil {
		log.Printf("healthz: failed to write response: %v", err)
	}
}
