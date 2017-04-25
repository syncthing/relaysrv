FROM debian
# Syncthing-Relay Server
MAINTAINER Lars K.W. Gohlke <lkwg82@gmx.de>
# based on https://hub.docker.com/r/blakev/syncthing-relaysrv/

ENV DEBUG           false
ENV SERV_PORT       22067

# 10 mbps
ENV RATE_GLOBAL     10000000
# 500 kbps
ENV RATE_SESSION    500000

ENV TIMEOUT_MSG     1m45s
ENV TIMEOUT_NET     3m30s
ENV PING_INT        1m15s

ENV PUBLIC_IP       "0.0.0.0"
ENV PROVIDED_BY     "syncthing-relay docker container"

RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get install ca-certificates curl wget -y \
    # download latest releases \
    && wget -O /tmp/relaysrv.tar.gz $(curl -s https://api.github.com/repos/syncthing/relaysrv/releases/latest | grep 'browser_' | cut -d\" -f4 | grep linux-amd64) \
    # extract downloads \
    && tar -xzvf /tmp/relaysrv.tar.gz \
    && rm /tmp/relaysrv.tar.gz \
    && mv relaysrv* /srv/relaysrv \
    # cleanup \
    && apt-get purge -y wget curl \
    && apt-get clean -y \
    && apt-get autoclean -y \
    && apt-get autoremove -y \
    && cp -R /usr/share/locale/en\@* /tmp/ && rm -rf /usr/share/locale/* && mv /tmp/en\@* /usr/share/locale/ \
    && rm -rf /var/cache/debconf/*-old \
    && rm -rf /var/lib/apt/lists/* \
    && rm -rf /usr/share/doc/* \
    && rm -rf /tmp/* /var/tmp/* 


RUN useradd relaysrv \
    && chown relaysrv /srv/relaysrv/relaysrv \
    && mkdir /keys \
    && chown relaysrv /keys

USER relaysrv

EXPOSE 22070 ${SERV_PORT}

CMD /srv/relaysrv/relaysrv -keys="/keys" \
    -listen="${PUBLIC_IP}:${SERV_PORT}" \
    -debug="${DEBUG}" \
    -global-rate="${RATE_GLOBAL}" \
    -per-session-rate="${RATE_SESSION}" \
    -message-timeout="${TIMEOUT_MSG}" \
    -network-timeout="${TIMEOUT_NET}" \
    -ping-interval="${PING_INT}" \
    -provided-by="${PROVIDED_BY}"
