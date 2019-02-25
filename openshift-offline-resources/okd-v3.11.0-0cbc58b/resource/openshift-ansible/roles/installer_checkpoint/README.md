OpenShift-Ansible Installer Checkpoint
======================================

A complete OpenShift cluster installation is comprised of many different
components which can take 30 minutes to several hours to complete.  If the
installation should fail, it could be confusing to understand at which component
the failure occurred.  Additionally, it may be desired to re-run only the
component which failed instead of starting over from the beginning.  Components
which came after the failed component would also need to be run individually.

Design
------

The Installer Checkpoint implements an Ansible callback plugin to allow
displaying and logging of the installer status at the end of a playbook run.

To ensure the callback plugin is loaded, regardless of ansible.cfg file
configuration, the plugin has been placed inside the installer_checkpoint role
which must be called early in playbook execution. The `init/main.yml` playbook
is run first for all entry point playbooks, therefore, the initialization of the
checkpoint plugin has been placed at the beginning of that file.

Playbooks use the [set_stats][set_stats] Ansible module to set a custom stats
variable indicating the status of the phase being executed.

The installer_checkpoint.py callback plugin extends the Ansible
`v2_playbook_on_stats` method, which is called at the end of a playbook run, to
display the status of each phase which was run.  The INSTALLER STATUS report is
displayed immediately following the PLAY RECAP.

Usage
-----

In order to indicate the beginning of a component installation, a play must be
added to the beginning of the main playbook for the component to set the phase
status to "In Progress".  Additionally, a play must be added after the last play
for that component to set the phase status to "Complete".  

The following example shows the first play of the etcd install using the
`set_stats` module for setting the required checkpoint data points.

* `title` - Name of the component phase
* `playbook` - Entry point playbook used to run only this component
* `status` - "In Progress" or "Complete"

```yaml
# playbooks/openshift-etcd/private/config.yml
---
- name: etcd Install Checkpoint Start
  hosts: all
  gather_facts: false
  tasks:
  - name: Set etcd install 'In Progress'
    run_once: true
    set_stats:
      data:
        installer_phase_etcd:
          title: "etcd Install"
          playbook: "playbooks/openshift-etcd/config.yml"
          status: "In Progress"
          start: "{{ lookup('pipe', 'date +%Y%m%d%H%M%SZ') }}"

#...
# Various plays here
#...

- name: etcd Install Checkpoint End
  hosts: all
  gather_facts: false
  tasks:
  - name: Set etcd install 'Complete'
    run_once: true
    set_stats:
      data:
        installer_phase_etcd:
          status: "Complete"
          end: "{{ lookup('pipe', 'date +%Y%m%d%H%M%SZ') }}"
``` 

Examples
--------

Example display of a successful playbook run:

```
PLAY RECAP *********************************************************************
master01.example.com : ok=158  changed=16   unreachable=0    failed=0
node01.example.com   : ok=469  changed=74   unreachable=0    failed=0
node02.example.com   : ok=157  changed=17   unreachable=0    failed=0
localhost            : ok=24   changed=0    unreachable=0    failed=0


INSTALLER STATUS ***************************************************************
Initialization             : Complete (0:02:14)
Health Check               : Complete (0:01:10)
etcd Install               : Complete (0:02:01)
Master Install             : Complete (0:11:43)
Master Additional Install  : Complete (0:00:54)
Node Install               : Complete (0:14:11)
Hosted Install             : Complete (0:03:28)
```

Example display if a failure occurs during execution:

```
INSTALLER STATUS ***************************************************************
Initialization             : Complete (0:02:14)
Health Check               : Complete (0:01:10)
etcd Install               : Complete (0:02:58)
Master Install             : Complete (0:09:20)
Master Additional Install  : In Progress (0:20:04)
    This phase can be restarted by running: playbooks/openshift-master/additional_config.yml
```

[set_stats]: http://docs.ansible.com/ansible/latest/set_stats_module.html
