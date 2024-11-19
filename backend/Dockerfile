# builder
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

RUN apk --no-cache add ca-certificates make
RUN make mysql-migrate-setup


# runner
FROM alpine:3.18
WORKDIR /app
RUN apk --no-cache add ca-certificates make

COPY --from=builder /app/main .
COPY --from=builder /app/Makefile .
COPY --from=builder /app/start.sh .
COPY --from=builder /app/migrations migrations
COPY --from=builder /go/bin/migrate /usr/local/bin

EXPOSE 8080
CMD ["/app/start.sh"]