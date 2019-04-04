# OpenShift External Dynamic Provisioners

## Required Vars
* `openshift_provisioners_install_provisioners`: When `True` the openshift_provisioners role will install provisioners that have their "master" var (e.g. `openshift_provisioners_efs`) set `True`. When `False` will uninstall provisioners that have their var set `True`.

## Optional Vars
* `openshift_provisioners_image_prefix`: The prefix for the provisioner images to use. Defaults to 'docker.io/openshift/origin-'.
* `openshift_provisioners_image_version`: The image version for the provisioner images to use. Defaults to 'latest'.
* `openshift_provisioners_project`: The namespace that provisioners will be installed in. Defaults to 'openshift-infra'.

## AWS EFS

### Prerequisites
* An IAM user assigned the AmazonElasticFileSystemReadOnlyAccess policy (or better)
* An EFS file system in your cluster's region
* [Mount targets](http://docs.aws.amazon.com/efs/latest/ug/accessing-fs.html) and [security groups](http://docs.aws.amazon.com/efs/latest/ug/accessing-fs-create-security-groups.html) such that any node (in any zone in the cluster's region) can mount the EFS file system by its [File system DNS name](http://docs.aws.amazon.com/efs/latest/ug/mounting-fs-mount-cmd-dns-name.html)

### Required Vars
* `openshift_provisioners_efs_fsid`: The [File system ID](http://docs.aws.amazon.com/efs/latest/ug/gs-step-two-create-efs-resources.html) of the EFS file system, e.g. fs-47a2c22e.
* `openshift_provisioners_efs_region`: The Amazon EC2 region of the EFS file system.
* `openshift_provisioners_efs_aws_access_key_id`: The AWS access key of the IAM user, used to check that the EFS file system specified actually exists.
* `openshift_provisioners_efs_aws_secret_access_key`: The AWS secret access key of the IAM user, used to check that the EFS file system specified actually exists.

### Optional Vars
* `openshift_provisioners_efs`: When `True` the AWS EFS provisioner will be installed or uninstalled according to whether `openshift_provisioners_install_provisioners` is `True` or `False`, respectively. Defaults to `False`.
* `openshift_provisioners_efs_path`: The path of the directory in the EFS file system in which the EFS provisioner will create a directory to back each PV it creates. It must exist and be mountable by the EFS provisioner. Defaults to '/persistentvolumes'.
* `openshift_provisioners_efs_name`: The `provisioner` name that `StorageClasses` specify. Defaults to 'openshift.org/aws-efs'.
* `openshift_provisioners_efs_nodeselector`: A map of labels (e.g. {"node":"infra","region":"west"} to select the nodes where the pod will land.
* `openshift_provisioners_efs_supplementalgroup`: The supplemental group to give the pod in case it is needed for permission to write to the EFS file system. Defaults to '65534'.
