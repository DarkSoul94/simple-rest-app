package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/DarkSoul94/simple-rest-app/app"
	"github.com/DarkSoul94/simple-rest-app/models"
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

type hBook struct {
	ID           uint64 `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Author       string `json:"author,omitempty"`
	CreationDate string `json:"creation_date,omitempty"`
}

func (h *Handler) CreateBook(c *gin.Context) {
	var (
		newBook hBook
		err     error
	)

	err = c.BindJSON(&newBook)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"status": "error", "error": err})
		return
	}

	err = h.uc.CreateBook(h.toModelBook(newBook))
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"status": "error", "error": err})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"status": "success"})
}

func (h *Handler) GetBook(c *gin.Context) {
	var (
		id   uint64
		book models.Book
		err  error
	)
	id, err = strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"status": "error", "error": "Invalid value in param 'id'"})
		return
	}

	book, err = h.uc.GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"status": "error", "error": err})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"status": "success", "book": h.toHBook(book)})
}

func (h *Handler) GetAllBooks(c *gin.Context) {
	var (
		books  []models.Book
		hBooks []hBook = make([]hBook, 0)
		err    error
	)
	books, err = h.uc.GetBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"status": "error", "error": err})
		return
	}

	for _, book := range books {
		hBooks = append(hBooks, h.toHBook(book))
	}

	c.JSON(http.StatusOK, map[string]interface{}{"status": "success", "books": hBooks})
}

func (h *Handler) UpdateBook(c *gin.Context) {
	var (
		book hBook
		err  error
	)

	err = c.BindJSON(&book)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"status": "error", "error": err})
		return
	}

	err = h.uc.UpdateBook(h.toModelBook(book))
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"status": "error", "error": err})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"status": "success"})
}

func (h *Handler) DeleteBook(c *gin.Context) {
	var (
		id  uint64
		err error
	)

	id, err = strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"status": "error", "error": "Invalid value in param 'id'"})
		return
	}

	err = h.uc.DeleteBook(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"status": "error", "error": err})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"status": "success"})
}

func (h *Handler) toModelBook(book hBook) models.Book {
	date, _ := time.ParseInLocation("2006-01-02", book.CreationDate, time.Local)
	return models.Book{
		ID:           book.ID,
		Name:         book.Name,
		Author:       book.Author,
		CreationDate: date,
	}
}

func (h *Handler) toHBook(book models.Book) hBook {
	return hBook{
		ID:           book.ID,
		Name:         book.Name,
		Author:       book.Author,
		CreationDate: book.CreationDate.Format("02.01.2006"),
	}
}
