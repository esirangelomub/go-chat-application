package middlewares

import (
	"github.com/go-chi/cors"
	"net/http"
)

func Handler(next http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Adjust this to your needs
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}).Handler(next)
}
