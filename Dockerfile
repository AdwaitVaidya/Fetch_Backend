# Use official Golang image as the base image
FROM golang:latest

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o Fetch_Backend .

EXPOSE 8080

CMD ["./Fetch_Backend"]