<p align="center">
</p>
<p align="center"><h1 align="center">FAMPAYASSIGNMENT</h1></p>
<p align="center"><!-- default option, no dependency badges. -->
</p>
<p align="center">
	<!-- default option, no dependency badges. -->
</p>
<br>

##  Project Structure

```sh
└── fampayassignment/
    ├── Dockerfile
    ├── apis.go
    ├── database.go
    ├── docker-compose.yml
    ├── go.mod
    ├── go.sum
    ├── main.go
    ├── models.go
    └── youtubeSearch.go
```
##  Language, Database and Optimizations used

**Tools used:**
```sh
1.Golang
2.Docker
3.Postgres
```
**Database Schema**
```sh
videotable
video_id (primary key)
title (TEXT NOT NULL)
description (TEXT)
published_at (TIMESTAMP)
thumbnails (JSONB)
created_at (TIMESTAMP)
```

**Optimizations and Improvements**
```sh
1.Created index on published_at column to sort data in desc, to retrive data faster in sorted manner
2.Created full text index on title and descripton to improve search to match better for partial queries
```

###  Prerequisites

Before getting started with fampayassignment, ensure your runtime environment meets the following requirements:

- **Programming Language:** Go
- **Package Manager:** Go modules
- **Container Runtime:** Docker


### How to run the application?
 
**Build from source:**

1. Clone the fampayassignment repository:
```sh
❯ git clone https://github.com/ayushmcodes/fampayassignment
```

2. Navigate to the project directory:
```sh
❯ cd fampayassignment
```

3. You can either run using Goland or Docker:


**Steps to run using `Golang`** &nbsp; [<img align="center" src="https://img.shields.io/badge/Go-00ADD8.svg?style={badge_style}&logo=go&logoColor=white" />](https://golang.org/)

```sh
1.Navigate to database.go and set host as localhost
2.go run .

```


**Using `docker`** &nbsp; [<img align="center" src="https://img.shields.io/badge/Docker-2CA5E0.svg?style={badge_style}&logo=docker&logoColor=white" />](https://www.docker.com/)

```sh
1.Navigate to database.go and set host as db(it is already initalized as db)
2.docker-compose up --build
```

```sh
Once the above cmds are executed successfully, application will start running and
create connection with db, create new video table and indexes, post that our
apis will go live
```

### Project Documentation

**Scheduled YoutTube Search**
```sh
Youtube search happens for a predifiend query for every 10seconds in a paginated manner
and response is stored in postgres

func getYouTubeSearchResults(predefininedSearchQuery string): This is the function which is responsible for performing search.

func rotateAPIKey(): This is the function which is responsible for performing key rotation once the current key's quota limit gets exhausted
```

**API's Provided**
```sh
1.GET video api
http://localhost:8080/api/v1/getVideos?page=1&limit=15

This api will fetch video data in sorted manner which is stored in database

2.Basic Search api
http://localhost:8080/api/v1/videos/search?q=live cricket

This api will provide the video data based on the search query provided.


3.Optimized Search api
http://localhost:8080/api/v2/videos/search?q=live cricket

This api is an optimized version of above api, which provided results which matches the partial search queries.This optimization is done using Full Text Index.
```

Snapshot of tea data In database

![image](https://github.com/user-attachments/assets/d5118c2c-9466-4875-9a41-1d48d59ecdaf)

When using basic search api

![image](https://github.com/user-attachments/assets/adca32e0-87db-417a-a21b-140e08c285b0)

When using optimized search api based on full index search

![image](https://github.com/user-attachments/assets/e053d913-0a20-4a56-9b25-83894a277b2e)

When using get video api

![image](https://github.com/user-attachments/assets/021e361d-e17c-4516-8778-03ccf423d388)

