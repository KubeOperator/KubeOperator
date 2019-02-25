# pylint: skip-file
# flake8: noqa


def main():
    '''
    ansible oc module for managing OpenShift configmap objects
    '''

    module = AnsibleModule(
        argument_spec=dict(
            kubeconfig=dict(default='/etc/origin/master/admin.kubeconfig', type='str'),
            state=dict(default='present', type='str',
                       choices=['present', 'absent', 'list']),
            debug=dict(default=False, type='bool'),
            namespace=dict(default='default', type='str'),
            name=dict(default=None, required=True, type='str'),
            from_file=dict(default=None, type='dict'),
            from_literal=dict(default=None, type='dict'),
        ),
        supports_check_mode=True,
    )


    rval = OCConfigMap.run_ansible(module.params, module.check_mode)
    if 'failed' in rval:
        module.fail_json(**rval)

    module.exit_json(**rval)

if __name__ == '__main__':
    main()
