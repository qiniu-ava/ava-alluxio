import sys
reload(sys)
sys.setdefaultencoding('utf-8')
import urllib
import urllib2
import json
import yaml
import os


def main():
    group_name = sys.argv[1]
    master_host = sys.argv[2]
    if not os.path.exists(os.environ['HOME'] + "/alluxio-exporter"):
        os.makedirs(os.environ['HOME'] + "/alluxio-exporter")
    try:
        url = "http://" + master_host + "/api/v1/master/info"
        req = urllib2.Request(url)
        res_data = urllib2.urlopen(req)
        res = json.loads(res_data.read().encode("utf-8"))
    except Exception as e:
        print "urlopen error: ", e
    content = {}
    content["alluxio"] = []
    content["alluxio"].append({"type": "master", "host": master_host})
    for worker in res["workers"]:
        content["alluxio"].append({"type": "worker",
            "host": worker["address"]["host"] + ":" + str(worker["address"]["webPort"])})
    with open(os.environ['HOME'] + "/alluxio-exporter/" + group_name + "-exporter.yml", "w") as f:
        yaml.safe_dump(content, f, default_flow_style=False, encoding='utf-8', allow_unicode=True)
    f.close()


if __name__ == '__main__':
    if len(sys.argv) != 3:
        print "usage:python exporter_conf.py <group_name> <cluster_master_ingress or master_host:webport>"
        sys.exit()
    main()