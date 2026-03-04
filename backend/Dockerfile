# ---- 构建阶段 ----
FROM golang:1.25-alpine AS builder

# 如果使用代理，可以通过构建参数传入（例如：--build-arg HTTP_PROXY=http://host.docker.internal:7897）
ARG HTTP_PROXY
ARG HTTPS_PROXY

WORKDIR /app

# 启用 Go Modules（默认使用 proxy.golang.org，可通过代理访问）
RUN go env -w GO111MODULE=on

# 复制依赖文件并下载（利用缓存）
COPY go.mod go.sum ./
RUN go mod download

# 复制所有源代码
COPY . .

# 编译主程序（入口为 cmd/main）
# 如果编译其他入口（如 CLIController），可以另外构建，这里只构建主服务
RUN CGO_ENABLED=0 GOOS=linux go build -o my-app ./cmd/main

# ---- 运行阶段 ----
FROM alpine:latest

# 安装 CA 证书（如需要调用外部 HTTPS API）
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/my-app .

# 声明端口（根据你的应用实际监听端口修改，例如 8080）
EXPOSE 8080

# 运行应用
CMD ["./my-app"]