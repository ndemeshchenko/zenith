FROM golang:1.20.5 as base

WORKDIR /app

COPY go.mod ./

RUN go mod download
RUN go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /main ./cmd/zenithd

FROM gcr.io/distroless/static-debian11

COPY --from=base /main .

#USER small-user:small-user
EXPOSE 8080
CMD ["./main"]
