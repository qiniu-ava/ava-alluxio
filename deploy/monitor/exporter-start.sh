#!/bin/bash

cd "$( dirname "${BASH_SOURCE[0]}" )" || exit
exporter=$1
cmd=$2
group=$3
master_host=$4

main_start() {
    case $exporter in
      all)
        if [ "$group" == "" ] || [ "$master_host" == "" ] ; then
          echo "usage: ./exporter-start.sh start all <group_name> <cluster_master_ingress or master_host:webport>"
          return
        fi
        /bin/bash alluxio-exporter.sh start "$group" "$master_host"
        /bin/bash cadvisor.sh start
        /bin/bash node-exporter.sh start
        /bin/bash jvm-exporter.sh start
        ;;
      alluxio-exporter)
        if [ "$group" == "" ] || [ "$master_host" == "" ] ; then
          echo "usage: ./exporter-start.sh start alluxio_exporter <group_name> <cluster_master_ingress or master_host:webport>"
          return
        fi
        /bin/bash alluxio-exporter.sh start "$group" "$master_host"
        ;;
      node-exporter)
        /bin/bash node-exporter.sh start
        ;;
      jvm-exporter)
        /bin/bash jvm-exporter.sh start
        ;;
      cadvisor)
        /bin/bash cadvisor.sh start
        ;;
    esac
}

remove() {
    case $exporter in
      all)
        /bin/bash alluxio-exporter.sh remove "$group"
        /bin/bash cadvisor.sh remove
        /bin/bash node-exporter.sh remove
        /bin/bash jvm-exporter.sh remove
        ;;
      alluxio_exporter)
        /bin/bash alluxio-exporter.sh remove
        ;;
      node-exporter)
        /bin/bash node-exporter.sh remove
        ;;
      jvm-exporter)
        /bin/bash jvm-exporter.sh remove
        ;;
      cadvisor)
        /bin/bash cadvisor.sh remove
        ;;
    esac
}

case $cmd in
  start)
    main_start
  ;;
  restart)
    remove
    main_start
  ;;
  remove)
    remove
  ;;
  rm)
    remove
  ;;
esac
