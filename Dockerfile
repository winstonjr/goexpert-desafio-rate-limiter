FROM golang:1.23 AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o rl ./cmd/rate_limiter/main.go

FROM scratch
WORKDIR /app
COPY --from=build /app/rl .
ENTRYPOINT ["./rl"]