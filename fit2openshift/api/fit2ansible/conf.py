import os

import yaml

def load_env():
    try:
        f = open("/opt/fit2openshift/conf/fit2openshift/fit2openshift_conf.yml")
        res = yaml.load(f)
    except Exception:
        print("配置文件错误!")
        raise Exception
    return res


