package main

import (
	"log"
	"net/http"

	"github.com/makifdb/quick-vid/handlers"
	"github.com/makifdb/quick-vid/middleware"
)

func main() {
	log.Println("Starting server on port 8080")

	router := http.NewServeMux()
	loggerMiddleware := middleware.LoggerMiddleware

	router.HandleFunc("/", loggerMiddleware(handlers.HomeHandler))
	router.HandleFunc("GET /api/transcript/{id}", loggerMiddleware(handlers.TranscriptHandler))

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Fatal(server.ListenAndServe())
}
