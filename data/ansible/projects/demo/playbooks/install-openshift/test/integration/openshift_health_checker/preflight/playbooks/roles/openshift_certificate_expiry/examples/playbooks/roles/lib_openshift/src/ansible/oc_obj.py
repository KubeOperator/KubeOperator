# pylint: skip-file
# flake8: noqa

# pylint: disable=too-many-branches
def main():
    '''
    ansible oc module for services
    '''

    module = AnsibleModule(
        argument_spec=dict(
            kubeconfig=dict(default='/etc/origin/master/admin.kubeconfig', type='str'),
            state=dict(default='present', type='str',
                       choices=['present', 'absent', 'list']),
            debug=dict(default=False, type='bool'),
            namespace=dict(default='default', type='str'),
            all_namespaces=dict(defaul=False, type='bool'),
            name=dict(default=None, type='str'),
            files=dict(default=None, type='list'),
            kind=dict(required=True, type='str'),
            delete_after=dict(default=False, type='bool'),
            content=dict(default=None, type='dict'),
            force=dict(default=False, type='bool'),
            selector=dict(default=None, type='str'),
            field_selector=dict(default=None, type='str'),
        ),
        mutually_exclusive=[["content", "files"], ["selector", "name"], ["field_selector", "name"]],

        supports_check_mode=True,
    )
    rval = OCObject.run_ansible(module.params, module.check_mode)
    if 'failed' in rval:
        module.fail_json(**rval)

    module.exit_json(**rval)

if __name__ == '__main__':
    main()
