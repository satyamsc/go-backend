FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/app

FROM alpine:3.20
WORKDIR /app
ENV DB_PATH=/data/devices.db
RUN mkdir -p /data
COPY --from=builder /app/app /usr/local/bin/app
COPY docs/swagger/openapi.yaml /app/openapi.yaml
EXPOSE 8080
VOLUME ["/data"]
ENTRYPOINT ["/usr/local/bin/app"]
