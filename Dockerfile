FROM golang:alpine as builder

RUN go env -w CGO_ENABLED=0

WORKDIR /build

COPY . .

RUN go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags '-extldflags "-static" -s -w' -o core ./cmd/core

FROM alpine:latest

RUN apk update && \
    apk upgrade --no-cache && \
    apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo 'Asia/Shanghai' >/etc/timezone && \
    rm -rf /var/cache/apk/*

WORKDIR /data

COPY --from=builder /build/core /usr/bin/core

RUN chmod +x /usr/bin/core

ENTRYPOINT [ "/usr/bin/core" ]
