package api

import "github.com/gin-gonic/gin"

type Handler interface {
	LoginAdmin(*gin.Context)
}

type handler struct {
}

func NewHandler() Handler {
	return &handler{}
}
