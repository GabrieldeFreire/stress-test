# Dockerfile
FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . .
RUN go build -o stress_test main.go

FROM alpine:latest as production
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /app/stress_test .
ENTRYPOINT ["./stress_test"]
