# pylint: skip-file
# flake8: noqa


# pylint: disable=wrong-import-position
from ansible.module_utils.six import string_types

def main():
    '''
    ansible oc adm module for ca create-server-cert
    '''

    module = AnsibleModule(
        argument_spec=dict(
            state=dict(default='present', type='str', choices=['present']),
            debug=dict(default=False, type='bool'),
            kubeconfig=dict(default='/etc/origin/master/admin.kubeconfig', type='str'),
            backup=dict(default=True, type='bool'),
            force=dict(default=False, type='bool'),
            # oc adm ca create-server-cert [options]
            cert=dict(default=None, type='str'),
            key=dict(default=None, type='str'),
            signer_cert=dict(default='/etc/origin/master/ca.crt', type='str'),
            signer_key=dict(default='/etc/origin/master/ca.key', type='str'),
            signer_serial=dict(default='/etc/origin/master/ca.serial.txt', type='str'),
            hostnames=dict(default=[], type='list'),
            expire_days=dict(default=None, type='int'),
        ),
        supports_check_mode=True,
    )

    results = CAServerCert.run_ansible(module.params, module.check_mode)
    if 'failed' in results:
        return module.fail_json(**results)

    return module.exit_json(**results)


if __name__ == '__main__':
    main()
