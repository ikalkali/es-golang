package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ikalkali/es-golang/entity/books"
	"github.com/ikalkali/es-golang/entity/queries"
	"github.com/ikalkali/es-golang/services"
	"github.com/ikalkali/es-golang/utils"
	"github.com/sirupsen/logrus"
)

var (
	BooksController booksControllerInterface = &booksController{}
)

type booksControllerInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
	GetAllBooks(w http.ResponseWriter, r *http.Request)
}

type booksController struct {
}

func (cont *booksController) Create(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Error(err.Error())
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()

	var book books.Books
	if err := json.Unmarshal(reqBody, &book); err != nil {
		logrus.Error(err.Error())
		w.Write([]byte(err.Error()))
		return
	}

	result, createErr := services.BooksService.Create(book)
	if createErr != nil {
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Println(result)
	utils.RespondJson(w, r, http.StatusOK, result)
}

func (cont *booksController) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := strings.TrimSpace(vars["id"])

	book, err := services.BooksService.Get(bookId)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Println(book)
	utils.RespondJson(w, r, http.StatusOK, book)
}

func (c *booksController) Search(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Cant read request body"))
		return
	}
	defer r.Body.Close()

	var query queries.EsQuery
	if err := json.Unmarshal(bytes, &query); err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	items, searchErr := services.BooksService.Search(query)
	if searchErr != nil {
		w.Write([]byte(searchErr.Error()))
		return
	}
	utils.RespondJson(w, r, http.StatusOK, items)
}

func (c *booksController) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	var (
		limit  int64
		offset int64
	)
	v := r.URL.Query()
	limitParams := v.Get("limit")
	offsetParams := v.Get("offset")
	if len(limitParams) == 0 {
		limit = 0
	} else {
		limit, _ = strconv.ParseInt(limitParams, 10, 64)
	}
	if len(offsetParams) == 0 {
		offset = 0
	} else {
		offset, _ = strconv.ParseInt(offsetParams, 10, 64)
	}
	books, searchErr := services.BooksService.GetAllBooks(limit, offset)
	if searchErr != nil {
		w.Write([]byte(searchErr.Error()))
		return
	}
	utils.RespondJson(w, r, http.StatusOK, books)
}
