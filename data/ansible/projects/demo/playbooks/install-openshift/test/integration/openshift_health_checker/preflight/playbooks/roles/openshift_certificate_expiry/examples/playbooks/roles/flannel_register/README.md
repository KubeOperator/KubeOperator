Role Name
=========

Register flannel configuration into etcd

Requirements
------------

This role assumes it's being deployed on a RHEL/Fedora based host with package
named 'flannel' available via yum, in version superior to 0.3.

Role Variables
--------------

| Name                | Default value                                      | Description                                     |
|---------------------|----------------------------------------------------|-------------------------------------------------|
| flannel_network     | {{ openshift.common.portal_net }} or 172.16.1.1/16 | interface to use for inter-host communication   |
| flannel_min_network | {{ min_network }} or 172.16.5.0                    | beginning of IP range for the subnet allocation |
| flannel_subnet_len  | 24                                                 | size of the subnet allocated to each host       |
| flannel_etcd_key    | /openshift.com/network                             | etcd prefix                                     |
| etcd_hosts          | etcd_urls                                          | a list of etcd endpoints                        |
| etcd_conf_dir       | {{ openshift.common.config_base }}/master          | SSL certificates directory                      |
| etcd_peer_ca_file   | {{ etcd_conf_dir }}/ca.crt                         | SSL CA to use for etcd                          |
| etcd_peer_cert_file | {{ etcd_conf_dir }}/master.etcd-client.crt         | SSL cert to use for etcd                        |
| etcd_peer_key_file  | {{ etcd_conf_dir }}/master.etcd-client.key         | SSL key to use for etcd                         |

Dependencies
------------

openshift_facts

Example Playbook
----------------

    - hosts: openshift_master
      roles:
         - { flannel_register }

License
-------

Apache License, Version 2.0

Author Information
------------------

Sylvain Baubeau <sbaubeau@redhat.com>
