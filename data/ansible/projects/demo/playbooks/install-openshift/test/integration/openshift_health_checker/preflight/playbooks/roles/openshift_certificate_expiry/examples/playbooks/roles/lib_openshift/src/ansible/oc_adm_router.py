# pylint: skip-file
# flake8: noqa


def main():
    '''
    ansible oc module for router
    '''

    module = AnsibleModule(
        argument_spec=dict(
            state=dict(default='present', type='str',
                       choices=['present', 'absent']),
            debug=dict(default=False, type='bool'),
            namespace=dict(default='default', type='str'),
            name=dict(default='router', type='str'),

            kubeconfig=dict(default='/etc/origin/master/admin.kubeconfig', type='str'),
            default_cert=dict(default=None, type='str'),
            cert_file=dict(default=None, type='str'),
            key_file=dict(default=None, type='str'),
            images=dict(default=None, type='str'), #'registry.access.redhat.com/openshift3/ose-${component}:${version}'
            latest_images=dict(default=False, type='bool'),
            labels=dict(default=None, type='dict'),
            ports=dict(default=['80:80', '443:443'], type='list'),
            replicas=dict(default=1, type='int'),
            selector=dict(default=None, type='str'),
            service_account=dict(default='router', type='str'),
            router_type=dict(default='haproxy-router', type='str'),
            host_network=dict(default=True, type='bool'),
            # external host options
            external_host=dict(default=None, type='str'),
            external_host_vserver=dict(default=None, type='str'),
            external_host_insecure=dict(default=False, type='bool'),
            external_host_partition_path=dict(default=None, type='str'),
            external_host_username=dict(default=None, type='str'),
            external_host_password=dict(default=None, type='str', no_log=True),
            external_host_private_key=dict(default=None, type='str', no_log=True),
            # Stats
            stats_user=dict(default=None, type='str'),
            stats_password=dict(default=None, type='str', no_log=True),
            stats_port=dict(default=1936, type='int'),
            # extra
            cacert_file=dict(default=None, type='str'),
            # edits
            edits=dict(default=[], type='list'),
        ),
        mutually_exclusive=[["router_type", "images"],
                            ["key_file", "default_cert"],
                            ["cert_file", "default_cert"],
                            ["cacert_file", "default_cert"],
                           ],

        required_together=[['cacert_file', 'cert_file', 'key_file']],
        supports_check_mode=True,
    )
    results = Router.run_ansible(module.params, module.check_mode)

    if 'failed' in results:
        module.fail_json(**results)

    module.exit_json(**results)


if __name__ == '__main__':
    main()
