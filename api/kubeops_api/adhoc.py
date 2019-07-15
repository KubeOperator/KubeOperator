from ansible_api.tasks import run_im_adhoc


def gather_host_info(host):
    hosts = [host.__dict__]
    result = run_im_adhoc(adhoc_data={'pattern': host.name, 'module': 'setup'},
                          inventory_data={'hosts': hosts, 'vars': {}})
    if not is_adhoc_success(result):
        raise Exception("get os info failed!")
    return result["raw"]["ok"][host.name]["setup"]["ansible_facts"]


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
    print(result)
    return result['raw']['ok'][host.name]['fetch']['dest']


def is_adhoc_success(result):
    return result.get('summary', {}).get('success', False)
