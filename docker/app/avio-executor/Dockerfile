FROM golang:1.9.3-alpine as builder

WORKDIR /go/src
COPY tools/golang/src/qiniu.com/ /go/src/qiniu.com/
RUN go install qiniu.com/app/avio-executor

FROM alpine
ADD http://devtools.dl.atlab.ai/docker/PRC-tz /etc/localtime
RUN apk add librdkafka-dev
COPY --from=builder /go/bin/avio-executor /opt/avio-executor
CMD ["/opt/avio-executor"]
ENV LC_ALL=C.UTF-8
