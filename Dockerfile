FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git && \
    go install github.com/gobuffalo/pop/v6/soda@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd/*.go
FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/app . 
COPY .env .
COPY templates ../templates


EXPOSE 8216

CMD ["./app"]
