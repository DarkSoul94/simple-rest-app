package usecase

import (
	"github.com/DarkSoul94/simple-rest-app/app"
	"github.com/DarkSoul94/simple-rest-app/models"
)

type usecase struct {
	repo app.Repository
}

// NewUsecase ...
func NewUsecase(repo app.Repository) app.Usecase {
	return &usecase{
		repo: repo,
	}
}

func (u *usecase) CreateBook(book models.Book) error {
	return u.repo.CreateBook(book)
}

func (u *usecase) GetBookByID(id uint64) (models.Book, error) {
	return u.repo.GetBookByID(id)
}

func (u *usecase) GetBooks() ([]models.Book, error) {
	return u.repo.GetBooks()
}

func (u *usecase) UpdateBook(book models.Book) error {
	return u.repo.UpdateBook(book)
}

func (u *usecase) DeleteBook(id uint64) error {
	return u.repo.DeleteBook(id)
}
