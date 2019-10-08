from ansible_api.tasks import run_im_adhoc
from time import sleep


def gather_host_info(host, retry=1):
    hosts = [host.__dict__]
    adhoc_success = False
    attempts = 0
    while attempts < retry and not adhoc_success:
        print("gather_host_info，times %d" % (attempts+1))
        try:
            result = run_im_adhoc(adhoc_data={'pattern': host.name, 'module': 'setup'},
                                  inventory_data={'hosts': hosts, 'vars': {}})
            success = is_adhoc_success(result)
            if success:
                return result["raw"]["ok"][host.name]["setup"]["ansible_facts"]
            attempts += 1
            sleep(1)
        except Exception:
            attempts += 1
            if attempts == retry:
                raise Exception("get os info failed!")
            else:
                sleep(1)

def get_cluster_token(host):
    hosts = [host.__dict__]
    shell = "kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep admin-user | awk '{print $1}') | grep token: | awk '{print $2}'"
    result = run_im_adhoc(adhoc_data={'pattern': host.name, 'module': 'shell', 'args': shell},
                          inventory_data={'hosts': hosts, 'vars': {}})
    return result.get('raw').get('ok')[host.name]['command']['stdout']


def fetch_cluster_config(host, dest):
    hosts = [host.__dict__]
    args = {'src': '/root/.kube/config',
            'dest': dest}
    result = run_im_adhoc(adhoc_data={'pattern': host.name,
                                      'module': 'fetch',
                                      'args': args},
                          inventory_data={'hosts': hosts, 'vars': {}})
    if not is_adhoc_success(result):
        raise Exception("get cluster config failed!")
    return result['raw']['ok'][host.name]['fetch']['dest']


def storage_health_check(host, module, command):
    hosts = [host.__dict__]
    result = run_im_adhoc(adhoc_data={'pattern': host.name, 'module': module, 'args': command},
                          inventory_data={"hosts": hosts, "vars": {}})
    return is_adhoc_success(result)


def is_adhoc_success(result):
    return result.get('summary', {}).get('success', False)
