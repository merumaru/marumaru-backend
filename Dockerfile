FROM golang:1.13.4

WORKDIR /app
COPY . .

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -o marumaru cmd/marumaru-backend/marumaru-backend.go

EXPOSE     8080

ENTRYPOINT  ["/app/marumaru"]