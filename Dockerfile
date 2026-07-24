FROM golang:1.26-alpine AS build

ENV GOEXPERIMENT=jsonv2 \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build \
        -trimpath \
        -ldflags="-s -w" \
        -o /out/server \
        ./cmd/server

FROM alpine:latest AS runner

RUN apk add --no-cache \
    vips \
    ca-certificates \
    tzdata
    
USER 65532:65532

ENV PORT=5174 \
    ENV=production

COPY --from=build /out/server /server

EXPOSE 5174

ENTRYPOINT ["/server"]