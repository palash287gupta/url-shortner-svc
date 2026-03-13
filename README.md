# URL Shortener Service

A simple URL shortener service built with Go.

## How to Run

```bash
go run main.go
```

The server will start on port 8080.

## API Usage

### Shorten a URL

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://www.google.com"}'
```

Response:
```json
{"short_url":"http://localhost:8080/abc123"}
```

### Access Short URL

Just visit the short URL or:
```bash
curl -L http://localhost:8080/abc123
```

It will redirect you to the original URL.

## Features

- Generates 6-character short codes
- Returns the same short code if URL is already shortened
- Thread-safe in-memory storage
