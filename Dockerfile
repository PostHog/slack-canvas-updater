FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd

RUN go build -o slack-canvas-updater cmd/updater/main.go

CMD ["./slack-canvas-updater"]
