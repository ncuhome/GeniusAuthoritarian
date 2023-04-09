FROM alpine:latest

RUN apk update && \
    apk upgrade --no-cache && \
    apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo 'Asia/Shanghai' >/etc/timezone && \
    rm -rf /var/cache/apk/*

COPY ./runner /usr/bin/runner

RUN chmod +x /usr/bin/runner

WORKDIR /data

ENTRYPOINT [ "/usr/bin/runner" ]
