# pylint: skip-file
# flake8: noqa

def main():
    '''
    ansible oc module for project
    '''

    module = AnsibleModule(
        argument_spec=dict(
            kubeconfig=dict(default='/etc/origin/master/admin.kubeconfig', type='str'),
            state=dict(default='present', type='str',
                       choices=['present', 'absent', 'list']),
            debug=dict(default=False, type='bool'),
            name=dict(default=None, require=True, type='str'),
            display_name=dict(default=None, type='str'),
            node_selector=dict(default=None, type='list'),
            description=dict(default=None, type='str'),
            admin=dict(default=None, type='str'),
            admin_role=dict(default='admin', type='str'),
        ),
        supports_check_mode=True,
    )

    rval = OCProject.run_ansible(module.params, module.check_mode)
    if 'failed' in rval:
        return module.fail_json(**rval)

    return module.exit_json(**rval)


if __name__ == '__main__':
    main()
