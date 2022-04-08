package http

import (
	"github.com/DarkSoul94/simple-rest-app/app"
	"github.com/gin-gonic/gin"
)

// Handler ...
type Handler struct {
	uc app.Usecase
}

// NewHandler ...
func NewHandler(uc app.Usecase) *Handler {
	return &Handler{
		uc: uc,
	}
}

func (h *Handler) CreateBook(c *gin.Context) {

}
func (h *Handler) GetBook(c *gin.Context) {

}
func (h *Handler) GetAllBooks(c *gin.Context) {

}
func (h *Handler) UpdateBook(c *gin.Context) {

}
func (h *Handler) DeleteBook(c *gin.Context) {

}
