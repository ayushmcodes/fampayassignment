package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	_ "github.com/lib/pq"
)

func main() {
	// Initialize the database connection and create table and indexes
	initConnectionToDB()
	defer closeDB()

	var predefinedSearchQuery string = "cricket"

	// Run the ticker in a goroutine to perform youtube search every 10sec and store data in DB
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Println("-------------TICKER-------------------")
				getYouTubeSearchResults(predefinedSearchQuery)
			}
		}
	}()

	
	// Define the API routes
	//api to get videos store in db
	http.HandleFunc("/api/v1/getVideos", getPaginatedVideos)
	//api to perform search on basis of title and description
	http.HandleFunc("/api/v1/videos/search", searchVideos)
	//api to perform optimized search on basis of title and description
	http.HandleFunc("/api/v2/videos/search", optimizedSearchVideos)
	
	// Start the server
	fmt.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}