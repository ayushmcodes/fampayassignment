package main

import (
	_ "github.com/lib/pq"
)

// Video represents the structure for video data to store in DB
type Video struct {
	VideoID    string            `json:"video_id"`
	Title      string            `json:"title"`
	Description string           `json:"description"`
	PublishedAt string           `json:"published_at"`
	Thumbnails  map[string]string `json:"thumbnails"` // JSON for thumbnails
}

// YouTubeResponse represents the structure for youtube search api response data
type YouTubeResponse struct {
	Items []struct {
		Id struct {
			VideoId string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			PublishedAt string `json:"publishedAt"`
			Thumbnails struct {
				Default struct {
					URL string `json:"url"`
				} `json:"default"`
				Medium struct {
					URL string `json:"url"`
				} `json:"medium"`
				High struct {
					URL string `json:"url"`
				} `json:"high"`
			} `json:"thumbnails"`
		} `json:"snippet"`
	} `json:"items"`

	NextPageToken string `json:"nextPageToken"`
}