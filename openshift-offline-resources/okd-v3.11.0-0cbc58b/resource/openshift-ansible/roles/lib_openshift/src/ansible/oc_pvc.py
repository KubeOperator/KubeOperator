# pylint: skip-file
# flake8: noqa

#pylint: disable=too-many-branches
def main():
    '''
    ansible oc module for pvc
    '''

    module = AnsibleModule(
        argument_spec=dict(
            kubeconfig=dict(default='/etc/origin/master/admin.kubeconfig', type='str'),
            state=dict(default='present', type='str',
                       choices=['present', 'absent', 'list']),
            debug=dict(default=False, type='bool'),
            name=dict(default=None, required=True, type='str'),
            namespace=dict(default=None, required=True, type='str'),
            volume_capacity=dict(default='1G', type='str'),
            storage_class_name=dict(default=None, required=False, type='str'),
            selector=dict(default=None, required=False, type='dict'),
            access_modes=dict(default=['ReadWriteOnce'], type='list'),
        ),
        supports_check_mode=True,
    )

    rval = OCPVC.run_ansible(module.params, module.check_mode)

    if 'failed' in rval:
        module.fail_json(**rval)

    return module.exit_json(**rval)


if __name__ == '__main__':
    main()
