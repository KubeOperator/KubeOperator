# pylint: skip-file
# flake8: noqa

def main():
    '''
    ansible oc module for approving certificate signing requests
    '''

    module = AnsibleModule(
        argument_spec=dict(
            kubeconfig=dict(default='/etc/origin/master/admin.kubeconfig', type='str'),
            state=dict(default='approve', type='str',
                       choices=['approve', 'deny', 'list']),
            debug=dict(default=False, type='bool'),
            nodes=dict(default=None, type='list'),
            timeout=dict(default=30, type='int'),
            approve_all=dict(default=False, type='bool'),
            service_account=dict(default='system:serviceaccount:openshift-infra:node-bootstrapper', type='str'),
            fail_on_timeout=dict(default=False, type='bool'),
        ),
        supports_check_mode=True,
        mutually_exclusive=[['approve_all', 'nodes']],
    )

    if module.params['nodes'] == []:
        module.fail_json(**dict(failed=True, msg='Please specify hosts.'))

    rval = OCcsr.run_ansible(module.params, module.check_mode)

    # If we timed out then we weren't finished. Fail if user requested to fail.
    if (module.params['timeout'] > 0 and
            module.params['fail_on_timeout'] and
            rval['timeout']):
        return module.fail_json(msg='Timed out accepting certificate signing requests. Failing as requested.', **rval)

    if 'failed' in rval:
        return module.fail_json(**rval)

    return module.exit_json(**rval)


if __name__ == '__main__':
    main()
