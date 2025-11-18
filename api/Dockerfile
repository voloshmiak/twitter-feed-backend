FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod tidy && \
    go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main ./cmd/api

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add curl

COPY --from=builder /app/main .
COPY --from=builder /app/migrations migrations

EXPOSE 8090

CMD ["./main"]