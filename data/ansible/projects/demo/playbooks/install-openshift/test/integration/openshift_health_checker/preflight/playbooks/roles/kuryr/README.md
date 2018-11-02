## OpenStack Kuryr

Install Kuryr CNI components (kuryr-controller, kuryr-cni) on Master and worker
nodes. Kuryr uses OpenStack Networking service (Neutron) to provide network for
pods. This allows to have interconnectivity between pods and OpenStack VMs.

## Requirements

* Ansible 2.2+
* Centos/ RHEL 7.3+

## Current Kuryr restrictions when used with OpenShift

* Openshift Origin only
* OpenShift on OpenStack Newton or newer (only with Trunk ports)

## Key Ansible inventory Kuryr master configuration parameters

* ``openshift_use_kuryr=True``
* ``openshift_use_openshift_sdn=False``
* ``openshift_sdn_network_plugin_name='cni'``
* ``kuryr_cni_link_interface=eth0``
* ``kuryr_openstack_auth_url=keystone_url``
* ``kuryr_openstack_user_domain_name=Default``
* ``kuryr_openstack_user_project_name=Default``
* ``kuryr_openstack_project_id=project_uuid``
* ``kuryr_openstack_username=kuryr``
* ``kuryr_openstack_password=kuryr_pass``
* ``kuryr_openstack_ca=/etc/ssl/ca.crt` (defaults to ``OS_CACERT`` env var)
* ``kuryr_openstack_pod_sg_id=pod_security_group_uuid``
* ``kuryr_openstack_pod_subnet_id=pod_subnet_uuid``
* ``kuryr_openstack_pod_service_id=service_subnet_uuid``
* ``kuryr_openstack_pod_project_id=pod_project_uuid``
* ``kuryr_openstack_worker_nodes_subnet_id=worker_nodes_subnet_uuid``
* ``kuryr_openstack_pool_driver=nested``
* ``kuryr_openstack_pool_max=0``
* ``kuryr_openstack_pool_min=1``
* ``kuryr_openstack_pool_batch=5``
* ``kuryr_openstack_pool_update_frequency=20``
* ``openshift_kuryr_precreate_subports=5``
* ``openshift_kuryr_device_owner=compute:kuryr``

## OpenShift API loadbalancer

Kuryr is connecting to OpenShift API through the load balancer created by
OpenStack playbook. Both Octavia and Neutron LBaaS v2 hardcode 50 seconds as
client and server inactivity timeout. This is a low value for Kuryr, which is
watching K8s API forever. If connection will get closed by the LB, Kuryr will
restart it with some messages about it in the logs.

If you have access to your OpenStack cloud configuration you can disable the
timeouts by providing custom HA proxy templates to your LBaaS v2 or Octavia
installations. It's controlled by ``[haproxy]jinja_config_template`` option in
Neutron LBaaS v2 and ``[haproxy_amphora]haproxy_template`` in Octavia's config.
Please note that such configuration change will affect all the load balancers
created in the cloud.

## Kuryr resources

* [Kuryr documentation](https://docs.openstack.org/kuryr-kubernetes/latest/)
* [Installing Kuryr containerized](https://docs.openstack.org/kuryr-kubernetes/latest/installation/containerized.html)
