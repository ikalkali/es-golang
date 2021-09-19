package services

import (
	"fmt"

	"github.com/ikalkali/es-golang/entity/books"
	"github.com/ikalkali/es-golang/entity/queries"
)

var (
	BooksService booksServiceInterface = &bookService{}
)

type booksServiceInterface interface {
	Create(books.Books) (*books.Books, error)
	Get(string) (*books.Books, error)
	Search(queries.EsQuery) ([]books.Books, error)
	GetAllBooks(int64, int64) ([]books.Books, error)
}

type bookService struct{}

func (s *bookService) Create(book books.Books) (*books.Books, error) {
	if err := book.Save(); err != nil {
		return nil, err
	}
	return &book, nil
}

func (s *bookService) Get(id string) (*books.Books, error) {
	book := books.Books{Id: id}
	if err := book.Get(); err != nil {
		return nil, err
	}
	return &book, nil
}

func (s *bookService) Search(query queries.EsQuery) ([]books.Books, error) {
	books := books.Books{}
	return books.Search(query)
}

func (s *bookService) GetAllBooks(limit int64, offset int64) ([]books.Books, error) {
	fmt.Println("Get all books service")
	books := books.Books{}
	return books.GetAllBooks(limit, offset)
}
