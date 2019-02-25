# pylint: skip-file
# flake8: noqa

def main():
    '''
    ansible oc module for service accounts
    '''

    module = AnsibleModule(
        argument_spec=dict(
            kubeconfig=dict(default='/etc/origin/master/admin.kubeconfig', type='str'),
            state=dict(default='present', type='str',
                       choices=['present', 'absent', 'list']),
            debug=dict(default=False, type='bool'),
            name=dict(default=None, required=True, type='str'),
            namespace=dict(default=None, required=True, type='str'),
            secrets=dict(default=None, type='list'),
            image_pull_secrets=dict(default=None, type='list'),
        ),
        supports_check_mode=True,
    )

    rval = OCServiceAccount.run_ansible(module.params, module.check_mode)
    if 'failed' in rval:
        module.fail_json(**rval)

    module.exit_json(**rval)

if __name__ == '__main__':
    main()
