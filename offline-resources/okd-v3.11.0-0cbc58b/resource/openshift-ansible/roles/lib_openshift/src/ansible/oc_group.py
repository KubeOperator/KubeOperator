# pylint: skip-file
# flake8: noqa

#pylint: disable=too-many-branches
def main():
    '''
    ansible oc module for group
    '''

    module = AnsibleModule(
        argument_spec=dict(
            kubeconfig=dict(default='/etc/origin/master/admin.kubeconfig', type='str'),
            state=dict(default='present', type='str',
                       choices=['present', 'absent', 'list']),
            debug=dict(default=False, type='bool'),
            name=dict(default=None, type='str'),
            namespace=dict(default='default', type='str'),
            # addind users to a group is handled through the oc_users module
            #users=dict(default=None, type='list'),
        ),
        supports_check_mode=True,
    )

    rval = OCGroup.run_ansible(module.params, module.check_mode)

    if 'failed' in rval:
        return module.fail_json(**rval)

    return module.exit_json(**rval)

if __name__ == '__main__':
    main()
