Openshift descheduler
=====================

Install the descheduler

Role Variables
--------------
Check defaults/main.yml

Installing Descheduler
--------------------

```
openshift_descheduler_install=true
```

```
ansible-playbook -i <inventory-file> playbooks/openshift-descheduler/config.yml
```

Uninstalling Descheduler
--------------------

```
openshift_descheduler_install=false
```

```
ansible-playbook -i <inventory-file> playbooks/openshift-descheduler/config.yml
```

Notes
-----

This is currently experimental software.  This role allows users to install the descheduler and the necessary authorization pieces that allow the descheduler to function.

https://github.com/openshift/descheduler

License
-------

Apache License, Version 2.0

Author Information
------------------

Openshift
