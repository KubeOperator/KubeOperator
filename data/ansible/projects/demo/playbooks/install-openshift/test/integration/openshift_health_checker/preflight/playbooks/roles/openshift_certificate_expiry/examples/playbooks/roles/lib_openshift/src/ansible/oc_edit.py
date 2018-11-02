# pylint: skip-file
# flake8: noqa


def main():
    '''
    ansible oc module for editing objects
    '''

    module = AnsibleModule(
        argument_spec=dict(
            kubeconfig=dict(default='/etc/origin/master/admin.kubeconfig', type='str'),
            state=dict(default='present', type='str',
                       choices=['present']),
            debug=dict(default=False, type='bool'),
            namespace=dict(default='default', type='str'),
            name=dict(default=None, required=True, type='str'),
            kind=dict(required=True, type='str'),
            file_name=dict(default=None, type='str'),
            file_format=dict(default='yaml', type='str'),
            content=dict(default=None, type='dict'),
            force=dict(default=False, type='bool'),
            separator=dict(default='.', type='str'),
            edits=dict(default=None, type='list'),
        ),
        supports_check_mode=True,
        mutually_exclusive=[['content', 'edits']],
        required_one_of=[['content', 'edits']],
    )

    rval = Edit.run_ansible(module.params, module.check_mode)
    if 'failed' in rval:
        module.fail_json(**rval)

    module.exit_json(**rval)

if __name__ == '__main__':
    main()
