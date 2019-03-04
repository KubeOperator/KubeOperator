import os

import yaml

from openshift_api.models import Setting


def load_env():
    try:
        f = open("/opt/fit2openshift/conf/fit2openshift/fit2openshift_conf.yml")
        res = yaml.load(f)
    except Exception:
        print("配置文件错误!")
        raise Exception
    return res


def set_host():
    hostname = Setting.objects.filter(key='hostname').first()
    if hostname:
        os.putenv("REGISTORY_HOSTNAME", hostname.value)
