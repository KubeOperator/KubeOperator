# Build AMI

When seeking to deploy a working openshift cluster using these plays, a few
items must be in place.

These are:

1. Create an instance, using a specified ssh key.
2. Run openshift-ansible setup roles to ensure packages and services are correctly configured.
3. Create the AMI.
4. If encryption is desired
  - A KMS key is created with the name of $clusterid
  - An encrypted AMI will be produced with $clusterid KMS key
5. Terminate the instance used to configure the AMI.

More AMI specific options can be found in ['openshift_aws/defaults/main.yml'](../../roles/openshift_aws/defaults/main.yml).  When creating an encrypted AMI please specify use_encryption:
```
# openshift_aws_ami_encrypt: True  # defaults to false
```

**Note**:  This will ensure to take the recently created AMI and encrypt it to be used later.  If encryption is not desired then set the value to false (defaults to false). The AMI id will be fetched and used according to its most recent creation date.
