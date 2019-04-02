# pylint: skip-file
# flake8: noqa

def main():
    '''
    ansible oc module for user
    '''

    module = AnsibleModule(
        argument_spec=dict(
            kubeconfig=dict(default='/etc/origin/master/admin.kubeconfig', type='str'),
            state=dict(default='present', type='str',
                       choices=['present', 'absent', 'list']),
            debug=dict(default=False, type='bool'),
            username=dict(default=None, type='str'),
            full_name=dict(default=None, type='str'),
            # setting groups for user data will not populate the
            # 'groups' field in the user data.
            # it will call out to the group data and make the user
            # entry there
            groups=dict(default=[], type='list'),
        ),
        supports_check_mode=True,
    )

    results = OCUser.run_ansible(module.params, module.check_mode)

    if 'failed' in results:
        module.fail_json(**results)

    module.exit_json(**results)

if __name__ == '__main__':
    main()
