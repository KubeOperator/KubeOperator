# pylint: skip-file
# flake8: noqa


def main():
    '''
    ansible oc adm module for user policy
    '''

    module = AnsibleModule(
        argument_spec=dict(
            state=dict(default='present', type='str',
                       choices=['present', 'absent']),
            debug=dict(default=False, type='bool'),
            resource_name=dict(required=True, type='str'),
            namespace=dict(default='default', type='str'),
            role_namespace=dict(default=None, type='str'),
            rolebinding_name=dict(default=None, type='str'),
            kubeconfig=dict(default='/etc/origin/master/admin.kubeconfig', type='str'),

            user=dict(required=True, type='str'),
            resource_kind=dict(required=True, choices=['role', 'cluster-role', 'scc'], type='str'),
        ),
        supports_check_mode=True,
    )

    results = PolicyUser.run_ansible(module.params, module.check_mode)

    if 'failed' in results:
        module.fail_json(**results)

    module.exit_json(**results)


if __name__ == "__main__":
    main()
