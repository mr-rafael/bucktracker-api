package api

import (
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
	writer.Write([]byte("OK"))
}
