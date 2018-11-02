# pylint: skip-file
# flake8: noqa

def main():
    '''
    ansible oc module for services
    '''

    module = AnsibleModule(
        argument_spec=dict(
            kubeconfig=dict(default='/etc/origin/master/admin.kubeconfig', type='str'),
            state=dict(default='present', type='str',
                       choices=['present', 'absent', 'list']),
            debug=dict(default=False, type='bool'),
            namespace=dict(default='default', type='str'),
            name=dict(default=None, type='str'),
            annotations=dict(default=None, type='dict'),
            labels=dict(default=None, type='dict'),
            selector=dict(default=None, type='dict'),
            clusterip=dict(default=None, type='str'),
            portalip=dict(default=None, type='str'),
            ports=dict(default=None, type='list'),
            session_affinity=dict(default='None', type='str'),
            service_type=dict(default='ClusterIP', type='str'),
            external_ips=dict(default=None, type='list'),
        ),
        supports_check_mode=True,
    )

    rval = OCService.run_ansible(module.params, module.check_mode)
    if 'failed' in rval:
        return module.fail_json(**rval)

    return module.exit_json(**rval)


if __name__ == '__main__':
    main()
