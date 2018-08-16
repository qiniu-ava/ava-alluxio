FROM golang:1.9.3-alpine as builder

WORKDIR /go/src
COPY tools/golang/src/qiniu.com/ /go/src/qiniu.com/
RUN go install qiniu.com/app/avio-walker

FROM alpine
ADD http://devtools.dl.atlab.ai/docker/PRC-tz /etc/localtime
RUN apk add librdkafka-dev
COPY --from=builder /go/bin/avio-walker /opt/avio-walker
ENV LC_ALL=C.UTF-8

CMD ["/opt/avio-walker"]
