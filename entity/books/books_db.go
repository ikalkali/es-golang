package books

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ikalkali/es-golang/elasticsearch"
	"github.com/ikalkali/es-golang/entity/queries"
)

const (
	indexItems = "books"
	typeItem   = "_doc"
)

func (b *Books) Save() error {
	result, err := elasticsearch.Client.Index(indexItems, typeItem, b)
	if err != nil {
		return err
	}
	b.Id = result.Id
	return nil
}

func (b *Books) Get() error {
	bookId := b.Id
	result, err := elasticsearch.Client.Get(indexItems, typeItem, b.Id)
	if err != nil {
		return err
	}

	bytes, err := result.Source.MarshalJSON()
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &b); err != nil {
		return err
	}
	b.Id = bookId
	return nil
}

func (b *Books) Search(query queries.EsQuery) ([]Books, error) {
	result, err := elasticsearch.Client.Search(indexItems, query.Build())
	if err != nil {
		return nil, err
	}

	books := make([]Books, result.TotalHits())
	for index, hit := range result.Hits.Hits {
		bytes, _ := hit.Source.MarshalJSON()
		var book Books
		if err := json.Unmarshal(bytes, &book); err != nil {
			return nil, err
		}
		book.Id = hit.Id
		books[index] = book
	}

	if len(books) == 0 {
		return nil, errors.New("no result found with the matching criteria")
	}
	return books, nil
}

func (b *Books) GetAllBooks(limit int64, offset int64) ([]Books, error) {
	fmt.Println("Get all books db")
	result, err := elasticsearch.Client.GetAll(indexItems, limit, offset)
	if err != nil {
		return nil, err
	}

	books := make([]Books, len(result.Hits.Hits))
	fmt.Println("TOTAL HITS", result.Hits)
	for index, hit := range result.Hits.Hits {
		bytes, _ := hit.Source.MarshalJSON()
		var book Books
		if err := json.Unmarshal(bytes, &book); err != nil {
			return nil, err
		}
		book.Id = hit.Id
		books[index] = book
	}

	if len(books) == 0 {
		return nil, errors.New("no document found")
	}
	return books, nil
}
