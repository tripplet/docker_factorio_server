FROM frolvlad/alpine-glibc:alpine-3.5

MAINTAINER https://github.com/dtandersen/docker_factorio_server

ENV PORT=34197 \
    VERSION=0.15.11 \
    SHA1=bf20dc746f49f13684c91558b99aa3e14e6da2c5

RUN mkdir /opt && \
    apk add --update --no-cache tini pwgen && \
    apk add --update --no-cache --virtual .build-deps curl && \
    curl -sSL https://www.factorio.com/get-download/$VERSION/headless/linux64 \
        -o /tmp/factorio_headless_x64_$VERSION.tar.xz && \
    echo "$SHA1  /tmp/factorio_headless_x64_$VERSION.tar.xz" | sha1sum -c && \
    tar xf /tmp/factorio_headless_x64_$VERSION.tar.xz --directory /opt && \
    rm /tmp/factorio_headless_x64_$VERSION.tar.xz && \
    ln -s /factorio/saves /opt/factorio/saves && \
    ln -s /factorio/mods /opt/factorio/mods && \
    ln -s /factorio/log /opt/factorio/log && \
    apk del .build-deps

VOLUME /factorio

EXPOSE $PORT/udp 27015/tcp

COPY ./docker-entrypoint.sh /

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/docker-entrypoint.sh"]
