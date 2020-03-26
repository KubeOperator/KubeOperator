from ansible_api.tasks import run_im_adhoc
from time import sleep


def drain_worker_node(host, worker_name):
    hosts = [host.__dict__]
    shell = "kubectl drain {} --delete-local-data --force --ignore-daemonsets && kubectl delete node {}".format(
        worker_name, worker_name)
    result = run_im_adhoc(adhoc_data={'pattern': host.name, 'module': 'shell', 'args': shell},
                          inventory_data={'hosts': hosts, 'vars': {}})
    if not is_adhoc_success(result):
        raise Exception("drain! node {} failed!".format(worker_name))


def gather_host_info(ip, port, username, retry=1, **kwargs):
    hosts = [{
        "ip": ip,
        "username": username,
        "password": kwargs.get('password', None),
        "private_key_path": kwargs.get('private_key_path', None),
        "name": "default",
        "port": port
    }]
    attempts = 0
    while attempts < retry:
        result = run_im_adhoc(adhoc_data={'pattern': "default", 'module': 'setup'},
                              inventory_data={'hosts': hosts, 'vars': {}})
        if is_adhoc_success(result):
            return result["raw"]["ok"]["default"]["setup"]["ansible_facts"]
        sleep(5)
        attempts += 1
        if attempts == retry:
            raise Exception('get os info fail')


def test_host(ip, port, username, **kwargs):
    r = False
    hosts = [{
        "ip": ip,
        "username": username,
        "password": kwargs.get('password', None),
        "private_key_path": kwargs.get('private_key_path', None),
        "name": "default",
        "port": port
    }]
    result = run_im_adhoc(adhoc_data={'pattern': "default", 'module': 'ping'},
                          inventory_data={'hosts': hosts, 'vars': {}})
    if is_adhoc_success(result):
        r = True
    return r


def get_cluster_token(host):
    hosts = [host.__dict__]
    shell = "kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep tiller | awk '{print $1}') | grep token: | awk '{print $2}'"
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


def is_adhoc_success(result):
    return result.get('summary', {}).get('success', False)


def get_host_time(ip, port, username, hostname, **kwargs, ):
    hosts = [{
        "ip": ip,
        "username": username,
        "password": kwargs.get('password', None),
        "private_key_path": kwargs.get('private_key_path', None),
        "name": hostname,
        "port": port
    }]
    shell = 'date'
    result = run_im_adhoc(adhoc_data={'pattern': hostname, 'module': 'shell', 'args': shell},
                          inventory_data={'hosts': hosts, 'vars': {}})
    if is_adhoc_success(result):
        return result["raw"]["ok"][hostname]["command"]["stdout_lines"][0]
