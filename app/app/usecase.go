package app

import "github.com/DarkSoul94/simple-rest-app/models"

type Usecase interface {
	CreateBook(book models.Book) error
	GetBookByID(id uint64) (models.Book, error)
	GetBooks() ([]models.Book, error)
	UpdateBook(book models.Book) error
	DeleteBook(id uint64) error
}
