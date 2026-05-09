FROM node:20-alpine AS admin-builder

RUN corepack enable && corepack prepare pnpm@latest --activate

WORKDIR /app/admin

COPY frontend/admin/package.json frontend/admin/pnpm-lock.yaml ./

RUN pnpm install --frozen-lockfile

COPY frontend/admin/ .

RUN pnpm build


FROM golang:alpine AS backend-builder

WORKDIR /app/Mou1ght

RUN apk --no-cache add ca-certificates gcc musl-dev && update-ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY --from=admin-builder /app/internal/service/router/adminui/dist ./internal/service/router/adminui/dist

RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o mou1ght-server ./cmd/main.go


FROM alpine AS runtime

LABEL authors="Yoake"

RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /app

COPY --from=backend-builder /app/Mou1ght/mou1ght-server mou1ght-server
COPY --from=backend-builder /etc/ssl/certs/ /etc/ssl/certs/

RUN mkdir -p /app/data/cache/upload

EXPOSE 10420

VOLUME ["/app/data"]

CMD ["./mou1ght-server"]
