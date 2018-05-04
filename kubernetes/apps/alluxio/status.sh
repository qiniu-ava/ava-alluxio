
status_master() {
  echo "checking status of alluxio master"
  kubectl get pods | grep alluxio | grep master | grep bowen
}

status_worker() {
  echo "checking status of alluxio worker"
  kubectl get pods | grep alluxio | grep worker | grep bowen
}

status_client() {
  echo "checking status of alluxio client"
  kubectl get pods | grep alluxio | grep client | grep bowen
}

status_zookeeper() {
  echo "checking status of alluxio zookeeper"
  kubectl get pods | grep alluxio | grep zk
}

print_usage() {
  echo "usage ./status [<all/master/worker/zookeeper/zk>]"
}

app=$1

if [ "$app" = "" ]; then 
  app=all
fi

case ${app} in
  all)
    status_master
    echo -e "\n\n"
    status_worker
    echo -e "\n\n"
    status_client
    echo -e "\n\n"
    status_zookeeper
    ;;
  master)
    status_master
    ;;
  worker)
    status_worker
    ;;
  client)
    status_client
    ;;
  zookeeper)
    status_zookeeper
    ;;
  zk)
    status_zookeeper
    ;;
  *)
    print_usage;
    ;;
esac

