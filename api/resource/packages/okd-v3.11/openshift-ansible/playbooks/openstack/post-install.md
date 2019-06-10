# Post-Install

* [Configure DNS](#configure-dns)
* [Get the `oc` Client](#get-the-oc-client)
* [Log in Using the Command Line](#log-in-using-the-command-line)
* [Access the UI](#access-the-ui)
* [Run Custom Post-Provision Actions](#run-custom-post-provision-actions)
* [Manage Trunk ports and subports](#manage-trunk-ports-and-subports)


## Configure DNS

OpenShift requires two public DNS records to function fully. The first one points to
the master/load balancer and provides the UI/API access. The other one is a
wildcard domain that resolves app route requests to the infra node. A private DNS
server and records are not required and not managed here.

If you followed the default installation from the README section, there is no
DNS configured. You should add two entries to the `/etc/hosts` file on the
Ansible host (where you to do a quick validation. A real deployment will
however require a DNS server with the following entries set.

In either case, the IP addresses for the API and routers will be printed
out at the end of the deployment.

The first one is your API/UI address and the second one is the router address.
Depending on your load balancer configuration they may or may not be the same.

In this example, we will use `10.40.128.130` for the `public_api_ip` and
`10.40.128.134` for `public_router_ip`.

Add the following entries to your `/etc/hosts`:

```
10.40.128.130 console.openshift.example.com
10.40.128.134 cakephp-mysql-example-test.apps.openshift.example.com
```

This points the cluster domain (as defined in the
`openshift_master_cluster_public_hostname` Ansible variable in `OSEv3`) to the
master node and any routes for deployed apps to the infra node.

If you deploy another app, it will end up with a different URL (e.g.
myapp-test.apps.openshift.example.com) and you will need to add that too.  This
is why a real deployment should always run a DNS where the second entry will be
a wildcard `*.apps.openshift.example.com).

This will be sufficient to validate the cluster here.

Take a look at the [External DNS][external-dns] section for
configuring a DNS service.


## Get the `oc` Client

The OpenShift command line client (called `oc`) can be downloaded and extracted
from `openshift-origin-client-tools` on the OpenShift release page:

https://github.com/openshift/origin/releases/latest/

You can also copy it from the master node:

    $ ansible -i inventory masters[0] -m fetch -a "src=/bin/oc dest=oc"

Once you obtain the `oc` binary, remember to put it in your `PATH`.


## Log in Using the Command Line

Once the `oc` client is available, you can login using the URLs specified in `/etc/hosts`:

```
oc login --insecure-skip-tls-verify=true https://console.openshift.example.com:8443 -u user -p password
oc new-project test
oc new-app --template=cakephp-mysql-example
oc status -v
curl http://cakephp-mysql-example-test.apps.openshift.example.com
```

This will trigger an image build. You can run `oc logs -f
bc/cakephp-mysql-example` to follow its progress.

Wait until the build has finished and both pods are deployed and running:

```
$ oc status -v
In project test on server https://master-0.openshift.example.com:8443

http://cakephp-mysql-example-test.apps.openshift.example.com (svc/cakephp-mysql-example)
  dc/cakephp-mysql-example deploys istag/cakephp-mysql-example:latest <-
    bc/cakephp-mysql-example source builds https://github.com/openshift/cakephp-ex.git on openshift/php:7.0
    deployment #1 deployed about a minute ago - 1 pod

svc/mysql - 172.30.144.36:3306
  dc/mysql deploys openshift/mysql:5.7
    deployment #1 deployed 3 minutes ago - 1 pod

Info:
  * pod/cakephp-mysql-example-1-build has no liveness probe to verify pods are still running.
    try: oc set probe pod/cakephp-mysql-example-1-build --liveness ...
View details with 'oc describe <resource>/<name>' or list everything with 'oc get all'.

```

You can now look at the deployed app using its route:

```
$ curl http://cakephp-mysql-example-test.apps.openshift.example.com
```

Its `title` should say: "Welcome to OpenShift".


## Access the UI

You can access the OpenShift cluster with a web browser by going to:

https://master-0.openshift.example.com:8443

Note that for this to work, the OpenShift nodes must be accessible
from your computer and its DNS configuration must use the cluster's
DNS.


## Run Custom Post-Provision Actions

A custom playbook can be run like this:

```
ansible-playbook --private-key ~/.ssh/openshift -i inventory/ openshift-ansible-contrib/playbooks/provisioning/openstack/custom-actions/custom-playbook.yml
```

If you'd like to limit the run to one particular host, you can do so as follows:

```
ansible-playbook --private-key ~/.ssh/openshift -i inventory/ openshift-ansible-contrib/playbooks/provisioning/openstack/custom-actions/custom-playbook.yml -l app-node-0.openshift.example.com
```

You can also create your own custom playbook. Here are a few examples:

### Add Additional YUM Repositories

```
---
- hosts: app
  tasks:

  # enable EPL
  - name: Add repository
    yum_repository:
      name: epel
      description: EPEL YUM repo
      baseurl: https://download.fedoraproject.org/pub/epel/$releasever/$basearch/
```

This example runs against app nodes. The list of options include:

  - OSEv3 (all created hosts: app, infra, masters, etcd, glusterfs, lb, nfs)
  - openstack_nodes (all OpenShift hosts: app, infra, masters, etcd)
  - openstack_compute_nodes
  - openstack_master_nodes
  - openstack_infra_nodes

### Attach Additional RHN Pools

```
---
- hosts: OSEv3
  tasks:
  - name: Attach additional RHN pool
    become: true
    command: "/usr/bin/subscription-manager attach --pool=<pool ID>"
    register: attach_rhn_pool_result
    until: attach_rhn_pool_result.rc == 0
    retries: 10
    delay: 1
```

This playbook runs against all cluster nodes. In order to help prevent slow connectivity
problems, the task is retried 10 times in case of initial failure.
Note that in order for this example to work in your deployment, your servers must use the RHEL image.

### Add Extra Docker Registry URLs

This playbook is located in the [custom-actions](https://github.com/openshift/openshift-ansible-contrib/tree/master/playbooks/provisioning/openstack/custom-actions) directory.

It adds URLs passed as arguments to the docker configuration program.
Going into more detail, the configuration program (which is in the YAML format) is loaded into an ansible variable
([lines 27-30](https://github.com/openshift/openshift-ansible-contrib/blob/master/playbooks/provisioning/openstack/custom-actions/add-docker-registry.yml#L27-L30))
and in its structure, `registries` and `insecure_registries` sections are expanded with the newly added items
([lines 56-76](https://github.com/openshift/openshift-ansible-contrib/blob/master/playbooks/provisioning/openstack/custom-actions/add-docker-registry.yml#L56-L76)).
The new content is then saved into the original file
([lines 78-82](https://github.com/openshift/openshift-ansible-contrib/blob/master/playbooks/provisioning/openstack/custom-actions/add-docker-registry.yml#L78-L82))
and docker is restarted.

Example usage:
```
ansible-playbook -i <inventory> openshift-ansible-contrib/playbooks/provisioning/openstack/custom-actions/add-docker-registry.yml  --extra-vars '{"registries": "reg1", "insecure_registries": ["ins_reg1","ins_reg2"]}'
```

### Add Extra CAs to the Trust Chain

This playbook is also located in the [custom-actions](https://github.com/openshift/openshift-ansible-contrib/blob/master/playbooks/provisioning/openstack/custom-actions) directory.
It copies passed CAs to the trust chain location and updates the trust chain on each selected host.

Example usage:
```
ansible-playbook -i <inventory> openshift-ansible-contrib/playbooks/provisioning/openstack/custom-actions/add-cas.yml --extra-vars '{"ca_files": [<absolute path to ca1 file>, <absolute path to ca2 file>]}'
```

Please consider contributing your custom playbook back to openshift-ansible!

A library of custom post-provision actions exists in `openshift-ansible-contrib/playbooks/provisioning/openstack/custom-actions`. Playbooks include:

* [add-yum-repos.yml](https://github.com/openshift/openshift-ansible-contrib/blob/master/playbooks/provisioning/openstack/custom-actions/add-yum-repos.yml): adds a list of custom yum repositories to every node in the cluster
* [add-rhn-pools.yml](https://github.com/openshift/openshift-ansible-contrib/blob/master/playbooks/provisioning/openstack/custom-actions/add-rhn-pools.yml): attaches a list of additional RHN pools to every node in the cluster
* [add-docker-registry.yml](https://github.com/openshift/openshift-ansible-contrib/blob/master/playbooks/provisioning/openstack/custom-actions/add-docker-registry.yml): adds a list of docker registries to the docker configuration on every node in the cluster
* [add-cas.yml](https://github.com/openshift/openshift-ansible-contrib/blob/master/playbooks/provisioning/openstack/custom-actions/add-rhn-pools.yml): adds a list of CAs to the trust chain on every node in the cluster


[external-dns]: ./configuration.md#dns-configuration

## Manage Trunk ports and subports

Running OpenShift on top of OpenStack VMs without the problem of double
encapsulation is achieved by using kuryr which leverages Neutron Trunk Ports
feature.

With the Trunk Ports (also known as VLAN aware VMs), we can create a trunk,
associate a parent port for it that will be used by the VM (in our case, the
master, infra and app-node VMs), and then we can create a normal Neutron port
and attach it to the trunk to become a subport of it. These subports are later
used by the pods running inside those VMs to connect to the Neutron networks.

Next we show a few example of how to manage trunks, parents and subports.

#### Create a Trunk port

```
# Create a Neutron port
openstack port create --network VM_NETWORK parent-port-0
# Create a trunk with that port as parent port
openstack network trunk create --parent-port parent-port-0 trunk-0
```

Note you need to first create the port, then the trunk with that port, and only
then you can create the VM by using the parent port created, in the example
above parent-port-0.

#### Attach subports to the trunk

```
# Create a Neutron port
openstack port create --network POD_NETWORK subport-0
# Attach the port as a subport of the trunk
openstack network trunk set --subport port=subport-0,segmentation-type=vlan,segmentation-id=101 trunk-0
```

#### Remove subports

In order to remove the subports Neutron ports, you need to first detach them
from the trunk, and then delete them:

```
# Detach subport from trunk
openstack network trunk unset --subport subport-0 trunk-0
# Remove port (as usual)
openstack port delete subport-0
```

#### Create subports for the Kuryr Ports Pools

Kuryr Ports Pool is a feature to speed up containers boot up time by reducing
the number of interactions between Kuryr-controller and Neutron API -- which in
turn reduces the load on the Neutron server, also improving the overall
performance. To achieve this, the Kuryr-controller maintains a pool of neutron
ports ready to be used -- instead of creating a port/subport upon pod creation.
For the nested case where pods will be created inside an OpenShift cluster
installed on top of OpenStack VMs, there will be several pools, one for each
pair of:
- Trunk port (i.e., VM belonging to OpenShift cluster)
- Set of security groups used by the pods
- Neutron Network used by the pods
- Project ID used to create the pods (i.e., OpenStack tenant)

Note, default kuryr drivers creates all the pods with the same security groups
set, subnets and project. Thus, in practice there is a pool per trunk port,
i.e., per VM belonging to the OpenShift cluster.

In order to manually populate one specific pool, the next can be done:

```
# Create port with the right project_id, security group set and network
openstack port create --network POD_NETWORK --security-group SG_1
--security-group SG_2 subport-1
# Attach the subport to the trunk where you want to add the port to the pool
openstack network trunk set --subport port=subport-1,segmentation-type=vlan,segmentation-id=1 APP_NODE_VM_TRUNK
```

Note you need to choose a segmentation id that is not already in use at that
trunk. To see the current subports attached to that trunk, and their associated
segmentation ids, you can do:

```
openstack network trunk show APP_NODE_VM_TRUNK
+-----------------+--------------------------------------------------------------------------------------------------+
| Field           | Value                                                                                            |
+-----------------+--------------------------------------------------------------------------------------------------+
| admin_state_up  | UP                                                                                               |
| created_at      | 2018-03-28T15:06:54Z                                                                             |
| description     |                                                                                                  |
| id              | 9048c109-c1aa-4a41-9508-71b2ba98f3b0                                                             |
| name            | APP_NODE_VM_TRUNK                                                                                |
| port_id         | 4180a2e5-e184-424a-93d4-54b48490f50d                                                             |
| project_id      | a05f6ec0abd04cba80cd160f8baaac99                                                                 |
| revision_number | 43                                                                                               |
| status          | ACTIVE                                                                                           |
| sub_ports       | port_id='1de77073-7127-4c39-a47b-cef15f98849c', segmentation_id='101', segmentation_type='vlan'  |
| tags            | []                                                                                               |
| tenant_id       | a05f6ec0abd04cba80cd160f8baaac99                                                                 |
| updated_at      | 2018-03-29T06:12:39Z                                                                             |
+-----------------+--------------------------------------------------------------------------------------------------+
```

Finally, next time the kuryr-controller pod gets restarted it will recover the
subports attached to each trunk, and add them to their respective pools -- if
they are not in used by a pod already. This can also be forced by manually
restarting the kuryr-controller by killing the running pod:

```
kubectl -n openshift-infra delete pod kuryr-controller-XXXXX
```
