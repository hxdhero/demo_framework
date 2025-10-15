# 构建阶段
FROM docker.m.daocloud.io/golang:1.24.0-alpine AS builder
ENV GO111MODULE=on
ENV GIN_MODE=release
WORKDIR /var/www/lls_saas_api
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/tmp/go-cache,sharing=private GOPROXY=https://goproxy.cn go mod download
COPY . .
RUN --mount=type=cache,target=/tmp/go-cache,sharing=private GOPROXY=https://goproxy.cn CGO_ENABLED=0 GOOS=linux GOARCH=amd64 /usr/local/go/bin/go build -p 4 -o server cmd/api/main.go

# 运行阶段
FROM docker.m.daocloud.io/golang:1.24.0-alpine
ENV GO111MODULE=on
ENV GIN_MODE=release
WORKDIR /var/www/lls_saas_api
COPY . .
COPY --from=builder /var/www/lls_saas_api/server /var/www/lls_saas_api/server
CMD ["./server"]
