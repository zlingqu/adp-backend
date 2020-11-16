# 运行
# FROM docker.dm-ai.cn/devops/base-image-golang-run-env:tag-v0.0.6 AS RUN
# RUN apk update && apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone
FROM alpine:3.12.1
WORKDIR /app
ADD build .
CMD ["/app/adp-backend"]