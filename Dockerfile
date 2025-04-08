# 构建阶段
FROM golang:1.21-alpine AS builder

# 安装必要的构建工具
RUN apk add --no-cache git ca-certificates tzdata

# 设置工作目录
WORKDIR /src

# 优先复制依赖文件，利用缓存
COPY go.mod go.sum ./
ENV GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/server ./cmd/server/main.go

# 最终镜像
FROM alpine:3.18

# 安装必要的运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 创建非root用户
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# 创建应用目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/server /app/

# 复制证书文件
COPY --from=builder /src/certs /app/certs

# 设置适当的权限
RUN chown -R appuser:appgroup /app && \
    chmod -R 550 /app && \
    chmod -R 550 /app/certs

# 切换到非root用户
USER appuser

# 暴露gRPC默认端口
EXPOSE 50054

# 启动应用
ENTRYPOINT ["/app/server"]
