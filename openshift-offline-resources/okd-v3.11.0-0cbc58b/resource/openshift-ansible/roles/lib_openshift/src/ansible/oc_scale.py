# pylint: skip-file
# flake8: noqa

def main():
    '''
    ansible oc module for scaling
    '''

    module = AnsibleModule(
        argument_spec=dict(
            kubeconfig=dict(default='/etc/origin/master/admin.kubeconfig', type='str'),
            state=dict(default='present', type='str', choices=['present', 'list']),
            debug=dict(default=False, type='bool'),
            kind=dict(default='dc', choices=['dc', 'rc'], type='str'),
            namespace=dict(default='default', type='str'),
            replicas=dict(default=None, type='int'),
            name=dict(default=None, type='str'),
        ),
        supports_check_mode=True,
    )
    rval = OCScale.run_ansible(module.params, module.check_mode)
    if 'failed' in rval:
        module.fail_json(**rval)

    module.exit_json(**rval)


if __name__ == '__main__':
    main()
