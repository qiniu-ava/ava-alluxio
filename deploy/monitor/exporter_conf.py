import yaml
import sys
reload(sys)
sys.setdefaultencoding('utf-8')
import urllib
import urllib2
import json
import yaml
import os
def main():
    with open('./alluxio-exporter.yml','r') as f:
        data = yaml.load(f)
    content = {}
    content["alluxio"]=[]
    for ele in data:
        group_name = ele["group"]
        master_host = ele["master_host"]
        if not os.path.exists(os.environ['HOME'] + "/alluxio-exporter"):
            os.makedirs(os.environ['HOME'] + "/alluxio-exporter")
        try:
            url = "http://" + master_host + "/api/v1/master/info"
            req = urllib2.Request(url)
            res_data = urllib2.urlopen(req)
            res = json.loads(res_data.read().encode("utf-8"))
        except Exception as e:
            print "urlopen error: ", e
        worker_host=[]
        for worker in res["workers"]:
            worker_host.append(worker["address"]["host"] + ":" + str(worker["address"]["webPort"]))
        content["alluxio"].append({"name": group_name, "master_host": master_host,"worker_host" : worker_host})
    with open(os.environ['HOME'] + "/alluxio-exporter/" + "exporter.yml", "w") as f:
        yaml.safe_dump(content, f, default_flow_style=False, encoding='utf-8', allow_unicode=True)
    f.close()

if __name__ == '__main__':
    main()
