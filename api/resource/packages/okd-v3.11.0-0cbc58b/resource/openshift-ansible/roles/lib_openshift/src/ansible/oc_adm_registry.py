# pylint: skip-file
# flake8: noqa

def main():
    '''
    ansible oc module for registry
    '''

    module = AnsibleModule(
        argument_spec=dict(
            state=dict(default='present', type='str',
                       choices=['present', 'absent']),
            debug=dict(default=False, type='bool'),
            namespace=dict(default='default', type='str'),
            name=dict(default=None, required=True, type='str'),

            kubeconfig=dict(default='/etc/origin/master/admin.kubeconfig', type='str'),
            images=dict(default=None, type='str'),
            latest_images=dict(default=False, type='bool'),
            labels=dict(default=None, type='dict'),
            ports=dict(default=['5000'], type='list'),
            replicas=dict(default=1, type='int'),
            selector=dict(default=None, type='str'),
            service_account=dict(default='registry', type='str'),
            mount_host=dict(default=None, type='str'),
            volume_mounts=dict(default=None, type='list'),
            env_vars=dict(default={}, type='dict'),
            edits=dict(default=[], type='list'),
            enforce_quota=dict(default=False, type='bool'),
            force=dict(default=False, type='bool'),
            daemonset=dict(default=False, type='bool'),
            tls_key=dict(default=None, type='str'),
            tls_certificate=dict(default=None, type='str'),
        ),

        supports_check_mode=True,
    )

    results = Registry.run_ansible(module.params, module.check_mode)
    if 'failed' in results:
        module.fail_json(**results)

    module.exit_json(**results)


if __name__ == '__main__':
    main()
