package middleware

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/mythosmystery/typenotes-go-graphql/auth"
	"github.com/mythosmystery/typenotes-go-graphql/database"
	"github.com/mythosmystery/typenotes-go-graphql/graph/model"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func Auth(db *database.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("x-token")
			if token == "" {
				log.Println("token is empty")
				next.ServeHTTP(w, r)
				return
			}

			claims, err := auth.ParseToken(token, os.Getenv("TOKEN_SECRET"))
			if err != nil || claims == nil {
				log.Println("invalid token, refreshing", err)
				refreshToken := r.Header.Get("x-refresh-token")
				if refreshToken == "" {
					log.Println("refresh token is empty")
					next.ServeHTTP(w, r)
					return
				}
				newToken, newRefreshToken, userID, err := auth.RefreshTokens(refreshToken)
				if err != nil {
					log.Println("error refreshing token", err)
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
				w.Header().Set("X-Token", newToken)
				w.Header().Set("X-Refresh-Token", newRefreshToken)
				user, err := model.GetUserById(userID, db)
				if err != nil {
					log.Println("error getting user", err)
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
				ctx := context.WithValue(r.Context(), userCtxKey, user)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if ForContext(r.Context()) == nil && claims != nil {
				user, err := model.GetUserById(claims.UserID, db)
				if err != nil {
					log.Println("error getting user", err)
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
				ctx := context.WithValue(r.Context(), userCtxKey, user)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func ForContext(ctx context.Context) *model.User {
	user := ctx.Value(userCtxKey)
	if user == nil {
		return nil
	}
	return user.(*model.User)
}
