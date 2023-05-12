FROM golang:latest

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

RUN go build -o ./cmd/server/main ./cmd/server/main.go


EXPOSE 8080
CMD ["./cmd/server/main"]


