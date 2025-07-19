FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/cpu-app ./pkg/observability/cpu/main

FROM debian:bullseye-slim
COPY --from=builder /app/cpu-app .
RUN chmod +x /cpu-app
EXPOSE 8484
CMD ["./cpu-app"]