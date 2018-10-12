FROM golang:alpine3.8 AS builder
COPY tools/golang/src/qiniu.com /go/src/qiniu.com/
RUN CGO_ENABLED=0 go build -o /go/bin/jvm-exporter qiniu.com/app/jvm-exporter

FROM alpine:3.8
COPY --from=builder /go/bin/jvm-exporter /
COPY tools/golang/src/qiniu.com/app/jvm-exporter/gcMetric.sh /tools/gcMetric.sh
ARG DOCKER_CLI_VERSION="17.06.2-ce"
ENV DOWNLOAD_URL="https://download.docker.com/linux/static/stable/x86_64/docker-$DOCKER_CLI_VERSION.tgz"
RUN apk --update add curl \
    && mkdir -p /tmp/download \
    && curl -L $DOWNLOAD_URL | tar -xz -C /tmp/download \
    && mv /tmp/download/docker/docker /usr/local/bin/ \
    && rm -rf /tmp/download \
    && apk del curl \
    && rm -rf /var/cache/apk/*
ENTRYPOINT [ "/jvm-exporter" ]
