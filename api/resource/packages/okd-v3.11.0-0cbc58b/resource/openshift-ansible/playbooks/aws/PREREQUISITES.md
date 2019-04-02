# Prerequisites

When seeking to deploy a working openshift cluster using these plays, a few
items must be in place.

These are:

1) vpc
2) security group to build the AMI in.
3) ssh keys to log into instances

These items can be provisioned ahead of time, or you can utilize the plays here
to create these items.

If you wish to provision these items yourself, or you already have these items
provisioned and wish to utilize existing components, please refer to
provisioning_vars.yml.example.

If you wish to have these items created for you, continue with this document.

# Running prerequisites.yml

Warning:  Running these plays will provision items in your AWS account (if not
present), and you may incur billing charges.  These plays are not suitable
for the free-tier.

## Step 1:
Ensure you have specified all the necessary provisioning variables.  See
provisioning_vars.example.yml and README.md for more information.

## Step 2:
```
$ ansible-playbook -i inventory.yml prerequisites.yml -e @provisioning_vars.yml
```

This will create a VPC, security group, and ssh_key.  These plays are idempotent,
and multiple runs should result in no additional provisioning of these components.

You can also verify that you will successfully utilize existing components with
these plays.
