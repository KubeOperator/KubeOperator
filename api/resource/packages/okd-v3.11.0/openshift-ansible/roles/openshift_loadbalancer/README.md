OpenShift HAProxy Loadbalancer
==============================

OpenShift HaProxy Loadbalancer Configuration

Requirements
------------

* Ansible 2.2

This role is intended to be applied to the [lb] host group which is
separate from OpenShift infrastructure components.

This role is not re-entrant. All haproxy configuration lives in a single file.

Role Variables
--------------

From this role:

| Name                                   | Default value |                                                       |
|----------------------------------------|---------------|-------------------------------------------------------|
| openshift_loadbalancer_limit_nofile    | 100000        | Limit number of open files.                           |
| openshift_loadbalancer_global_maxconn  | 20000         | Maximum per-process number of concurrent connections. |
| openshift_loadbalancer_default_maxconn | 20000         | Maximum per-process number of concurrent connections. |
| openshift_loadbalancer_frontends       | none          | List of frontends. See example below.                 |
| openshift_loadbalancer_backends        | none          | List of backends. See example below.                  |
| openshift_image_tag                    | none          | Image tag for containerized haproxy image.            |

Dependencies
------------

* openshift_facts
* os_firewall
* openshift_repos

Example Playbook
----------------

```
- name: Configure loadbalancer hosts
  hosts: lb
  roles:
  - role: openshift_loadbalancer
    openshift_loadbalancer_frontends:
    - name: atomic-openshift-api
      mode: tcp
      options:
      - tcplog
      binds:
      - "*:8443"
      default_backend: atomic-openshift-api
    openshift_loadbalancer_backends:
    - name: atomic-openshift-api
      mode: tcp
      option: tcplog
      balance: source
      servers:
      - name: master1
        address: "192.168.122.221:8443"
	opts: check
      - name: master2
        address: "192.168.122.222:8443"
	opts: check
      - name: master3
        address: "192.168.122.223:8443"
	opts: check
    openshift_image_tag: v3.6.153
```

License
-------

Apache License, Version 2.0

Author Information
------------------

Jason DeTiberus (jdetiber@redhat.com)
