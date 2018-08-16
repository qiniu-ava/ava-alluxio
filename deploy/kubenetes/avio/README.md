# avio service 的部署

## 前置依赖

avio service 依赖 jq kubernets 集群 avio-kafka-svc 消息队列服务和 alluxio-proxy-srv 服务。两者的部署流程分别见 [kafka](../mq/README.md) 和 [proxy](../alluxio/README.md)。

## 部署流程

```shell

# 部署 executor
./executor/run.sh

# 部署 walker
./walker/run.sh

# 部署 apiserver
./apiserver/run.sh

```
