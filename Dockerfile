FROM golang:alpine as builder

ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct \
    GO111MODULE=on \
    CGO_ENABLED=1

#设置时区参数
ENV TZ=Asia/Shanghai
RUN echo "https://mirrors.aliyun.com/alpine/v3.4/main/" > /etc/apk/repositories \
    && apk --no-cache add tzdata zeromq \
    && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo '$TZ' > /etc/timezone

WORKDIR /app
ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . /app

RUN go build -o scheduling-system ./cmd

EXPOSE 3000

ENTRYPOINT ["./system"]