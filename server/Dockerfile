FROM golang:1.27-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/kanaos ./cmd/kanaos

FROM alpine:3.22
RUN apk add --no-cache ca-certificates && addgroup -S kanaos && adduser -S -G kanaos kanaos

WORKDIR /app
COPY --from=build /out/kanaos /usr/local/bin/kanaos

USER kanaos
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -qO- http://127.0.0.1:8080/api/health >/dev/null || exit 1

ENTRYPOINT ["kanaos"]