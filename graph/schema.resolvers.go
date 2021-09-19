package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/ikalkali/es-golang/entity/books"
	"github.com/ikalkali/es-golang/graph/generated"
	"github.com/ikalkali/es-golang/services"
)

func (r *queryResolver) Books(ctx context.Context) ([]*books.Books, error) {
	fmt.Println(ctx)
	bookList, err := services.BooksService.GetAllBooks(0, 0)
	books := make([]*books.Books, len(bookList))
	if err != nil {
		return nil, err
	}
	for idx, _ := range bookList {
		books[idx] = &bookList[idx]
	}
	return books, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
