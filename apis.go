package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

// Function to get paginated videos sorted by published_at
func getPaginatedVideos(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	// Default values if not provided
	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "10"
	}

	// Convert page and limit to integers
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		http.Error(w, "Invalid page parameter", http.StatusBadRequest)
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}

	// Calculate the offset based on page and limit
	offset := (pageInt - 1) * limitInt

	// SQL query to fetch videos sorted by published_at in descending order
	query := `
		SELECT video_id, title, description, published_at, thumbnails
		FROM videotable
		ORDER BY published_at DESC
		LIMIT $1 OFFSET $2
	`

	// Prepare the result set
	rows, err := db.Query(query, limitInt, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying the database: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Slice to hold the video data
	var videos []Video

	for rows.Next() {
		var video Video
		var thumbnailsJSON []byte

		// Scan the row into variables
		err := rows.Scan(&video.VideoID, &video.Title, &video.Description, &video.PublishedAt, &thumbnailsJSON)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning row: %v", err), http.StatusInternalServerError)
			return
		}

		// Unmarshal thumbnails from JSON
		err = json.Unmarshal(thumbnailsJSON, &video.Thumbnails)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error unmarshalling thumbnails: %v", err), http.StatusInternalServerError)
			return
		}

		// Append to the videos slice
		videos = append(videos, video)
	}

	// Check if there were any rows
	if len(videos) == 0 {
		http.Error(w, "No videos found", http.StatusNotFound)
		return
	}

	// Get the total count of videos in the database for pagination metadata
	var totalCount int
	err = db.QueryRow("SELECT COUNT(*) FROM videotable").Scan(&totalCount)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting total count: %v", err), http.StatusInternalServerError)
		return
	}

	//response structure
	response := struct {
		Page       int     `json:"page"`
		Limit      int     `json:"limit"`
		TotalCount int     `json:"total_count"`
		Videos     []Video `json:"videos"`
	}{
		Page:       pageInt,
		Limit:      limitInt,
		TotalCount: totalCount,
		Videos:     videos,
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Send the JSON response
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}

// Function to send data in paginated format, and it will be triggered on api call
func searchVideos(w http.ResponseWriter, r *http.Request) {
	// Get the search term from the query parameters
	searchTerm := r.URL.Query().Get("q")

	// If no search term is provided, return an error
	if searchTerm == "" {
		http.Error(w, "Search term is required", http.StatusBadRequest)
		return
	}

	// SQL query to search for videos by title or description
	query := `
		SELECT video_id, title, description, published_at, thumbnails
		FROM videotable
		WHERE title ILIKE $1 OR description ILIKE $1
	`

	// Prepare the result set
	rows, err := db.Query(query, "%"+searchTerm+"%")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying the database: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Slice to hold the video data
	var videos []Video

	// Loop through the rows and populate the video data
	for rows.Next() {
		var video Video
		var thumbnailsJSON []byte

		// Scan the row into variables
		err := rows.Scan(&video.VideoID, &video.Title, &video.Description, &video.PublishedAt, &thumbnailsJSON)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning row: %v", err), http.StatusInternalServerError)
			return
		}

		// Unmarshal thumbnails from JSON
		err = json.Unmarshal(thumbnailsJSON, &video.Thumbnails)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error unmarshalling thumbnails: %v", err), http.StatusInternalServerError)
			return
		}

		// Append to the videos slice
		videos = append(videos, video)
	}

	// Check if there were any rows
	if len(videos) == 0 {
		http.Error(w, "No videos found", http.StatusNotFound)
		return
	}

	// Response structure
	response := struct {
		Videos []Video `json:"videos"`
	}{
		Videos: videos,
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Send the JSON response
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}

// Optimized search, which searches on title and description.
//Here we have put a full text index on title and descripton
//so that it's able to search videos containing partial match for the search query in either video title or description.
func optimizedSearchVideos(w http.ResponseWriter, r *http.Request) {
	// Get the search term from the query parameters
	searchTerm := r.URL.Query().Get("q")

	// If no search term is provided, return an error
	if searchTerm == "" {
		http.Error(w, "Search term is required", http.StatusBadRequest)
		return
	}

	// SQL query to search for videos by title or description
	query := `
		SELECT video_id, title, description, published_at, thumbnails
		FROM videotable
		WHERE to_tsvector('english', title || ' ' || description) @@ plainto_tsquery('english', $1)
	`

	// Prepare the result set
	rows, err := db.Query(query, "%"+searchTerm+"%")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying the database: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Slice to hold the video data
	var videos []Video

	// Loop through the rows and populate the video data
	for rows.Next() {
		var video Video
		var thumbnailsJSON []byte

		// Scan the row into variables
		err := rows.Scan(&video.VideoID, &video.Title, &video.Description, &video.PublishedAt, &thumbnailsJSON)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning row: %v", err), http.StatusInternalServerError)
			return
		}

		// Unmarshal thumbnails from JSON
		err = json.Unmarshal(thumbnailsJSON, &video.Thumbnails)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error unmarshalling thumbnails: %v", err), http.StatusInternalServerError)
			return
		}

		// Append to the videos slice
		videos = append(videos, video)
	}

	// Check if there were any rows
	if len(videos) == 0 {
		http.Error(w, "No videos found", http.StatusNotFound)
		return
	}

	// Response structure
	response := struct {
		Videos []Video `json:"videos"`
	}{
		Videos: videos,
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Send the JSON response
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}
