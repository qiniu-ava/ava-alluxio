# alluxio 在 kubernetes 上部署的组件

## alluxio master dashboard

此服务用来查看 alluxio 运行状态，使用 [`qiniu-ava/ava-alluxio/docker/app/alluxio`](../../../docker/app/alluxio/) 编译出来的镜像，编译操作应使用[`qiniu-ava/ava-alluxio/docker-build.sh`](../../../dicker-build.sh) 来生成。

完整发布流程如下:

1. 修改 ./deploy/kubernetes/alluxio/dashboard.yml 中 alluxio-nginx 镜像的版本号为 2 中将会打上的 tag

2. 提交 commit 并基本当前 commit 打了 tag

3. 编译 alluxio-nginx 镜像

```shell
cd $PATH_TO_REPO
./docker-build.sh dashboard -p
```

4. 更新 dashboard 相关组件

```shell
cd $PATH_TO_REPO
kubectl apply -f ./deploy/kubernetes/alluxio/dashboard.yml
```

## alluxio proxy

此服务用来在 avio 中操作 alluxio 文件。使用线上 alluxio 镜像，发布流程如下:

```shell
kubectl apply -f ./deploy/kubernetes/alluxio/proxy.yml
```
