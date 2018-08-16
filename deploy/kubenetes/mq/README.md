# avio-kafka-svc 服务部署流程

在配置有 jq kubernetes 权限的主机上执行如下命令:

```shell

kubectl create -f ./zookeeper.yml   # kafka 依赖的 zk 服务
kubectl create -f ./kafka.yml

```
