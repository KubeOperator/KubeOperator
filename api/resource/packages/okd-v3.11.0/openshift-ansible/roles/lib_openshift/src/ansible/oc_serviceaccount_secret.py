# pylint: skip-file
# flake8: noqa

def main():
    '''
    ansible oc module to manage service account secrets.
    '''

    module = AnsibleModule(
        argument_spec=dict(
            kubeconfig=dict(default='/etc/origin/master/admin.kubeconfig', type='str'),
            state=dict(default='present', type='str',
                       choices=['present', 'absent', 'list']),
            debug=dict(default=False, type='bool'),
            namespace=dict(default=None, required=True, type='str'),
            secret=dict(default=None, type='str'),
            service_account=dict(required=True, type='str'),
        ),
        supports_check_mode=True,
    )

    rval = OCServiceAccountSecret.run_ansible(module.params, module.check_mode)
    if 'failed' in rval:
        module.fail_json(**rval)

    module.exit_json(**rval)

if __name__ == '__main__':
    main()
