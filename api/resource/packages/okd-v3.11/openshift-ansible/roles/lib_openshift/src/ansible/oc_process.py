# pylint: skip-file
# flake8: noqa


def main():
    '''
    ansible oc module for processing templates
    '''

    module = AnsibleModule(
        argument_spec=dict(
            kubeconfig=dict(default='/etc/origin/master/admin.kubeconfig', type='str'),
            state=dict(default='present', type='str', choices=['present', 'list']),
            debug=dict(default=False, type='bool'),
            namespace=dict(default='default', type='str'),
            template_name=dict(default=None, type='str'),
            content=dict(default=None, type='str'),
            params=dict(default=None, type='dict'),
            create=dict(default=False, type='bool'),
            reconcile=dict(default=True, type='bool'),
        ),
        supports_check_mode=True,
    )

    rval = OCProcess.run_ansible(module.params, module.check_mode)
    if 'failed' in rval:
        module.fail_json(**rval)

    module.exit_json(**rval)

if __name__ == '__main__':
    main()
