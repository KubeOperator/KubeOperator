# pylint: skip-file
# flake8: noqa

def main():
    '''
    ansible oc module for environment variables
    '''

    module = AnsibleModule(
        argument_spec=dict(
            kubeconfig=dict(default='/etc/origin/master/admin.kubeconfig', type='str'),
            state=dict(default='present', type='str',
                       choices=['present', 'absent', 'list']),
            debug=dict(default=False, type='bool'),
            kind=dict(default='rc', choices=['dc', 'rc', 'pods'], type='str'),
            namespace=dict(default='default', type='str'),
            name=dict(default=None, required=True, type='str'),
            env_vars=dict(default=None, type='dict'),
        ),
        mutually_exclusive=[["content", "files"]],

        supports_check_mode=True,
    )
    results = OCEnv.run_ansible(module.params, module.check_mode)

    if 'failed' in results:
        module.fail_json(**results)

    module.exit_json(**results)


if __name__ == '__main__':
    main()
