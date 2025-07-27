FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/memory-app ./pkg/observability/memory/heap/main

FROM debian:bullseye-slim
COPY --from=builder /app/memory-app .
RUN chmod +x /memory-app
EXPOSE 8484
CMD ["./memory-app"]