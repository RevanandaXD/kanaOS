# SVELTE KIT (adapter-node)
FROM node:26-alpine AS web-build
WORKDIR /app/web
RUN npm install -g pnpm@12
COPY app/client/package.json app/client/pnpm-lock.yaml app/client/pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile
COPY app/client/ ./
RUN pnpm build

# GOLANG API
FROM golang:1.27-alpine AS go-build
WORKDIR /src
COPY app/server/go.mod app/server/go.sum ./
RUN go mod download
COPY app/server/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/kanaos ./cmd/kanaos

# RUNTIME: satu container berisi web + API
FROM node:26-alpine
RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=go-build /out/kanaos /usr/local/bin/kanaos
COPY --from=web-build /app/web/build ./web/build
COPY --from=web-build /app/web/node_modules ./web/node_modules
COPY --from=web-build /app/web/package.json ./web/package.json

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

ENV NODE_ENV=production
ENV WEB_UPSTREAM=http://127.0.0.1:3000

EXPOSE 8080
ENTRYPOINT ["entrypoint.sh"]
