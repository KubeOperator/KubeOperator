# pylint: skip-file
# flake8: noqa

def main():
    '''
    ansible oc module for volumes
    '''

    module = AnsibleModule(
        argument_spec=dict(
            kubeconfig=dict(default='/etc/origin/master/admin.kubeconfig', type='str'),
            state=dict(default='present', type='str',
                       choices=['present', 'absent', 'list']),
            debug=dict(default=False, type='bool'),
            kind=dict(default='dc', choices=['dc', 'rc', 'pods'], type='str'),
            namespace=dict(default='default', type='str'),
            vol_name=dict(default=None, type='str'),
            name=dict(default=None, type='str'),
            mount_type=dict(default=None,
                            choices=['emptydir', 'hostpath', 'secret', 'pvc', 'configmap'],
                            type='str'),
            mount_path=dict(default=None, type='str'),
            # secrets require a name
            secret_name=dict(default=None, type='str'),
            # pvc requires a size
            claim_size=dict(default=None, type='str'),
            claim_name=dict(default=None, type='str'),
            # configmap requires a name
            configmap_name=dict(default=None, type='str'),
        ),
        supports_check_mode=True,
    )
    rval = OCVolume.run_ansible(module.params, module.check_mode)
    if 'failed' in rval:
        module.fail_json(**rval)

    module.exit_json(**rval)


if __name__ == '__main__':
    main()
