FROM alpine:latest

RUN apk update && \
    apk upgrade --no-cache && \
    apk add --no-cache tzdata openssh-server && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo 'Asia/Shanghai' >/etc/timezone && \
    rm -rf /var/cache/apk/* && \
    sed -i 's/^#?PermitRootLogin\s.*/PermitRootLogin no/g' /etc/ssh/sshd_config && \
    sed -i 's/^#?PubkeyAuthentication\s.*/PubkeyAuthentication yes/g' /etc/ssh/sshd_config && \
    sed -i 's/^#?PasswordAuthentication\s.*/PasswordAuthentication no/g' /etc/ssh/sshd_config && \
    sed -i 's/^#?PermitEmptyPasswords\s.*/PermitEmptyPasswords no/g' /etc/ssh/sshd_config && \
    sed -i 's/^#?X11Forwarding\s.*/X11Forwarding no/g' /etc/ssh/sshd_config

COPY ./runner /usr/bin/runner

RUN chmod 700 /usr/bin/runner
RUN chown root /usr/bin/runner

WORKDIR /data

ENTRYPOINT [ "/usr/bin/runner" ]
