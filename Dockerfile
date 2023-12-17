# Build stage
FROM golang:1.21.1-alpine3.14 AS builder
WORKDIR /go/src/app
COPY . .
RUN apk add --no-cache git && \
    go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app

# Final stage
FROM alpine:3.14
COPY --from=builder /go/bin/app /app
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
EXPOSE 8080
ENTRYPOINT [ "/entrypoint.sh" ]