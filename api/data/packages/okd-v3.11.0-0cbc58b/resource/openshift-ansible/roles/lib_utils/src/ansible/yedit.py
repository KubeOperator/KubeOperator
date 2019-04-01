# flake8: noqa
# pylint: skip-file


# pylint: disable=too-many-branches
def main():
    ''' ansible oc module for secrets '''

    module = AnsibleModule(
        argument_spec=dict(
            state=dict(default='present', type='str',
                       choices=['present', 'absent', 'list']),
            debug=dict(default=False, type='bool'),
            src=dict(default=None, type='str'),
            content=dict(default=None),
            content_type=dict(default='yaml', choices=['yaml', 'json']),
            key=dict(default='', type='str'),
            value=dict(),
            value_type=dict(default='', type='str'),
            update=dict(default=False, type='bool'),
            append=dict(default=False, type='bool'),
            index=dict(default=None, type='int'),
            curr_value=dict(default=None, type='str'),
            curr_value_format=dict(default='yaml',
                                   choices=['yaml', 'json', 'str'],
                                   type='str'),
            backup=dict(default=False, type='bool'),
            backup_ext=dict(default=".{}".format(time.strftime("%Y%m%dT%H%M%S")), type='str'),
            separator=dict(default='.', type='str'),
            edits=dict(default=None, type='list'),
        ),
        mutually_exclusive=[["curr_value", "index"], ['update', "append"]],
        required_one_of=[["content", "src"]],
    )

    # Verify we recieved either a valid key or edits with valid keys when receiving a src file.
    # A valid key being not None or not ''.
    if module.params['src'] is not None:
        key_error = False
        edit_error = False

        if module.params['key'] in [None, '']:
            key_error = True

        if module.params['edits'] in [None, []]:
            edit_error = True

        else:
            for edit in module.params['edits']:
                if edit.get('key') in [None, '']:
                    edit_error = True
                    break

        if key_error and edit_error:
            module.fail_json(failed=True, msg='Empty value for parameter key not allowed.')

    rval = Yedit.run_ansible(module.params)
    if 'failed' in rval and rval['failed']:
        module.fail_json(**rval)

    module.exit_json(**rval)


if __name__ == '__main__':
    main()
