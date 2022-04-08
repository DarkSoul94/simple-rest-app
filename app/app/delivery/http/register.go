package http

import (
	"github.com/DarkSoul94/simple-rest-app/app"
	"github.com/gin-gonic/gin"
)

// RegisterHTTPEndpoints ...
func RegisterHTTPEndpoints(router *gin.RouterGroup, uc app.Usecase) {
	h := NewHandler(uc)

	bookEndpoints := router.Group("/book")
	{
		bookEndpoints.POST("/", h.CreateBook)
		bookEndpoints.GET("/", h.GetAllBooks)
		bookEndpoints.GET("/:id", h.GetBook)
		bookEndpoints.PUT("/", h.UpdateBook)
		bookEndpoints.DELETE("/:id", h.DeleteBook)
	}
}
