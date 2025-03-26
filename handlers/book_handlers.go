package handlers

import (
	"encoding/json"
	"net/http"

	"book-api/models"
	"book-api/store"
	"book-api/utils"

	"github.com/gorilla/mux"
)

// BookHandler contains methods to handle book-related requests
type BookHandler struct {
	store *store.BookStore
}

// NewBookHandler creates a new book handler
func NewBookHandler(store *store.BookStore) *BookHandler {
	return &BookHandler{store: store}
}

// GetBooks handles GET /books
func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	books := h.store.GetAllBooks()
	utils.RespondWithJSON(w, http.StatusOK, books)
}

// GetBook handles GET /books/{id}
func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	book, found := h.store.GetBookByID(id)
	if !found {
		utils.RespondWithError(w, http.StatusNotFound, "Book not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, book)
}

// CreateBook handles POST /books
func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&book); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	createdBook, err := h.store.CreateBook(book)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, createdBook)
}

// UpdateBook handles PUT /books/{id}
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var book models.Book
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&book); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	updatedBook, err := h.store.UpdateBook(id, book)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, updatedBook)
}

// DeleteBook handles DELETE /books/{id}
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.store.DeleteBook(id); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Book deleted successfully"})
}

// RegisterRoutes registers all book-related routes
func (h *BookHandler) RegisterRoutes(router *mux.Router) {
	// Order matters! More specific routes should come before general ones
	router.HandleFunc("/books/search", h.SearchBooks).Methods("GET")
	router.HandleFunc("/books", h.GetBooks).Methods("GET")
	router.HandleFunc("/books", h.CreateBook).Methods("POST")
	router.HandleFunc("/books/{id}", h.GetBook).Methods("GET")
	router.HandleFunc("/books/{id}", h.UpdateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", h.DeleteBook).Methods("DELETE")
}
