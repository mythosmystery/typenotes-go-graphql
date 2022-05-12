package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/mythosmystery/typenotes-go-graphql/auth"
	"github.com/mythosmystery/typenotes-go-graphql/database"
	"github.com/mythosmystery/typenotes-go-graphql/graph"
	"github.com/mythosmystery/typenotes-go-graphql/graph/generated"
)

const defaultPort = "8080"

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db := database.Connect()

	router := chi.NewRouter()

	router.Use(auth.Middleware(db))

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{DB: db}}))

	router.Handle("/gql", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
