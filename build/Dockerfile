
FROM golang:1.20-alpine
WORKDIR /go/src/github.com/apulis/app/go-tech

ENV GOSUMDB=sum.golang.google.cn

ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.cn,direct"

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk --no-cache add git pkgconfig build-base

ADD go.mod .
ADD go.sum .
RUN go mod download
ADD . .

RUN  GO111MODULE=${GO111MODULE} go build -buildmode=pie  -o /go/bin/go-tech cmd/main.go


FROM alpine:3.11
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk --no-cache add ca-certificates libdrm curl zip
RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata
WORKDIR /app/
COPY --from=0 /go/bin/go-tech .
COPY --from=0 /go/src/github.com/apulis/app/go-tech/configs ./configs
CMD ["./go-tech", "api"]