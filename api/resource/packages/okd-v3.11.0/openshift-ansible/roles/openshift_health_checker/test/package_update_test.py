from openshift_checks.package_update import PackageUpdate


def test_package_update():
    return_value = {}

    def execute_module(module_name=None, module_args=None, *_):
        assert module_name == 'check_yum_update'
        assert 'packages' in module_args
        # empty list of packages means "generic check if 'yum update' will work"
        assert module_args['packages'] == []
        return {'foo': return_value}

    result = PackageUpdate(execute_module).run()
    assert result['foo'] is return_value
