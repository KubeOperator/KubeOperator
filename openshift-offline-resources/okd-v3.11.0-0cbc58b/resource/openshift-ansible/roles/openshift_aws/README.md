openshift_aws
==================================

Provision AWS infrastructure and instances.

This role contains many task-areas to provision resources and perform actions
against an AWS account for the purposes of dynamically building an openshift
cluster.

This role is primarily intended to be used with "import_role" and "tasks_from".

import_role can be called from the tasks section in a play.  See example
playbook below for reference.

These task-areas are:

* provision a vpc: vpc.yml
* provision elastic load balancers: elb.yml
* upload IAM ssl certificates to use with load balancers: iam_cert.yml
* provision an S3 bucket: s3.yml
* provision an instance to build an AMI: provision_instance.yml
* provision a security group in AWS: security_group.yml
* provision ssh keys and users in AWS: ssh_keys.yml
* provision an AMI in AWS: seal_ami.yml
* provision scale groups: scale_group.yml
* provision launch configs: launch_config.yml

Requirements
------------

* Ansible 2.3
* Boto

Appropriate AWS credentials and permissions are required.




Example Playbook
----------------

```yaml
- import_role:
    name: openshift_aws
    tasks_from: vpc.yml
  vars:
    openshift_aws_clusterid: test
    openshift_aws_region: us-east-1
```

License
-------

Apache License, Version 2.0

Author Information
------------------
