# pylint: skip-file
# flake8: noqa

def main():
    ''' ansible oc module for version '''

    module = AnsibleModule(
        argument_spec=dict(
            kubeconfig=dict(default='/etc/origin/master/admin.kubeconfig', type='str'),
            state=dict(default='list', type='str',
                       choices=['list']),
            debug=dict(default=False, type='bool'),
        ),
        supports_check_mode=True,
    )

    rval = OCVersion.run_ansible(module.params)
    if 'failed' in rval:
        module.fail_json(**rval)


    module.exit_json(**rval)


if __name__ == '__main__':
    main()
