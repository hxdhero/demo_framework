#!/bin/sh
echo "开始编译"
CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  /usr/local/go/bin/go build -ldflags="-X main.CompileTime=`date '+%Y-%m-%d_%H:%M:%S'`"  -o server cmd/api/main.go
echo "编译完成"


BuildDev(){
  echo "开始更新"
  supervisorctl stop lls_go:server
  supervisorctl start lls_go:server
  echo "更新完成"
}

BuildProd(){
  docker_img=xxx.cn-hangzhou.cr.aliyuncs.com/shun178/xxx:"$1"
  echo "image is $docker_img"
  docker build -t "$docker_img" -f Dockerfile  .
  docker push "$docker_img"
  docker rmi "$docker_img"
}

if [ "$1" = "prod" ]; then
  BuildProd "$2"
else
  BuildDev
fi