from ansible_api.tasks import run_im_adhoc


def gather_host_info(ip, username, password):
    hosts = [{
        "ip": ip,
        "username": username,
        "password": password,
        "name": "default"
    }]
    result = run_im_adhoc(adhoc_data={'pattern': "default", 'module': 'setup'},
                          inventory_data={'hosts': hosts, 'vars': {}})
    if not is_adhoc_success(result):
        raise Exception("get os info failed!")
    return result["raw"]["ok"]["default"]["setup"]["ansible_facts"]


def test_host(ip, username, password):
    r = False
    hosts = [{
        "ip": ip,
        "username": username,
        "password": password,
        "name": "default"
    }]
    result = run_im_adhoc(adhoc_data={'pattern': "default", 'module': 'ping'},
                          inventory_data={'hosts': hosts, 'vars': {}})
    if is_adhoc_success(result):
        r = True
    return r


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


def is_adhoc_success(result):
    return result.get('summary', {}).get('success', False)
