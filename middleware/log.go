package middleware

import (
	"log"
	"net/http"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s \n   Token: %s \n   Refresh Token: %s", r.RemoteAddr, r.Method, r.URL, r.Header.Get("X-Token"), r.Header.Get("X-Refresh-Token"))
		next.ServeHTTP(w, r)
	})
}
