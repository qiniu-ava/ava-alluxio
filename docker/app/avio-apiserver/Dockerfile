FROM golang:1.9.3-alpine as builder

WORKDIR /go/src
COPY tools/golang/src/qiniu.com/ /go/src/qiniu.com/
RUN go install qiniu.com/app/avio-apiserver

FROM alpine
ADD http://devtools.dl.atlab.ai/docker/PRC-tz /etc/localtime
COPY --from=builder /go/bin/avio-apiserver /opt/avio-apiserver
CMD ["/opt/avio-apiserver"]
ENV LC_ALL=C.UTF-8
