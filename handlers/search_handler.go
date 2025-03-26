package handlers

import (
	"net/http"
	"strings"
	"sync"

	"book-api/models"
	"book-api/utils"
)

// SearchBooks handles GET /books/search?q=keyword
func (h *BookHandler) SearchBooks(w http.ResponseWriter, r *http.Request) {
	// Get the query parameter
	query := r.URL.Query().Get("q")
	if query == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing search query parameter 'q'")
		return
	}

	// Convert query to lowercase for case-insensitive search
	query = strings.ToLower(query)

	// Get all books
	allBooks := h.store.GetAllBooks()

	// If few books, don't use concurrency
	if len(allBooks) < 100 {
		results := make([]models.Book, 0)
		for _, book := range allBooks {
			if strings.Contains(strings.ToLower(book.Title), query) ||
				strings.Contains(strings.ToLower(book.Description), query) {
				results = append(results, book)
			}
		}
		utils.RespondWithJSON(w, http.StatusOK, results)
		return
	}

	// For larger datasets, use concurrency
	// Determine number of workers based on CPU count or other factors
	numWorkers := 4

	// Calculate batch size
	batchSize := (len(allBooks) + numWorkers - 1) / numWorkers // Ceiling division

	// Create a channel to receive results
	resultChan := make(chan models.Book)

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Start workers to process batches concurrently
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		// Calculate start and end indices for this batch
		startIdx := i * batchSize
		endIdx := startIdx + batchSize
		if endIdx > len(allBooks) {
			endIdx = len(allBooks)
		}

		// Skip if this batch is empty
		if startIdx >= len(allBooks) {
			wg.Done()
			continue
		}

		// Create a batch of books
		batch := allBooks[startIdx:endIdx]

		// Start a goroutine to process this batch
		go func(books []models.Book) {
			defer wg.Done()

			// Search through books in this batch
			for _, book := range books {
				if strings.Contains(strings.ToLower(book.Title), query) ||
					strings.Contains(strings.ToLower(book.Description), query) {
					resultChan <- book
				}
			}
		}(batch)
	}

	// Start a goroutine to close the result channel when all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results from the channel
	results := make([]models.Book, 0)
	for book := range resultChan {
		results = append(results, book)
	}

	utils.RespondWithJSON(w, http.StatusOK, results)
}
