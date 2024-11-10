package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)

const (
	host     = "db"
	port     = 5432
	user     = "postgres"
	password = "hello"
	dbname   = "postgres"
)

var db *sql.DB

//Function to create connection with DB
func initConnectionToDB() {
	// Format connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open a connection to the database
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	fmt.Println("Successfully connected to the database!")
	createTableAndSetupIndexes()
}

// Function to close the database connection
func closeDB() {
	if err := db.Close(); err != nil {
		log.Fatal("Error closing database connection: ", err)
	}
}

// Function to create video table and indexes
func createTableAndSetupIndexes() {
	tableCreationQuery := `
	CREATE TABLE IF NOT EXISTS videotable (
		video_id VARCHAR PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT,
		published_at TIMESTAMP WITH TIME ZONE NOT NULL,
		thumbnails JSONB,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	`
	
	queryToCreateIndexOnPublishedDate:=`CREATE INDEX idx_published_at ON videotable (published_at DESC);`

	queryToCreateFullTextIndexOnTitleAndDescription:=`CREATE INDEX idx_video_text_search ON videotable USING gin(to_tsvector('english', title || ' ' || description));`

	


	_, err := db.Exec(tableCreationQuery)
	if err != nil {
		fmt.Println("Error executing query: ", err)
	}

	fmt.Println("Table created successfully or already exists.")

	_, err = db.Exec(queryToCreateIndexOnPublishedDate)
	if err != nil {
		fmt.Println("Error executing query: ", err)
	}

	fmt.Println("Index successfully created on PublishedDate")

	_, err = db.Exec(queryToCreateFullTextIndexOnTitleAndDescription)
	if err != nil {
		fmt.Println("Error executing query: ", err)
	}

	fmt.Println("Full Text Index successfully created on title and description")
}

func saveVideoToDB(video Video) {
	// Prepare the SQL INSERT query
	insertQuery := `
		INSERT INTO videotable (video_id, title, description, published_at, thumbnails)
		VALUES ($1, $2, $3, $4, $5)
	`

	// Convert thumbnails map to JSON
	thumbnailsJSON, err := json.Marshal(video.Thumbnails)
	if err != nil {
		log.Printf("Error marshalling thumbnails to JSON: %v\n", err)
		return
	}

	// Insert the video data into the database
	_, err = db.Exec(insertQuery, video.VideoID, video.Title, video.Description, video.PublishedAt, thumbnailsJSON)
	if err != nil {
		log.Printf("Error inserting video %s into database: %v\n", video.VideoID, err)
	} else {
		fmt.Println("Video inserted into database successfully.")
	}
}