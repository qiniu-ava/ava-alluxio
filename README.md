<!-- TOC -->

- [ava-alluxio](#ava-alluxio)
  - [前置文档](#%E5%89%8D%E7%BD%AE%E6%96%87%E6%A1%A3)
  - [本地开发](#%E6%9C%AC%E5%9C%B0%E5%BC%80%E5%8F%91)
    - [avio server](#avio-server)
    - [avio](#avio)
  - [部署](#%E9%83%A8%E7%BD%B2)
    - [前置条件](#%E5%89%8D%E7%BD%AE%E6%9D%A1%E4%BB%B6)
    - [部署 zookeeper](#%E9%83%A8%E7%BD%B2-zookeeper)
    - [部署 master](#%E9%83%A8%E7%BD%B2-master)
    - [部署 worker](#%E9%83%A8%E7%BD%B2-worker)
  - [工具](#%E5%B7%A5%E5%85%B7)
    - [生成 alluxio 包](#%E7%94%9F%E6%88%90-alluxio-%E5%8C%85)
    - [生成 alluxio 镜像](#%E7%94%9F%E6%88%90-alluxio-%E9%95%9C%E5%83%8F)
    - [alluxio 批量加载工具 avio](#alluxio-%E6%89%B9%E9%87%8F%E5%8A%A0%E8%BD%BD%E5%B7%A5%E5%85%B7-avio)
      - [生成 avio 二进制包](#%E7%94%9F%E6%88%90-avio-%E4%BA%8C%E8%BF%9B%E5%88%B6%E5%8C%85)
      - [发布 avio 新版本](#%E5%8F%91%E5%B8%83-avio-%E6%96%B0%E7%89%88%E6%9C%AC)
  - [帮助](#%E5%B8%AE%E5%8A%A9)

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

目前 alluxio 服务部署在 jq 机房 jq13 ~ 17 和 jq19 ~ 21 这几台机器，以下约定这几台机器为 `alluxio 机器组`。 alluxio 的 web dashboard 通过 jq-ava 中的 kubernetes 转发以提供外网服务。

**部署前，请先确保你已经申请了 `alluxio 机器组` 的 root 权限。*

### 前置条件

请先确保在部署的机器上已经安装了 git，并将本代码库 clone 到目标机器。

1. 使用 root 账号登录服务器

2. 执行如下命令以安装 git。

```shell
apt update && apt install git
```

3. 克隆本代码库到本机指定位置，请注意克隆前需要将本机的公钥添加到 deploy key 中。

```shell
mkdir -p /disk1/repos/
git clone git@github.com:qiniu-ava/ava-alluxio.git /disk1/repos/ava-alluxio
```

4. 执行本代码库中 `deploy/env/install.sh` 脚本安装其他必要的依赖并做设置。

### 部署 zookeeper

目前 zookeeper 集群部署在 jq13、jq14、jq15 三个节点。分别在前述三个节点中执行如下步骤部署或者升级 zookeeper:

1. 使用 root 账号登录服务器

2. 进入 deploy 相关目录，并更新最新代码

```shell
cd /disk1/repos/ava-alluxio/deploy/docker
git pull
```

3. 执行 zookeeper 的部署脚本

```shell
./alluxio.zookeeper.sh
```

### 部署 master

在确保上述 zookeeper 部署成功后，可在 `alluxio 机器组` 任意三个以上(最好是奇数)节点上部署 master 实例。分别在需要部署 master 实例的服务器上执行如下步骤:

1. 同部署 zookeeper 中的 *1*。

2. 同部署 zookeeper 中的 *2*。

3. 执行 master 的部署脚本或者更新脚本

```shell
./alluxio.master.sh start/restart
```

### 部署 worker

在确保上述 zookeeper 部署成功后，可在 `alluxio 机器组` 任意节点上部署 worker 实例。分别在需要部署 worker 实例的服务器上执行如下步骤:

1. 同部署 zookeeper 中的 *1*

2. 同部署 zookeeper 中的 *2*

3. 在 jq13 ~ 17 上，执行 worker 的部署脚本或者更新脚本

```shell
./alluxio.worker.sh start/restart
```

4. 在 jq19 ~ 21 上，执行读 worker 和写 worker 的部署脚本或者更新脚本

```shell
./alluxio.worker.fat.sh start/restart
./alluxio.worker.write.sh start/restart
```

5. 依照[下述方式](#%E7%94%9F%E6%88%90-alluxio-%E5%8C%85)生成 alluxio 包，创建 Jira issue 给 Kirk 组相关同学帮忙更新 k8s 集群中各节点上的 alluxio worker 实例。

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
