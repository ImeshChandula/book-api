package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"book-api/models"

	"github.com/google/uuid"
)

// responsible for storing and retrieving books
type BookStore struct {
	sync.RWMutex
	books    map[string]models.Book
	filename string
}

// creates a new book store
func NewBookStore(filename string) (*BookStore, error) {
	store := &BookStore{
		books:    make(map[string]models.Book),
		filename: filename,
	}

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Create empty file if it doesn't exist
		if err := store.saveToFile(); err != nil {
			return nil, err
		}
	} else {
		// Load existing data
		if err := store.loadFromFile(); err != nil {
			return nil, err
		}
	}

	return store, nil
}

// loads books from the file
func (s *BookStore) loadFromFile() error {
	data, err := ioutil.ReadFile(s.filename)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		// Empty file, initialize with empty map
		s.books = make(map[string]models.Book)
		return nil
	}

	return json.Unmarshal(data, &s.books)
}

// saves books to the file
func (s *BookStore) saveToFile() error {
	data, err := json.MarshalIndent(s.books, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(s.filename, data, 0644)
}

// returns all books
func (s *BookStore) GetAllBooks() []models.Book {
	s.RLock()
	defer s.RUnlock()

	books := make([]models.Book, 0, len(s.books))
	for _, book := range s.books {
		books = append(books, book)
	}
	return books
}

// GetBookByID returns a book by its ID
func (s *BookStore) GetBookByID(id string) (models.Book, bool) {
	s.RLock()
	defer s.RUnlock()

	book, found := s.books[id]
	return book, found
}

// CreateBook creates a new book
func (s *BookStore) CreateBook(book models.Book) (models.Book, error) {
	s.Lock()
	defer s.Unlock()

	// Generate a UUID if not provided
	if book.BookID == "" {
		book.BookID = uuid.New().String()
	} else if _, exists := s.books[book.BookID]; exists {
		return models.Book{}, fmt.Errorf("book with ID %s already exists", book.BookID)
	}

	s.books[book.BookID] = book

	if err := s.saveToFile(); err != nil {
		return models.Book{}, err
	}
	return book, nil
}

// UpdateBook updates an existing book
func (s *BookStore) UpdateBook(id string, book models.Book) (models.Book, error) {
	s.Lock()
	defer s.Unlock()

	if _, exists := s.books[id]; !exists {
		return models.Book{}, fmt.Errorf("book with ID %s not found", id)
	}

	// Ensure BookID is not changed
	book.BookID = id
	s.books[id] = book

	if err := s.saveToFile(); err != nil {
		return models.Book{}, err
	}
	return book, nil
}

// DeleteBook deletes a book by its ID
func (s *BookStore) DeleteBook(id string) error {
	s.Lock()
	defer s.Unlock()

	if _, exists := s.books[id]; !exists {
		return fmt.Errorf("book with ID %s not found", id)
	}

	delete(s.books, id)

	return s.saveToFile()
}
