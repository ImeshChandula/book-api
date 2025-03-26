package main

import (
	"fmt"
	"log"
	"net/http"

	"book-api/handlers"
	"book-api/store"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the book store
	bookStore, err := store.NewBookStore("books.json")
	if err != nil {
		log.Fatalf("Failed to initialize book store: %v", err)
	}

	// Initialize the book handler
	bookHandler := handlers.NewBookHandler(bookStore)

	// Initialize the main router
	mainRouter := mux.NewRouter()

	// Create a subRouter with the /api prefix
	apiRouter := mainRouter.PathPrefix("/api").Subrouter()

	// Register routes on the API subRouter
	bookHandler.RegisterRoutes(apiRouter)

	// Start the server
	port := "8080"
	fmt.Printf("Server running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mainRouter))
}
