# syntax=docker/dockerfile:1.7

FROM node:20-alpine AS web-builder
WORKDIR /web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

FROM golang:1.25-alpine AS go-builder
WORKDIR /src
RUN apk add --no-cache build-base
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
COPY --from=web-builder /web/dist ./web/dist
RUN CGO_ENABLED=1 GOOS=linux go build -o /out/fotoboo-api ./cmd/api

FROM alpine:3.21
RUN apk --no-cache add ca-certificates && \
    addgroup -S fotoboo && \
    adduser -S -G fotoboo fotoboo

WORKDIR /app
COPY --from=go-builder /out/fotoboo-api ./fotoboo-api
COPY --from=go-builder /src/web/dist ./web/dist

RUN mkdir -p /app/data/photos && chown -R fotoboo:fotoboo /app

ENV PORT=8080
ENV WEB_DIR=/app/web
ENV STORAGE_PATH=/app/data/photos
ENV DB_PATH=/app/data/fotoboo.db

USER fotoboo

EXPOSE 8080

CMD ["./fotoboo-api"]
