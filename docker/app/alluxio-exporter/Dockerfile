FROM golang:alpine3.8 AS builder
COPY tools/golang/src/qiniu.com /go/src/qiniu.com/
RUN CGO_ENABLED=0 go build -o /go/bin/alluxio-exporter qiniu.com/app/alluxio-exporter

FROM alpine:3.8
COPY --from=builder /go/bin/alluxio-exporter /
COPY tools/golang/src/qiniu.com/app/alluxio-exporter/exporter.yml /exporter
ENTRYPOINT [ "/alluxio-exporter" ]
