FROM golang:1.19.3-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o ./cmd/zenithd/main ./cmd/zenithd

EXPOSE 8080

CMD ["/app/cmd/zenithd"]