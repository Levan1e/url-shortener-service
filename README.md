# URL Shortener Service

It's a service you can use for links shortening, that is written in Go.

## Project description:

URL Shortener Service is a service for generating short links and restoring original URLs.

### Features:
+ URL shortening with short link output.
+ Restoring the original URL using a short link.
+ Healthcheck endpoint.
+ In-Memory storage with automatic state saving to a JSON file upon shutdown.
+ Docker and Docker Compose support

## Launch options: 
+ ## In-memory launch (default):
        go run ./cmd/shortener/main.go  
### Features of In-Memory Storage:
+ When the service terminates, the state is saved in storage.json.
+ The next time the service is started, the data is automatically loaded from the file.


+ ## Database Postgre launch:
        go run ./cmd/shortener/main.go -storage=postgres

## API endpoints:
+ ### POST /api/v1/shorten
  URL Shortening
  
  ### Request:
        { "url": "http://example.com"}
  ### Retrieval:
        {"shorten_url": "abc1234567"}
  
+ ### GET /api/v1/{shorten_url}
  Getting the original URL from a shortened link.
  
  ### Retrieval:
        { "url": "http://example.com"}
  
+ ### GET /api/v1/health
  Checking the service status

  ### Retrieval:
        {"status": "ok"}     
