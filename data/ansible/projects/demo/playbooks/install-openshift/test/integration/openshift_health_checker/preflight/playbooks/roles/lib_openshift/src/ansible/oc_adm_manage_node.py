# pylint: skip-file
# flake8: noqa


def main():
    '''
    ansible oadm module for manage-node
    '''

    module = AnsibleModule(
        argument_spec=dict(
            debug=dict(default=False, type='bool'),
            kubeconfig=dict(default='/etc/origin/master/admin.kubeconfig', type='str'),
            node=dict(default=None, type='list'),
            selector=dict(default=None, type='str'),
            pod_selector=dict(default=None, type='str'),
            schedulable=dict(default=None, type='bool'),
            list_pods=dict(default=False, type='bool'),
            evacuate=dict(default=False, type='bool'),
            dry_run=dict(default=False, type='bool'),
            force=dict(default=False, type='bool'),
            grace_period=dict(default=None, type='int'),
        ),
        mutually_exclusive=[["selector", "node"], ['evacuate', 'list_pods'], ['list_pods', 'schedulable']],
        required_one_of=[["node", "selector"]],

        supports_check_mode=True,
    )
    results = ManageNode.run_ansible(module.params, module.check_mode)

    if 'failed' in results:
        module.fail_json(**results)

    module.exit_json(**results)


if __name__ == "__main__":
    main()
