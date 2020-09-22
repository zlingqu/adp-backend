# 编译
FROM docker.dm-ai.cn/devops/base-image-golang-compile:master-2-1e85e1e99bfee20f6f0cc5de5a74ce339100d4bd AS COMPILE
WORKDIR /app
ADD src ./
RUN ls && go env -w GOPRIVATE=gitlab.dm-ai.cn && go env -w GO111MODULE=on \
    && export GOPROXY=https://mirrors.aliyun.com/goproxy/ buildDir=./build \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/main -v ./ \
    && cp -r 3rd-api/kubernetes/key ./build && cp -r 3rd-api/jenkins/conf ./build

# 运行
FROM docker.dm-ai.cn/devops/base-image-golang-run-env:tag-v0.0.6 AS RUN
RUN apk update && apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone
WORKDIR /app
COPY --from=COMPILE /app/build .
CMD ["/app/main"]