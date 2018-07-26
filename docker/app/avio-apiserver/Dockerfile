FROM golang:1.9.3-alpine as builder

WORKDIR /go/src
ARG VERSION=dev
COPY tools/golang/src/qiniu.com/ /go/src/qiniu.com/
RUN go install -ldflags "-X qiniu.com/ava/com/version.version=${VERSION}" qiniu.com/app/avio-executor

FROM alpine
ADD http://devtools.dl.atlab.ai/docker/PRC-tz /etc/localtime
COPY --from=builder /go/bin/avio-executor /opt/avio-executor
CMD ["/opt/avio-executor"]
ENV LC_ALL=C.UTF-8
