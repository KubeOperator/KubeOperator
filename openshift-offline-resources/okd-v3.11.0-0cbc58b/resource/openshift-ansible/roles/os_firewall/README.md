OS Firewall
===========

OS Firewall manages firewalld and iptables installation.
case.

Note: firewalld is not supported on Atomic Host
https://bugzilla.redhat.com/show_bug.cgi?id=1403331

Requirements
------------

Ansible 2.2

Role Variables
--------------

| Name                      | Default |                                        |
|---------------------------|---------|----------------------------------------|
| os_firewall_use_firewalld | False   | If false, use iptables                 |

Dependencies
------------

None.

Example Playbook
----------------

Use iptables:
```
---
- hosts: servers
  task:
  - import_role:
      name: os_firewall
    vars:
      os_firewall_use_firewalld: false
```

Use firewalld:
```
---
- hosts: servers
  vars:
  tasks:
  - import_role:
      name: os_firewall
    vars:
      os_firewall_use_firewalld: true
```

License
-------

Apache License, Version 2.0

Author Information
------------------
Jason DeTiberus - jdetiber@redhat.com
