package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	_ "github.com/lib/pq"
)

const baseURL = "https://www.googleapis.com/youtube/v3/search"

var availableAPIKeys = [...]string{
	"AIzaSyCpwME7-qtXZQMiLAofnCWwgkeJi4lhujk",
    "AIzaSyDDsYYDYLJPAhcpVNYwHiSw-b2Cmv32CYk", 
    "AIzaSyCAUDDX1qFJQhi_oEBQ3w105kHieETpeEg",
}

var indexOfCurrentAPIInUse int = 0 

//Function to perform key rotation when current key limit gets exhausted
func rotateAPIKey(){
	indexOfCurrentAPIInUse=(indexOfCurrentAPIInUse+1)%len(availableAPIKeys)
}

var pageToken string = ""
var publishedAfter string = "2024-01-01T00:00:00Z"

//Function that will run in every 10seconds, make an HTTP call to youtube's search api, retrive data and store in db
func getYouTubeSearchResults(predefininedSearchQuery string) {

	fmt.Println("API Key In Use "+availableAPIKeys[indexOfCurrentAPIInUse])
	queryParams := url.Values{}
	queryParams.Set("key", availableAPIKeys[indexOfCurrentAPIInUse])
	queryParams.Set("part", "snippet")
	queryParams.Set("type", "video")
	queryParams.Set("order", "date")
	queryParams.Set("publishedAfter", publishedAfter)
	queryParams.Set("q", predefininedSearchQuery)
	queryParams.Set("pageToken", pageToken)

	// Combine base URL with query parameters
	fullURL := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

	// Make the GET request
	resp, err := http.Get(fullURL)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	//In case youtube api sends 403 due to quota limit reach, perform key rotation
	if resp.StatusCode == http.StatusForbidden{
		fmt.Println("Quota has been reached for "+availableAPIKeys[indexOfCurrentAPIInUse]+", performing key rotation")
		rotateAPIKey()
		return
	}

	// Unmarshal the JSON response into the YouTubeResponse struct
	var youtubeResponse YouTubeResponse
	err = json.Unmarshal(body, &youtubeResponse)
	
	if err != nil {
		fmt.Println("Error unmarshalling JSON response:", err)
		return
	}

	for _, item := range youtubeResponse.Items {
		newVideo := Video{
			VideoID:    item.Id.VideoId,
			Title:      item.Snippet.Title,
			Description: item.Snippet.Description,
			PublishedAt: item.Snippet.PublishedAt,
			Thumbnails: map[string]string{
				"default": item.Snippet.Thumbnails.Default.URL,
				"medium":  item.Snippet.Thumbnails.Medium.URL,
				"high":    item.Snippet.Thumbnails.High.URL,
			},
		}

		// Save video data to postgresql database
		saveVideoToDB(newVideo)
	}
	pageToken = youtubeResponse.NextPageToken
}