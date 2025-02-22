package handlers

import (
	"log"
	"net/http"
)

func RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("req: %s %s", req.Method, req.URL)
		next.ServeHTTP(w, req)
	})
}
