<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [avio](#avio)
  - [概述](#%E6%A6%82%E8%BF%B0)
  - [相关概念](#%E7%9B%B8%E5%85%B3%E6%A6%82%E5%BF%B5)
    - [alluxio-fuse 文件系统](#alluxio-fuse-%E6%96%87%E4%BB%B6%E7%B3%BB%E7%BB%9F)
  - [设计](#%E8%AE%BE%E8%AE%A1)
    - [使用示例](#%E4%BD%BF%E7%94%A8%E7%A4%BA%E4%BE%8B)
    - [子命令列表](#%E5%AD%90%E5%91%BD%E4%BB%A4%E5%88%97%E8%A1%A8)
    - [preload](#preload)
    - [ls](#ls)
    - [stat](#stat)
    - [mv](#mv)
    - [cp](#cp)
    - [rm](#rm)
  - [帮助](#%E5%B8%AE%E5%8A%A9)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# avio

## 概述

`avio` 是为方便 `ava` 用户在『[工作台](http://portal.qiniu.com/ava/workspaces)』和『[训练](http://portal.qiniu.com/ava/trains)』中使用挂载的 bucket 中的文件而开发的命令行工具，其接口类似 `linux` 系统中基础的文件和目录管理命令行工具，并通过提高并发效率来提供比直接使用这些命令更好的性能。

## 相关概念

### alluxio-fuse 文件系统

在 `ava` 用户启动的 `pod` 中，可以在 `/workspace/mnt/buckets/` 目录下挂载一些 bucket。在这些 `bucket` 的目录下，用户可以像访问本地文件系统一样访问 `bucket` 中的文件。这是因为我们在相应的路径挂载了 `alluxio-fuse` 的文件系统，`alluxio-fuse` 文件系统底层基于分布式缓存系统 `alluxio` 和用户文件系统挂载工具 `fuse`。与本地文件系统不同的是，`alluxio-fuse` 文件系统中的原信息和文件数据实际都存储在远端的分布式存储系统中，访问时会有大量的网络交互，因此执行一些命令行如 mv/ls/cp 等相比在本地文件系统上会感觉到延时。


## 设计

### 使用示例

``` shell
avio <cmd> [options] [options]...
```

### 子命令列表

| cmd | description | status |
| :--- | :--- | :---: |
| help | 查看帮助文档 | WIP |
| preload | 将指定目录或者文件预加载到 alluxio 中，便于之后快速读取 | WIP |
| ls | 列出指定目录或者文件的信息 | WIP |
| stat | 查看指定目录或者文件在 alluxio 系统中的状况，比如是否已经加载缓存中，若已经加载进来位于哪一层级 | WIP |
| mv | 将指定目录或者文件移动到/出 alluxio 中 | WIP |
| cp | 将指定目录或者文件复制到/出 alluxio 中 | WIP |
| rm | 删除指定目录或者文件 | WIP |


### preload

avio preload [options] [path]
  + options:
    + -d --depth 递归深度，默认值为 4，最大值为 10
    + -i --input-file 所有需要 preload 的文件的列表文件路径
    + --is-jsonlist 只在 -i/--input-file 被设置时有效，可选项：true/false ，默认值为 false
    + -p --pool 并行 preload 的并发数，默认值为 50，最大值为 200
    + -l --log 当前任务的日志文件，默认值为 /var/log/avio/preload-<datetime>-<pid>.log
  + path:
    + path 和 -i/--input-file 参数是互斥的，若两者都没有设置则将会 preload 当前目录

**在任务启动之前会尝试查看找到的第一个文件所在的文件系统的 mountpoint，若发现此 mountpoint 不是 alluxio-fuse 类型的文件系统，将退出程序(且退出码非 0 )。*

### ls
avio ls [options] [path]
  + options
    + -f --full 当指定的路径下文件或子目录数量超过5000时，是否全量拉取，默认值为 false
    + -h 以易读的方式显示文件大小
    + -r 递归的列出目录下子目录中的文件
    + -c 只显示目录下的文件或子目录数

### stat

avio stat [options] [path]
  + options:
    + -d --depth 递归深度，默认值为递归到 4，最大值为 10
    + -t --show-tier 显示在缓存中的层级
    + --report 生成统计数据到指定文件，默认为 $HOME/.avio/report/stat_<start_unix_timestamp>.log
  + path:
    + 要统计的目录或者文件

### mv
avio mv [options] <source> <target>  
  + options:
    + -p --pool 并行 mv 的并发数，默认值为 50，最大值为 200
    + -l --log 当前任务的日志文件，默认值为 /var/log/avio/mv-<datetime>-<pid>.log
  + source:
    + 源地址
  + target:
    + 目的地址

**目前能支持的通配符只有\*，且 source 和 target 中至少需要有一个是在 alluxio-fuse 文件系统中*

### cp

avio cp [options] <source> <target>
  + options:
    + -p --pool 并行 cp 的并发数，默认值为 50，最大值为 200
    + -l --log 当前任务的日志文件，默认值为 /var/log/avio/cp-<datetime>-<pid>.log
  + source:
    + 源地址
  + target:
    + 目的地址

**目前能支持的通配符只有\*，且 source 和 target 中至少需要有一个是在 alluxio-fuse 文件系统中*

### rm
avio rm <path>
  + path:
    + 要删除的目录或者文件

**在任务启动之前会尝试查看找到的第一个文件所在的文件系统的 mountpoint，若发现此 mountpoint 不是 alluxio-fuse 类型的文件系统，将退出程序(且退出码非 0 )。*

## 帮助

1. 查看指定目录是否是在 alluxio-fuse 文件系统中
``` shell
df -h <path>
```

2. 局限性

由于 `alluxio-fuse` 文件系统中的文件树们是将 `kodo` 中的对象使用 `/` 作为分隔符来模拟出来符合 `linux` 命名规范的，而 `kodo` 中允许任意的 `/` 符号顺序和重复次数，而在 `linux` 文件系统中不允许出现文件名为空以及未转义的 `/` 符号，这使得 `alluxio-fuse` 文件系统在一些特殊路径下不可用。如：
  + `kodo` 中名为 `a/b//c` 的文件在我们的文件系统中无法访问
  + 所挂载的 `bucket` 中名称以 `/` 开头的对象都无法访问
  + `kodo` 中允许同时存在 `a/b` 和 `a/b/` 两个对象，这将导致 `alluxio-fuse` 文件系统中的冲突，目前出现此情况时我们会将 `a/b` 路径优先当为之前已经访问过的类型(是文件还是文件夹)，若之前没有访问过则将认为是文件。
