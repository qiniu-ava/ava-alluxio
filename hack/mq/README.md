# 说明

此 kafka 是本地调试时启动的，和线上采用不同的 zookeeper/kafka 镜像，线上版本是为 kubernetes 而设计的，而本地版本暂时直接部署在本地 docker 下。

本地版本的 kafka 基于 [wurstmeister/kafka](https://github.com/wurstmeister/kafka-docker)，zookeeper 基于 [wurstmeister/zookeeper](https://github.com/wurstmeister/zookeeper-docker)

## 运行

请先确保在本地安装了 docker-compose，安装过 docker-for-macos 的设备上已自动安装。

启动 kafka 服务:

``` shell

cd <path_to_hack_mq>
./run.sh

```

## 停止

``` shell

cd <path_to_hack_mq>
docker-compose stop

```
