#! /usr/bin/python3
"""
apt-get update
apt-get install -y python3-pip
pip3 install google-apputils
pip3 install protobuf
protoc --proto_path=. --python_out=. AsyncCache.proto
"""

from AsyncCache_pb2 import AsyncCacheRequest, LocalBlockOpenResponse
import sys
import socket
import struct
import time

if len(sys.argv) != 4:
    print("Usage:", sys.argv[0], "ip port file_list")
    sys.exit(-1)

ip = sys.argv[1]
port = int(sys.argv[2])
flist = sys.argv[3]

address = (ip, port)
clientsocket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
clientsocket.connect(address)

TH_HI = 100
TH_LOW = 50
ID_PURGE_BLOCK = 0
ID_QUERY_BLOCK = -1


def async_cache(id, length):
    data = AsyncCacheRequest()
    data.length = length
    data.block_id = id
    s = data.SerializeToString()
    # (long)length - (int)id - (int)packet_length - packet
    clientsocket.sendall(struct.pack('>QII', 16 + len(s), 112, len(s)) + s)


def evict(id):
    async_cache(id, ID_PURGE_BLOCK)


def trim(idlist, low):
    rsp = LocalBlockOpenResponse()
    while len(idlist) > low:
        not_done = []
        for id in idlist:
            async_cache(id, ID_QUERY_BLOCK)
            result = clientsocket.recv(16)
            if result == b'':
                raise RuntimeError("socket connection broken")
            alen, retid, plen = struct.unpack('>QII', result)
            if retid != 106:
                raise RuntimeError("not expect id")
            result = clientsocket.recv(plen)
            rsp.ParseFromString(result)
            if rsp.path[0] != '/':
                print("done ", id)
            else:
                not_done.append(id)
                evict(id)
        idlist = not_done
        time.sleep(2)
    return idlist


idlist = []
with open(flist, "r") as lines:
    for line in lines:
        id = int(line.split()[0])
        if len(idlist) >= TH_HI:
            idlist = trim(idlist, TH_LOW)
        print("evicting ", id)
        evict(id)
        idlist.append(id)

trim(idlist, 0)
