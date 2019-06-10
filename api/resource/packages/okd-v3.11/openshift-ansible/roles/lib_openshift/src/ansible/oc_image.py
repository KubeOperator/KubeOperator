# pylint: skip-file
# flake8: noqa


def main():
    '''
    ansible oc module for image import
    '''

    module = AnsibleModule(
        argument_spec=dict(
            kubeconfig=dict(default='/etc/origin/master/admin.kubeconfig', type='str'),
            state=dict(default='present', type='str',
                       choices=['present', 'list']),
            debug=dict(default=False, type='bool'),
            namespace=dict(default='default', type='str'),
            registry_url=dict(default=None, type='str'),
            image_name=dict(default=None, required=True, type='str'),
            image_tag=dict(default=None, type='str'),
            force=dict(default=False, type='bool'),
        ),

        supports_check_mode=True,
    )

    rval = OCImage.run_ansible(module.params, module.check_mode)

    if 'failed' in rval:
        module.fail_json(**rval)

    module.exit_json(**rval)

if __name__ == '__main__':
    main()
