import sys
reload(sys)
sys.setdefaultencoding('utf-8')
from pandora import api
import json


def json_psarser(r):
    rJsonStr = "total: " + str(r["total"]) + ", "
    datas = r["data"]
    dataNew = ""
    stopFlag = True if len(datas) > 0 else False
    if len(datas) > 0:
        for data in datas:
            jsonDic = "dockerID: " + data["dockerID"].encode('utf-8') + ", " + "hostname: " + data["hostname"].encode('utf-8') + ", " + "localip: " + data["localip"].encode('utf-8') + "\n" + data["raw"] + "\n"
            dataNew += jsonDic
    rJsonStr += dataNew
    return rJsonStr, r["scroll_id"], stopFlag


def search(repo, headers, scroll, client, offset=1, method="POST"):
    url = "/v5/repos/" + repo + "/search"
    print offset
    if scroll == "" :
        data = json.dumps({"size":10, "from": offset})
    else:
        data = json.dumps({"size":10, "from": offset, "scroll": scroll})
    response = client.do_request(method, url, data=data, headers=headers)
    r = json.loads(response.read())
    return json_psarser(r)


def scrollDownload(headers, repo, scroll_id, scroll, client, method="POST"):
    url = "/v5/repos/" + repo + "/scroll"
    data = json.dumps({"scroll_id": scroll_id})
    response = client.do_request(method, url, data=data, headers=headers)
    r = json.loads(response.read().encode('utf-8'))
    return json_psarser(r)


def read_log(fp, repo, headers, scroll, client, offset=1):
    r, scroll_id, stopFlag = search(repo, headers, scroll, client, offset)
    fp.write(r + "\n")
    count = 1
    while stopFlag:
        r, scroll_id, stopFlag = scrollDownload(headers, repo, scroll_id, scroll, client)
        fp.write(r + "\n")
        count += 1
        print count, scroll_id


def main():
    repo = raw_input("Enter your repo: ")
    ak = raw_input("Enter your pandora ak: ")
    sk = raw_input("Enter your pandora sk: ")
    offset = raw_input("Enter the offset(default: 1): ")
    try:
        offset = int(offset)
    except Exception as e:
        print "string2int error: ", e
        return "please input a number"
    scroll = raw_input("Enter your scroll time: ")
    if repo == "" or ak == "" or sk == "" :
        return "please retry with correct repo and aksk"
    endpoint = 'https://nb-insight.qiniuapi.com'
    client = api.Client(endpoint, ak, sk)
    # encodedSign = sign(method, contentType, resource, sk)
    headers = {"Content-Type": "application/json"}
    with open('alluxio-log.txt', 'w') as f:
        try:
            read_log(f, repo, headers, scroll, client, offset)
        except Exception as e:
            print "readLog Error: ", e
            res = "readLog failure"
        else:
            res = "readLog success"
        finally:
            f.close()
            return res
    

if __name__ == '__main__':
    res = main()
    print res
