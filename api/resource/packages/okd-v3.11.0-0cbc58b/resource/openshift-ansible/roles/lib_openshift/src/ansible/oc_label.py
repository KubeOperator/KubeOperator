# pylint: skip-file
# flake8: noqa

def main():
    ''' ansible oc module for labels '''

    module = AnsibleModule(
        argument_spec=dict(
            kubeconfig=dict(default='/etc/origin/master/admin.kubeconfig', type='str'),
            state=dict(default='present', type='str',
                       choices=['present', 'absent', 'list', 'add']),
            debug=dict(default=False, type='bool'),
            kind=dict(default='node', type='str',
                      choices=['node', 'pod', 'namespace']),
            name=dict(default=None, type='str'),
            namespace=dict(default=None, type='str'),
            labels=dict(default=None, type='list'),
            selector=dict(default=None, type='str'),
        ),
        supports_check_mode=True,
        mutually_exclusive=(['name', 'selector']),
    )

    results = OCLabel.run_ansible(module.params, module.check_mode)

    if 'failed' in results:
        module.fail_json(**results)

    module.exit_json(**results)

if __name__ == '__main__':
    main()
