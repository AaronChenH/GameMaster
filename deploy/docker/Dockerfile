# 构建阶段
FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o game-admin

# 运行阶段
FROM alpine:latest

RUN apk add --no-cache tzdata ca-certificates
COPY --from=builder /app/game-admin /app/
COPY configs/ /app/configs/
COPY templates/ /app/templates/

EXPOSE 8080
CMD ["/app/game-admin"] 