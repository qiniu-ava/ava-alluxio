# -*- coding:utf-8 -*-
'''
usage: python regression.py <ak> <sk> <bucket> <container_name>
Preconditions:
1. Known ak/sk of avatest account
2. In the bucket of the avatest account, the folder_1w folder is already included,
and there are more than 5000 files in it.
3. In the bucket of the avatest account, there are already multiple files and folders.
4. The following operations are performed under alluxio-proxy. If the program is running
in a container, the container needs priviledge permission when it starts.
'''
import sys
import subprocess
import os

# execute the command
def execute(title, cmd, i):
    print "------------"+str(i)+"."+title+"------------"
    child = subprocess.Popen([cmd], shell=True)
    child.wait()
    if child.returncode == 0:
        print "Success to "+title+" !"
    else:
        print "Fail to "+title+" !"

# execute the command without output and just get the returncode
def execute_ignore_output(title, cmd, i):
    print "------------"+str(i)+"."+title+"------------"
    dev_null = open(os.devnull, 'w')
    returncode = subprocess.call([cmd], shell=True, stdout=dev_null)
    if returncode == 0:
        print "Success to "+title+" !"
    else:
        print "Fail to "+title+" !"

# state the command and its title. execute command one by one
def main():
    docker = "docker exec "+sys.argv[4]+" "
    title0 = "create path in alluxio"
    title1 = "create the same path again"
    title2 = "mount bucket-1"
    title3 = "mount bucket-2"
    title4 = "mount bucket to the same path"
    title5 = "mount bucket to the nonexistent path"
    title6 = "different users mount the same bucket"
    title7 = "unmount mountpoint"
    title8 = "unmount nonexistent mountpoint"
    title9 = "fuse bucket to the local path"
    title10 = "fuse unmount local file path"
    title11 = "list files"
    title12 = "list nonexistent files"
    title13 = "write to file"
    title14 = "cat files in kodo"
    title15 = "change the file name"
    title16 = "remove file"
    title17 = "create large file in alluxio"
    title18 = "read large file in alluxio"
    title19 = "qrsctl put file and list file"
    title20 = "ls folder with more than 5000 files"

    cmd0 = docker+"/opt/alluxio/bin/alluxio fs mkdir /a/b/1381351869 /a/b/1381102889"
    cmd1 = docker+"/opt/alluxio/bin/alluxio fs mkdir /a/b/1381351869"
    cmd2 = docker+"/opt/alluxio/bin/alluxio fs mount \
        --option fs.oss.accessKeyId=wb3L0e4GBOf_Kq4IVS8y9Csq9fC3u8UDmqb2S-pj \
        --option fs.oss.accessKeySecret=XWo3jAwdW7ETQnJsnWTIRedrmEc6Au-8jt2Xj9KV \
        --option fs.oss.userId=1381351869 \
        --option fs.oss.endpoint=petvhg8di.bkt.clouddn.com \
        /a/b/1381351869/alluxio-test oss://"+sys.argv[3]
    cmd3 = docker+"/opt/alluxio/bin/alluxio fs mount \
        --option fs.oss.accessKeyId=wb3L0e4GBOf_Kq4IVS8y9Csq9fC3u8UDmqb2S-pj \
        --option fs.oss.accessKeySecret=XWo3jAwdW7ETQnJsnWTIRedrmEc6Au-8jt2Xj9KV \
        --option fs.oss.userId=1381102889 \
        --option fs.oss.endpoint=petvhg8di.bkt.clouddn.com \
        /a/b/1381102889/alluxio-test oss://"+sys.argv[3]
    cmd4 = docker+"/opt/alluxio/bin/alluxio fs mount \
        --option fs.oss.accessKeyId=wb3L0e4GBOf_Kq4IVS8y9Csq9fC3u8UDmqb2S-pj \
        --option fs.oss.accessKeySecret=XWo3jAwdW7ETQnJsnWTIRedrmEc6Au-8jt2Xj9KV \
        --option fs.oss.userId=1381102889 \
        --option fs.oss.endpoint=petvhg8di.bkt.clouddn.com \
        /a/b/1381102889/alluxio-test oss://"+sys.argv[3]
    cmd5 = docker+"/opt/alluxio/bin/alluxio fs mount \
        --option fs.oss.accessKeyId=wb3L0e4GBOf_Kq4IVS8y9Csq9fC3u8UDmqb2S-pj \
        --option fs.oss.accessKeySecret=XWo3jAwdW7ETQnJsnWTIRedrmEc6Au-8jt2Xj9KV \
        --option fs.oss.userId=1381102889 \
        --option fs.oss.endpoint=petvhg8di.bkt.clouddn.com \
        /a/b/not_exists_path/alluxio-test oss://"+sys.argv[3]
    cmd6 = docker+"/opt/alluxio/bin/alluxio fs mount \
        --option fs.oss.accessKeyId=wb3L0e4GBOf_Kq4IVS8y9Csq9fC3u8UDmqb2S-pj \
        --option fs.oss.accessKeySecret=XWo3jAwdW7ETQnJsnWTIRedrmEc6Au-8jt2Xj9KV \
        --option fs.oss.userId=1381102889 \
        --option fs.oss.endpoint=petvhg8di.bkt.clouddn.com \
        /a/b/1381351869/alluxio-test oss://"+sys.argv[3]
    cmd7 = docker+"/opt/alluxio/bin/alluxio fs unmount /a/b/1381351869/alluxio-test"
    cmd8 = docker+"/opt/alluxio/bin/alluxio fs unmount /a/b/not_exists_path/alluxio-test"
    cmd9 = docker+"mkdir -p /test-alluxio-fuse" \
        +"&&"+docker+"/opt/alluxio/integration/fuse/bin/alluxio-fuse \
        mount /test-alluxio-fuse /a/b/1381102889/alluxio-test"
    cmd10 = docker+"/opt/alluxio/integration/fuse/bin/alluxio-fuse umount /test-alluxio-fuse" \
        +"&&"+docker+ "/opt/alluxio/integration/fuse/bin/alluxio-fuse \
        mount /test-alluxio-fuse /a/b/1381102889/alluxio-test "
    cmd11 = docker+"/opt/alluxio/bin/alluxio fs ls /a/b/1381102889/alluxio-test"
    cmd12 = docker+"/opt/alluxio/bin/alluxio fs ls /a/b/1381102889/alluxio-test/path_not_exists"
    cmd13 = docker+'bash -c "echo test-alluxio-fuse-write > /test-alluxio-fuse/test-write"'
    cmd14 = docker+"/opt/alluxio/bin/alluxio fs cat /a/b/1381102889/alluxio-test/test-write"
    cmd15 = docker+"/opt/alluxio/bin/alluxio fs mv /a/b/1381102889/alluxio-test/test-write \
        /a/b/1381102889/alluxio-test/write-test"+"&&"+docker+"/opt/alluxio/bin/alluxio fs mv \
        /a/b/1381102889/alluxio-test/write-test /a/b/1381102889/alluxio-test/test-write"
    cmd16 = docker+"rm /test-alluxio-fuse/test-write"
    cmd17 = docker+'dd if=/dev/urandom of=/test-alluxio-fuse/test-write-2G bs=4M count=512'
    cmd18 = docker+'dd if=/test-alluxio-fuse/test-write-2G of=/dev/null bs=4M'
    cmd19 = docker+"curl -O http://devtools.qiniu.com/linux/amd64/qrsctl?ref=developer.qiniu.com \
        -o qrsctl" \
        +"&&"+docker+"chmod +x qrsctl" \
        +"&&"+docker+"mv qrsctl /usr/local/bin" \
        +"&&"+docker+"qrsctl login "+sys.argv[1]+" "+sys.argv[2] \
        +"&&"+docker+"mkdir -p /tmp/" \
        +"&&"+docker+"dd if=/dev/urandom of=/tmp/1m bs=1M count=1" \
        +"&&"+docker+"qrsctl put "+sys.argv[3]+" test-ls/1 /tmp/1m" \
        +"&&"+docker+'echo "first time"' \
        +"&&"+docker+"ls -lh /test-alluxio-fuse/test-ls/1" \
        +"&&"+docker+"dd if=/dev/urandom of=/tmp/4m bs=1M count=4 " \
        +"&&"+docker+"qrsctl put "+sys.argv[3]+" test-ls/1 /tmp/4m" \
        +"&&"+docker+'echo "second time"' \
        +"&&"+docker+"ls -lh /test-alluxio-fuse/test-ls/1"
    cmd20 = docker+'ls -l /test-alluxio-fuse/folder_1w'

    title = [title0, title1, title2, title3, title4, title5, title6, title7, title8, title9,\
        title10, title11, title12, title13, title14, title15, title16, title17, title18, title19,\
        title20]
    cmd = [cmd0, cmd1, cmd2, cmd3, cmd4, cmd5, cmd6, cmd7, cmd8, cmd9, cmd10, cmd11, cmd12, \
        cmd13, cmd14, cmd15, cmd16, cmd17, cmd18, cmd19, cmd20]

    for i in range(len(title)):
        if i == 0:
            print "*******************Mount相关*******************"
        if i == 11:
            print "*******************文件相关********************"
        if i == 11 or i == 20:
            execute_ignore_output(title[i], cmd[i], i)
            continue
        execute(title[i], cmd[i], i)

if __name__ == '__main__':
    if len(sys.argv) != 5:
        print "usage:python regression.py <ak> <sk> <bucket> <container_name>"
        sys.exit()
    main()
