FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o image-processor ./backend/cmd/main.go

FROM alpine:3.17

WORKDIR /app

COPY --from=builder /app/image-processor .
COPY store-master.json .

EXPOSE 8080

CMD ["./image-processor"]