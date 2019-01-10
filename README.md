<!-- TOC -->

- [ava-alluxio](#ava-alluxio)
    - [前置文档](#前置文档)
    - [本地开发](#本地开发)
        - [avio server](#avio-server)
        - [avio](#avio)
    - [部署](#部署)
        - [前置条件](#前置条件)
        - [部署 zookeeper](#部署-zookeeper)
        - [部署 master](#部署-master)
        - [部署 worker](#部署-worker)
        - [部署 logkit-pro](#部署-logkit-pro)
            - [logkit-pro config](#logkit-pro-config)
            - [pandora log download](#pandora-log-download)
        - [部署 cadvisor](#部署-cadvisor)
        - [部署 node-exporter](#部署-node-exporter)
        - [部署 alluxio-exporter(deprecated)](#部署-alluxio-exporterdeprecated)
        - [多分组情况下部署 alluxio-exporter](#多分组情况下部署-alluxio-exporter)
        - [部署 jvm-exporter](#部署-jvm-exporter)
        - [配置 ava-prometheus](#配置-ava-prometheus)
        - [多分组情况下配置 ava-prometheus](#多分组情况下配置-ava-prometheus)
        - [部署 grafana](#部署-grafana)
    - [工具](#工具)
        - [生成 alluxio 包](#生成-alluxio-包)
        - [生成 alluxio 镜像](#生成-alluxio-镜像)
        - [alluxio 批量加载工具 avio](#alluxio-批量加载工具-avio)
            - [生成 avio 二进制包](#生成-avio-二进制包)
            - [发布 avio 新版本](#发布-avio-新版本)
    - [帮助](#帮助)

<!-- /TOC -->

# ava-alluxio

此代码库主要包含 ava 组提供的 alluxio 服务的本地开发和部署脚本等代码。

## 前置文档

- [`alluxio`](https://www.alluxio.org/docs/1.7/en/index.html)
- [`fuse`](https://github.com/libfuse/libfuse/tree/master/doc)

## 本地开发

### avio server

avio server 包括[多个模块](./docs/avio-server.md)，要在本地开发时正常使用这些模块，需要一次运行如下的脚本:

- 配置 `./hack/mq/env` 文件(参照 `./hack/mq/env.template` )，启动 kafka:

``` shell
cd ./hack/mq
./run.sh # 注意，为了方便 docker-compose 正常启动，请在此脚本所在路径下运行
```

- 配置 `./hack/alluxio/env` 文件(参照 `./hack/mq/env.template` )，启动 alluxio 的 master、worker、proxy

``` shell
./hack/alluxio/master.sh
./hack/alluxio/worker.sh
./hack/alluxio/proxy.sh
```

- 分别启动 avio-apiserver/avio-executor/avio-walker 模块

```shell
./hack/avio/run.sh apiserver # 终端1
./hack/avio/run.sh executor # 终端2
./hack/avio/run.sh walker # 终端3
```

### avio

@TODO

## 部署

目前 alluxio 服务部署在 jq 机房，以下约定部署了 alluxio 服务的机器为 `alluxio 机器组`。 alluxio 的 web dashboard 通过 jq-k8s 转发以提供外网服务。

**部署前，请先确保你已经申请了 `alluxio 机器组` 的 root 权限。*

### 前置条件

请先确保在部署的机器上已经挂载了 cephfs 和 RBD，并安装了 git，将本代码库 clone 到目标机器。

1. 使用 root 账号登录服务器

2. 执行如下命令确认是否已经挂载 cephfs 和 RBD

``` shell
df -h /disk-cephfs # 如果已经挂载则会显示 cephfs 类型的挂载点信息
df -h /disk-rbd # 如果已经挂载 则会显示 /dev/rbd[0-9] 的块设备已挂载在此挂载点
```

3. 执行如下命令以安装 git。

```shell
apt update && apt install -y git
```

4. 克隆本代码库到本机指定位置，请注意克隆前需要将本机的公钥添加到 deploy key 中。

```shell
# mkdir -p /disk1/repos/ 已废弃，现在代码放在 /disk-cephfs/workspace/repos/ava-alluxio 下面
git clone git@github.com:qiniu-ava/ava-alluxio.git /disk-cephfs/workspace/repos/ava-alluxio
```

5. 执行本代码库中 `deploy/env/install.sh` 脚本安装其他必要的依赖并做设置。

### 部署 zookeeper

目前 zookeeper 集群部署在 jq13、jq14、jq15 三个节点。分别在前述三个节点中执行如下步骤部署或者升级 zookeeper:

1. 使用 root 账号登录服务器

2. 进入 deploy 相关目录，并更新最新代码

```shell
cd /disk-cephfs/workspace/repos/ava-alluxio/deploy/docker
git pull
```

3. 执行 zookeeper 的部署脚本

```shell
./alluxio.zookeeper.sh
```

### 部署 master

在确保上述 zookeeper 部署成功后，可在 `alluxio 机器组` 节点上部署分组 master 的实例。分别在需要部署此分组 master 实例的服务器上执行如下步骤:

1. 同部署 zookeeper 中的 *1*。

2. 同部署 zookeeper 中的 *2*。

3. 执行 master 的部署脚本或者更新脚本

```shell
./<group_name>/master.sh start/restart
```

### 部署 worker

在确保上述 zookeeper 部署成功后，可在 `alluxio 机器组` 节点上部署分组 worker 实例。分别在需要部署此分组 worker 实例的服务器上执行如下步骤:

1. 同部署 zookeeper 中的 *1*

2. 同部署 zookeeper 中的 *2*

3. 在 <group_name>/worker.sh 脚本中声明的机器节点上，执行 worker 的部署脚本或者更新脚本

```shell
./<group_name>/worker.sh start/restart
```

4. 在 <group_name>/worker-write.sh 脚本中声明的机器节点上，执行写 worker 的部署脚本或者更新脚本

```shell
./<group_name>/worker-write.sh start/restart
```

5. 依照[下述方式](#%E7%94%9F%E6%88%90-alluxio-%E5%8C%85)生成 alluxio 包，创建 Jira issue 给 Kirk 组相关同学帮忙更新 k8s 集群中各节点上的 alluxio worker 实例

6. 如果有新增分组或者减少分组，需要更新 alluxio dashboard 的代理服务。

### 部署 logkit-pro

logkit-pro 部署执行步骤如下：

1. 同部署 zookeeper 中的 *1*

2. 进入 deploy 相关目录，并更新最新代码

```shell
cd /alluxio-share/workspace/repos/ava-alluxio
git pull
cd deploy/logkit
```

3. 执行 generate_shell_conf.sh 脚本自动生成 logkit 的 runner 配置和 log 读取脚本

```shell
./generate_shell_conf.sh
```

3. 执行 logkit 的初始化脚本logkit_init.sh

```shell
./logkit_init.sh
```

4. 执行 logkit 的启动脚本 start.sh, 启动时需要手动提供有 pandora 服务 AKSK 或提供 aksk 文件路径(文件内容格式见start.sh中的usage帮助)和账号名

```shell
cd ~/logkit/_package
./start.sh --ak=AK --sk=SK or ./start.sh --aksk=AKSKFilePath --mail=MAIL
```

5. 查看logkit日志, 验证服务运行是否成功。日志文件路径在 `logkit.conf` 文件中的 `log` 路径下。

#### logkit-pro config

推荐配置姿势:

1. 本机运行命令

```shell
wget https://pandora-dl.qiniu.com/logkit-pro-local_mac_${LOGKIT_VERSION}.tar.gz && \
tar xzvf logkit-pro-local_mac_${LOGKIT_VERSION}.tar.gz && \
rm logkit-pro-local_mac_${LOGKIT_VERSION}.tar.gz && \
cd logkit-pro-local_mac_${LOGKIT_VERSION}/
```

2. 查看logkit.conf中的bind_host

3. 本机运行logkit:

```shell
./logkit-pro -f logkit.conf
```

4. 本机浏览器进入第二步中的bind_host, 用户名和密码见 auth.conf

5. 点击添加日志收集, 按照网页前端提示, 完成配置, 并以该配置作为模版

6. 根据模版配置以及配置说明进一步定制化, 配置文档：[logkit 配置文档](https://github.com/qiniu/logkit/wiki/logkit%E4%B8%BB%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6)

#### pandora log download

使用方式:

1. 安装requests包

```shell
pip install requests
```

2. 运行logDownload.py

```shell
cd <logkit path>/pandora/
python logDownload.py
```

3. 根据提示输入pandora日志仓库名称, AK, SK, from, scroll。`from`是一个`int`值，表示从第几条日志开始下载; `scroll`是一个时间段，可输入的字符串如 `10s, 20m, 1h`。经测试 `from` 和 `scroll` 属性在该API中无效, 已经向pandora反馈.

4. 输入完成后, 开始下载, 下载时若没有错误会在屏幕打印下载次数和返回的`scroll_id`, 最后一次下载完会打印`readLog Success`表示下载成功; 若下载过程中出错, 会打印出错信息, 并提示`readLog Fail`.

### 部署 cadvisor

执行以下命令:

```shell
cd /alluxio-share/workspace/repos/ava-alluxio/deploy/monitor
./cadvisor.sh start/restart
```

### 部署 node-exporter

在 jq13 ~ 17, jq19 ~ 21 上, 执行以下命令:

```shell
cd /alluxio-share/workspace/repos/ava-alluxio/deploy/monitor
./node-export.sh start/restart
```

### 部署 alluxio-exporter(deprecated)

暂定在 jq17 上, 执行以下命令:

```shell
pip install PyYAML
cd /alluxio-share/workspace/repos/ava-alluxio/deploy/monitor
./alluxio-export.sh start/restart <group_name> <master_ingress or master_host:port>
```

补充:

* alluxio-exporter 启动时需要exporter.yml配置文件，其中需要alluxio组件类型和host地址，可参考 `ava-alluxio/tools/golang/qiniu.com/app/alluxio-exporter/exporter.yml`

* alluxio 容器命名应为 alluxio-master-<group_name> 或 alluxio-worker-<group_name>

### 多分组情况下部署 alluxio-exporter
以pod的方式启动，不需要指定节点。需要：
1. 维护分组信息文件 `alluxio-exporter.yml`。将其创建为一个`configmap`，并且将该configmap挂载到pod上。

```shell
kubectl create configmap alluxio-exporter-config --from-file=./deploy/monitor/alluxio-exporter.yml #创建configmap
```
2. 启动pod

```shell
kubectl create -f ./deploy/monitor/deployment.yml
```

### 部署 jvm-exporter

暂定在 jq17 上, 执行以下命令:

```shell
cd /alluxio-share/workspace/repos/ava-alluxio/deploy/monitor
./jvm-export.sh start/restart
```

### 配置 ava-prometheus

参考文档: [ServiceMoniter 监控配置参考文档](https://cf.qiniu.io/pages/viewpage.action?pageId=37716151)
配置文件可参考：[ServiceMoniter 监控文件](https://gitlab.qiniu.io/ava/ava-deploy/tree/master/apps/alluxio-monitor)

补充:

1. 对于不在k8s集群中的机器，需要使用`Endpoints + Service`来表示`metric`获取源。另外, `ServiceMoniter`中的`endpoints label`需要与`endpoints.yml`中的`ports label`一一对应

2. `ServiceMonitor` 可以选多个 `Service`，通过一个共同的 `label` 就可以了，实际抓取的是 `Service` 后面对应的 `Endpoints`

### 多分组情况下配置 ava-prometheus
对于alluxio exporter, 相对于以上的`Endpoints + Service + ServiceMonitor`配置方式，不同之处在于：把配置方式改变为：`Deployment + Service + ServiceMonitor`。在alluxio exporter中采集各个分组的metric，并且以pod的形式启动。所需的Deployment + Service + ServiceMonitor配置文件均放在deploy/monitor路径下。


### 部署 grafana

前提条件：

1. 已有的 `prometheus` 数据源, 或其他grafana可用的数据源

部署步骤：

1. 编写k8s中grafana服务所需的配置文件，可参考[alluxio-dashboard的配置文件](https://gitlab.qiniu.io/ava/ava-deploy/tree/master/apps/alluxio-dashboard), `configmap` 中 `datasource.ymal` 中 `url` 为数据源的 `url` , `configmap` 中 `grafana.ini` 需要加入 `email` 报警所需的 `[smtp]` 配置, 配置属性如下:

```ini
[smtp]
enabled = true
host = smtpHost:port
user = XXX@XXX.com
password = XXXXXXX
from_address = xuerui@qiniu.com
skip_verify = true
```

2. 将 `service.ymal, ingress.ymal, configmap.ymal, deployment.ymal` 放入同一个文件夹下, 文件夹路径为 `<directory path>`

3. 用 `kubectl` 命令连接到 `xs` 或 `jq` 集群, 然后运行命令:

```shell
kubectl create -f <directory path> # <directory path>为上一步创建的文件夹
```

4. 在浏览器中打开 `ingress.ymal` 中 `host` 路径，验证服务启动成功。

## 工具

### 生成 alluxio 包

```shell
./tools/scripts/build-alluxio.sh
cd .tmp/alluxio
hash=`git rev-parse --short=7 HEAD` && tar zcvf alluxio-1.7.2-${hash}.tar.gz ./alluxio-1.7.2-SNAPSHOT
```

如需要使用本地尚未合并到 [alluxio](github.com/qiniu-ava/alluxio) 中的代码，则可以在执行 build-alluxio.sh 脚本时指定 --local-alluxio 参数，如

```shell
./tools/scrips/build-alluxio.sh --local-alluxio=$HOME/qbox/alluxio
```

关于 ./tools/scrips/build-alluxio.sh 脚本更多的使用细则，请查看脚本中的 usage 说明。

若要生成在 kubernetes 集群中部署的 client 则需要使用最新 tag 的代码来生成压缩包，且以 tag 命名压缩包，如:

```shell

git tag ava-alluxio-<version>             # 先确保代码已全都合并到 qiniu-ava/alluxio 中
git push upstream --tags                  # 确保对应的 tag 推送到了 qiniu-ava/alluxio 源中
git checkout ava-alluxio-<version>
./tools/scripts/build-alluxio.sh
cd .tmp/alluxio
tar_name=`git describe --exact-match --tags $(git rev-parse --short=7 HEAD)` \
  && tar zcvf ${tar_name}.tar.gz ./alluxio-1.7.2-SNAPSHOT \
  && qrsctl put bowen tmp/${tar_name}.tar.gz ./${tar_name}.tar.gz

```

### 生成 alluxio 镜像

1. 在本地登录 reg-xs.qiniu.io registry

``` shell
docker login reg-xs.qiniu.io -u altab -p <atlab_password>
```

2. 执行如下脚本生成本地 alluxio 镜像

```shell
./tools/scripts/build-alluxio.sh
```

3. 执行如下脚本推送 alluxio 镜像

```shell
docker tag alluxio reg.xs.qiniu.io/atlab/alluxio
docker push reg.xs.qiniu.io/atlab/alluxio
```

### alluxio 批量加载工具 avio

avio 相关详情介绍请参阅 [CF文档](https://cf.qiniu.io/pages/viewpage.action?pageId=81986327)

#### 生成 avio 二进制包

``` shell
make tools-golang
```

成功执行后，在 tools/glang/bin/ 文件夹下将有 avio 和 linux_amd64/avio 两个二进制包，前者为当前系统的版本，后者为 linux amd64 架构的版本。

#### 发布 avio 新版本

请先确认已在本地安装好 qrsctl 且以 ava-test 账号登录，并确定要发布的版本号(可以在 tools/golang/src/qiniu.com/app/avio/main.go 中改动)。

``` shell
make tools-golang-deploy
```

发布后可以在 `http://devtools.dl.atlab.ai/ava/cli/avio` 下载，或者 `http://devtools.dl.atlab.ai/ava/cli/avio/<version>/avio-linux` 下载验证指定版本

## 帮助

关于 alluxio 部署过程中的一些注意问题，请查看文档 [Q&A](https://github.com/qiniu-ava/ava-alluxio/blob/develop/questions.md)。
