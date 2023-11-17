FROM alpine:latest

RUN apk update && \
    apk upgrade --no-cache && \
    apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo 'Asia/Shanghai' >/etc/timezone && \
    rm -rf /var/cache/apk/*

WORKDIR /data

COPY ./runner /usr/bin/runner

RUN chmod +x /usr/bin/runner

ENTRYPOINT [ "/usr/bin/runner" ]
