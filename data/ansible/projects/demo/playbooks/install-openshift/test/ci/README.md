* Copy `test/ci/vars.yml.sample` to `test/ci/vars.yml`
* Adjust it your liking - this would be the host configuration
* Adjust `inventory/group_vars/OSEv3/vars.yml` - this would be Origin-specific config
* Provision instances via `ansible-playbook -vv -i test/ci/inventory/ test/ci/launch.yml`
  This would place inventory file in `test/ci/inventory/hosts` and run prerequisites and deploy.

* Once the setup is complete run `ansible-playbook -vv -i test/ci/inventory/ test/ci/deprovision.yml`
