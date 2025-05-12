FROM golang:1.23-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /app/main .

FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata curl
RUN addgroup -S nonroot && adduser -S nonroot -G nonroot
WORKDIR /app
COPY --from=builder /app/main /app/main
RUN chown -R nonroot:nonroot /app
USER nonroot
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD curl --fail http://localhost:8080/news || exit 1
ENTRYPOINT ["/app/main"]