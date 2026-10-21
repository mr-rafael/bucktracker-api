
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server/main.go
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.27.3

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/server /app/server
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/internal/db/migrations /app/internal/db/migrations

EXPOSE 8080

CMD ["/app/server"]
