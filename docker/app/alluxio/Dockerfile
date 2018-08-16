FROM ubuntu:16.04

RUN sed -i s/archive.ubuntu.com/mirrors.163.com/g /etc/apt/sources.list && \
    sed -i s/security.ubuntu.com/mirrors.163.com/g /etc/apt/sources.list

RUN apt update && \
    apt install -y curl nginx && \
    apt autoclean && rm -rf /var/lib/apt/lists/*

COPY ./conf /etc/nginx
COPY ./entrypoint.sh /entrypoint.sh

ENTRYPOINT /entrypoint.sh
