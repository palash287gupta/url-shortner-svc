FROM golang:1.24

WORKDIR /app

COPY go.mod ./
COPY main.go ./

RUN go build -o url-shortner-svc main.go

EXPOSE 8080

CMD ["./url-shortner-svc"]
