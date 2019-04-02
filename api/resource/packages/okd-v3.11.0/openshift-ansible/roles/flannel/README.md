Role Name
=========

Configure flannel on openshift nodes

Requirements
------------

* Ansible 2.2
* This role assumes it's being deployed on a RHEL/Fedora based host with package
named 'flannel' available via yum or dnf (conditionally), in version superior
to 0.3.

Role Variables
--------------

| Name                 | Default value                           | Description                                   |
|----------------------|-----------------------------------------|-----------------------------------------------|
| flannel_interface    | ansible_default_ipv4.interface          | interface to use for inter-host communication |
| flannel_etcd_key     | /openshift.com/network                  | etcd prefix                                   |
| etcd_hosts           | etcd_urls                               | a list of etcd endpoints                      |
| etcd_cert_config_dir | {{ openshift.common.config_base }}/node | SSL certificates directory                    |
| etcd_peer_ca_file    | {{ etcd_conf_dir }}/ca.crt              | SSL CA to use for etcd                        |
| etcd_peer_cert_file  | Openshift SSL cert                      | SSL cert to use for etcd                      |
| etcd_peer_key_file   | Openshift SSL key                       | SSL key to use for etcd                       |

Dependencies
------------

Example Playbook
----------------

    - hosts: openshift_node
      roles:
        - { role: flannel, etcd_urls: ['https://127.0.0.1:2379'] }

License
-------

Apache License, Version 2.0

Author Information
------------------

Sylvain Baubeau <sbaubeau@redhat.com>
