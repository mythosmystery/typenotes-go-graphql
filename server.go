package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/mythosmystery/typenotes-go-graphql/database"
	"github.com/mythosmystery/typenotes-go-graphql/graph"
	"github.com/mythosmystery/typenotes-go-graphql/graph/generated"
	"github.com/mythosmystery/typenotes-go-graphql/middleware"
	"github.com/rs/cors"
)

const defaultPort = "3001"

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

	router.Use(middleware.Auth(db))
	router.Use(middleware.LogRequest)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{DB: db}}))

	router.Handle("/gql", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "x-token", "x-refresh-token"},
		ExposedHeaders:   []string{"X-Token", "X-Refresh-Token"},
		AllowCredentials: true,
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler.Handler(router)))
}
