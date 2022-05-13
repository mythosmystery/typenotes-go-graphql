package middleware

import (
	"context"
	"net/http"
	"os"

	"github.com/mythosmystery/typenotes-go-graphql/auth"
	"github.com/mythosmystery/typenotes-go-graphql/database"
	"github.com/mythosmystery/typenotes-go-graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
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
				next.ServeHTTP(w, r)
				return
			}
			claims, err := auth.ParseToken(token, os.Getenv("TOKEN_SECRET"))
			if err != nil {
				refreshToken := r.Header.Get("x-refresh-token")
				if refreshToken == "" {
					next.ServeHTTP(w, r)
				}
				newToken, newRefreshToken, err := auth.RefreshTokens(refreshToken)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
				}
				w.Header().Set("Access-Control-Expose-Headers", "x-token, x-refresh-token")
				w.Header().Set("x-token", newToken)
				w.Header().Set("x-refresh-token", newRefreshToken)
				res := db.User.FindOne(r.Context(), bson.M{"_id": claims.UserID})
				var user model.User
				err = res.Decode(user)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				ctx := context.WithValue(r.Context(), userCtxKey, user)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
			if ForContext(r.Context()) == nil {
				res := db.User.FindOne(r.Context(), bson.M{"_id": claims.UserID})
				var user model.User
				err = res.Decode(user)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				ctx := context.WithValue(r.Context(), userCtxKey, user)
				next.ServeHTTP(w, r.WithContext(ctx))
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
