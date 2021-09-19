package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/ikalkali/es-golang/controllers"
	"github.com/ikalkali/es-golang/elasticsearch"
	"github.com/ikalkali/es-golang/graph"
	"github.com/ikalkali/es-golang/graph/generated"
	"github.com/sirupsen/logrus"
)

var (
	router = mux.NewRouter()
)

func main() {
	elasticsearch.Init()

	mapUrls()

	srv := &http.Server{
		Addr:         "127.0.0.1:8081",
		WriteTimeout: 500 * time.Millisecond,
		ReadTimeout:  2 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      router,
	}

	logrus.Info("Starting app")
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

func mapUrls() {
	gqlServer := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	router.HandleFunc("/ping", func(rw http.ResponseWriter, r *http.Request) {
		v := r.URL.Query()
		fmt.Println(v.Get("ikal"))
		logrus.Info("Ping called")
		rw.Write([]byte("Pong"))
	}).Methods(http.MethodGet)

	router.HandleFunc("/books", controllers.BooksController.Create).Methods(http.MethodPost)
	router.HandleFunc("/books/{id}", controllers.BooksController.Get).Methods(http.MethodGet)
	router.HandleFunc("/books/search", controllers.BooksController.Search).Methods(http.MethodPost)
	router.HandleFunc("/all-books", controllers.BooksController.GetAllBooks).Methods(http.MethodGet)
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", gqlServer)
}
