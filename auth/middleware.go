package auth

import (
	"log"
	"net/http"

	"github.com/mythosmystery/typenotes-go-graphql/database"
)

func Middleware(db *database.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("x-token")
			if token == "" {
				next.ServeHTTP(w, r)
				return
			}
			log.Println("token:", token)
			next.ServeHTTP(w, r)
		})
	}
}
