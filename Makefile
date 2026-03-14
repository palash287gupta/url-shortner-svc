build:
	go build -o url-shortner-svc ./cmd/url-shortner-svc

test:
	go test ./...

run:
	go run ./cmd/url-shortner-svc/main.go ./cmd/url-shortner-svc/config.go

# Docker commands
docker-build:
	docker build -t url-shortner-svc .

docker-run:
	docker run -d -p 8080:8080 --name url-shortner-svc url-shortner-svc

docker-stop:
	docker stop url-shortner-svc
	docker rm url-shortner-svc

clean:
	rm -f url-shortner-svc
