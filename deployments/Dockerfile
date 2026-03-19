# ---- 构建阶段 ----
FROM golang:1.22-alpine AS builder
WORKDIR /app

# 先复制 go.mod，充分利用 Docker 层缓存
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/server ./cmd/server

# ---- 运行阶段（最小镜像）----
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/config ./config

EXPOSE 8080
ENTRYPOINT ["./server"]