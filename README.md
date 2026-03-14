# URL Shortener Service

A simple URL shortener service built with Go.

## How to Run

```bash
go run ./cmd/url-shortner-svc/main.go
```

Server will start on port 8080.

### Using Docker

Build the image:

```bash
make docker-build
```

Run the container:

```bash
make docker-run
```

Stop the container:

```bash
make docker-stop
```

Check logs:

```bash
docker logs -f url-shortner-svc
```

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

### Get Metrics

Get the top 3 most shortened domains:

```bash
curl http://localhost:8080/metrics
```

Response:

```txt
www.udemy.com: 6
www.youtube.com: 4
en.wikipedia.org: 2
```

## Features

- Generates 6-character short codes
- Returns the same short code if URL is already shortened
- Thread-safe in-memory storage
- Track and display top 3 most shortened domains
