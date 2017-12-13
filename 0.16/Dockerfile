FROM frolvlad/alpine-glibc:alpine-3.6

MAINTAINER https://github.com/dtandersen/docker_factorio_server

ENV PORT=34197 \
    RCON_PORT=27015 \
    VERSION=0.15.40 \
    SHA1=f79a975f6b8c0ee87e2fa60f7d1f7133f332c3ec

VOLUME /factorio

RUN mkdir /opt && \
    apk add --update --no-cache tini pwgen && \
    apk add --update --no-cache --virtual .build-deps curl && \
    curl -sSL https://www.factorio.com/get-download/$VERSION/headless/linux64 \
        -o /tmp/factorio_headless_x64_$VERSION.tar.xz && \
    echo "$SHA1  /tmp/factorio_headless_x64_$VERSION.tar.xz" | sha1sum -c && \
    tar xf /tmp/factorio_headless_x64_$VERSION.tar.xz --directory /opt && \
    chmod -R ugo=rwx /opt/factorio && \
    rm /tmp/factorio_headless_x64_$VERSION.tar.xz && \
    ln -s /factorio/saves /opt/factorio/saves && \
    ln -s /factorio/mods /opt/factorio/mods && \
    apk del .build-deps && \
    addgroup -g 845 -S fdaemon && \
    adduser -u 845 -G fdaemon -s /bin/sh -SDH fdaemon && \
    chown -R fdaemon:fdaemon /opt/factorio /factorio

USER fdaemon

EXPOSE $PORT/udp $RCON_PORT/tcp

COPY ./docker-entrypoint.sh /

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/docker-entrypoint.sh"]
