# pylint: skip-file
# flake8: noqa

def main():
    '''
    ansible oc module for clusterrole
    '''

    module = AnsibleModule(
        argument_spec=dict(
            kubeconfig=dict(default='/etc/origin/master/admin.kubeconfig', type='str'),
            state=dict(default='present', type='str',
                       choices=['present', 'absent', 'list']),
            debug=dict(default=False, type='bool'),
            name=dict(default=None, type='str'),
            rules=dict(default=None, type='list'),
        ),
        supports_check_mode=True,
    )

    results = OCClusterRole.run_ansible(module.params, module.check_mode)

    if 'failed' in results:
        module.fail_json(**results)

    module.exit_json(**results)

if __name__ == '__main__':
    main()
