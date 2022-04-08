package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/DarkSoul94/simple-rest-app/app"
	"github.com/DarkSoul94/simple-rest-app/models"
	"github.com/DarkSoul94/simple-rest-app/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sql.DB) app.Repository {
	return &repo{
		db: sqlx.NewDb(db, "postgres"),
	}
}

type dbBook struct {
	ID           uint64    `db:"id"`
	Name         string    `db:"name"`
	Author       string    `db:"author"`
	CreationDate time.Time `db:"creation_date"`
}

func (r *repo) CreateBook(book models.Book) error {
	var (
		dbBook dbBook
		query  string
		err    error
	)

	query = `INSERT INTO books (name, author, creation_date) VALUES ($1, $2, $3)`

	dbBook = r.toDbBook(book)
	_, err = r.db.Exec(query, dbBook.Name, dbBook.Author, dbBook.CreationDate)
	if err != nil {
		logger.LogError(
			"Create book",
			"app/repo/postgres/repo",
			fmt.Sprintf("name: %s, author: %s, date: %s", book.Name, book.Author, book.CreationDate),
			err,
		)
		return err
	}

	return nil
}

func (r *repo) GetBookByID(id uint64) (models.Book, error) {
	var (
		book  dbBook
		query string
		err   error
	)

	query = `SELECT * FROM books WHERE id = $1`
	err = r.db.Get(&book, query, id)
	if err != nil {
		return models.Book{}, err
	}

	return r.toModelsBook(book), nil
}

func (r *repo) GetBooks() ([]models.Book, error) {
	var (
		dbBooks []dbBook      = make([]dbBook, 0)
		mBooks  []models.Book = make([]models.Book, 0)
		query   string
		err     error
	)

	query = `SELECT * FROM books`
	err = r.db.Select(&dbBooks, query)
	if err != nil {
		return nil, err
	}

	for _, book := range dbBooks {
		mBooks = append(mBooks, r.toModelsBook(book))
	}

	return mBooks, nil
}

func (r *repo) UpdateBook(book models.Book) error {
	var (
		dbBook dbBook
		query  string
		err    error
	)

	query = `UPDATE books SET
						name = $1,
						author = $2,
						creation_date = $3
						WHERE id = $4`

	dbBook = r.toDbBook(book)

	_, err = r.db.Exec(query, dbBook.Name, dbBook.Author, dbBook.CreationDate, dbBook.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) DeleteBook(id uint64) error {
	var (
		query string
		err   error
	)

	query = `DELETE FROM books WHERE id = $1`
	_, err = r.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) Close() error {
	return r.db.Close()
}

func (r *repo) toDbBook(book models.Book) dbBook {
	return dbBook{
		ID:           book.ID,
		Name:         book.Name,
		Author:       book.Author,
		CreationDate: book.CreationDate,
	}
}

func (r *repo) toModelsBook(book dbBook) models.Book {
	return models.Book{
		ID:           book.ID,
		Name:         book.Name,
		Author:       book.Author,
		CreationDate: book.CreationDate,
	}
}
