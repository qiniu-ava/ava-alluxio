# 问题清单

## 部署

### 安装前置依赖时，可能会因为部署 apt 源无法访问在 apt update 时失败。

  请在 /etc/apt 目录下的 source.list 文件中注释掉不需要的源，并将 source.list.d 目录下不需要的文件改为非 .list 后缀即可。

### 部署完 master 机器， master 节点重启后 zookeeper 并没有选举出合适的新节点。

  请确定 zookeeper 中各节点是否处于正常状态，可登录到 zookeeper 容器中通过 zkCli.sh 查看 各节点选举状况。如下使用 jq13 为例:
  ```shell
  root@jq13: docker exec -it alluxio-zk /bin/bash
  zk-container: zkCli.sh
  > ls /leader/alluxio-ro
  ```
  各节点中应包含合适的 `master` 节点 hostname 的列表。尝试启停各 master 节点时，此列表页应当随之更新。
  若 zookeeper 中的选举列表未合理更新，请检查 zookeeper 中的 myid 和其他配置文件是否合适。其中 myid 应该是 1~255 之间的整数，且每个实例都不同。配置文件 ${conf_dir}/zoo.cfg 中应有如下的 server 列表

  ```
  server.<server_1_myid>=<server_1_hostname>:2888:3888
  server.<server_2_myid>=<server_2_hostname>:2888:3888
  server.<server_3_myid>=<server_3_hostname>:2888:3888
  ```

  且需要注意的是，在此列表中对应到当前节点的行中，hostname 应该是 0.0.0.0，这一点是由于我们将 zookeeper 部署在 docker 中引入的。

### master 节点重启后之前的元数据以及绑定的 mountpoint 全部丢失。

  出现此情况有几种可能性:
  1. master 集群配置时使用的 zookeeper 没有生效，如出现此情况，则可以从各 master 的日志中发现每个 master 都试图成为 primary。正常情况下应该是只有一个 primary，其他都是 secondary。
  2. master 的启动脚本中没有设置 --not-format 参数，导致每次启动时都将之前的 journal 清楚掉了。
  3. 各 master 节点所使用的 journal 目录没有实现共享。正常情况下，应该使用某种分布式存储，并在各服务器上挂载指定目录(这里我们使用 cephfs 挂载在 /alluxio-share/ 下)，并将其中的某个子目录(这里我们使用 /alluxio-share/alluxio/journal )映射到 docker 容器中作为 journal 的目录。
