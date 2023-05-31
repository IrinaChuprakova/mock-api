# build
FROM golang:1.19-alpine3.16 as build

WORKDIR /app
COPY . .

RUN go mod download
RUN go mod tidy
RUN go build -o ./api cmd/api/main.go

# app
FROM alpine:3.16

COPY --from=build app/api /app/
EXPOSE 8080
WORKDIR /app/
RUN mkdir storage

CMD ["./api"]
