FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.20.5-alpine AS build
#ARG GO_VERSION
ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod ./

RUN --mount=type=cache,target=/go/pkg/mod \
    apk add --no-cache git ca-certificates && \
    go mod download && \
    go mod verify

COPY . .

RUN --mount=readonly,target=. --mount=type=cache,target=/go/pkg/mod \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -a -o /main -ldflags '-w -extldflags "-static"' ./cmd/zenithd


#FROM --platform=${TARGETPLATFORM:-linux/amd64} gcr.io/distroless/static-debian11
FROM gcr.io/distroless/static-debian11

COPY --from=build /main .

EXPOSE 8080
CMD ["./main"]
