FROM golang:1.9.3-alpine as builder

WORKDIR /go/src
ARG VERSION=dev
COPY tools/golang/src/qiniu.com/ /go/src/qiniu.com/
RUN go install -ldflags "-X qiniu.com/ava/com/version.version=${VERSION}" qiniu.com/app/avio-walker

FROM alpine
RUN apk add librdkafka-dev
COPY --from=builder /go/bin/avio-walker /opt/avio-walker
COPY --from=builder /go/src/qiniu.com/app/avio-walker/avio-walker.conf /opt/app.conf
ENV LC_ALL=C.UTF-8

CMD ["/opt/avio-walker", "-f", "/opt/app.conf"]
