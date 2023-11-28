FROM alpine:latest

RUN apk update && \
    apk upgrade --no-cache && \
    apk add --no-cache tzdata openssh-server && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo 'Asia/Shanghai' >/etc/timezone && \
    rm -rf /var/cache/apk/* && \
    sed -i 's/.*PermitRootLogin.*/PermitRootLogin no/' /etc/ssh/sshd_config && \
    sed -i 's/.*PubkeyAuthentication.*/PubkeyAuthentication yes/' /etc/ssh/sshd_config && \
    sed -i 's/.*PasswordAuthentication.*/PasswordAuthentication no/' /etc/ssh/sshd_config && \
    sed -i 's/.*PermitEmptyPasswords.*/PermitEmptyPasswords no/' /etc/ssh/sshd_config && \
    sed -i 's/.*KbdInteractiveAuthentication.*/KbdInteractiveAuthentication no/' /etc/ssh/sshd_config && \
    sed -i 's/.*X11Forwarding.*/X11Forwarding no/' /etc/ssh/sshd_config && \
    sed -i 's/.*AllowTcpForwarding.*/AllowTcpForwarding yes/' /etc/ssh/sshd_config && \
    sed -i 's/.*GatewayPorts.*/GatewayPorts yes/' /etc/ssh/sshd_config && \
    ssh-keygen -A

COPY ./runner /usr/bin/runner

RUN chmod 700 /usr/bin/runner && \
    chown root /usr/bin/runner && \
    addgroup -g 1001 -S common

WORKDIR /data

ENTRYPOINT [ "/usr/bin/runner" ]
