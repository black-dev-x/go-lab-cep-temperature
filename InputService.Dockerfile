FROM golang:1.24 AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o main ./input-service/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/main /app/main

EXPOSE 3000
ENTRYPOINT ["./main"]
