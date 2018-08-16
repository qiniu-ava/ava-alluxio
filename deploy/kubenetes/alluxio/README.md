# alluxio 在 kubernetes 上部署的组件

## alluxio master dashboard

此服务用来查看 alluxio 运行状态，使用 [`qiniu-ava/ava-alluxio/docker/app/alluxio`](../../../docker/app/alluxio/) 编译出来的镜像，编译操作应使用[`qiniu-ava/ava-alluxio/docker-build.sh`](../../../dicker-build.sh) 来生成。完整发布流程如下:

```shell

cd $PATH_TO_REPO

./docker-build.sh dashboard
kubectl create -f ./deploy/kubernetes/alluxio/dashboard.yml

```

## alluxio proxy

此服务用来在 avio 中操作 alluxio 文件。使用线上 alluxio 镜像，发布流程如下:

```shell

kubectl create -f ./deploy/kubernetes/alluxio/proxy.yml

```
