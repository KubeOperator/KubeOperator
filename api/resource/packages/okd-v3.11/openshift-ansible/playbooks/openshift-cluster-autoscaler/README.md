# Cluster autoscaler

This directory contains Ansible playbooks for deployment of [cluster autoscaler](https://github.com/openshift/kubernetes-autoscaler).

It's assumed a user is familiar with the cluster deployment.
For cluster deployment in AWS cloud provider you can follow instructions provided in [AWS playbooks](https://github.com/openshift/openshift-ansible/tree/master/playbooks/aws).

## Running

The autoscaler can be deployed directly by running the corresponding playbook:

```sh
ansible-playbook -i inventory playbooks/openshift-cluster-autoscaler/config.yml
```

or during your cluster deployment by setting `openshift_cluster_autoscaler_install=True`
in the corresponding inventory file.

Currently, the autoscaler is deployed only on nodes that have the `infra` [node role](https://github.com/openshift/openshift-ansible#node-group-definition-and-mapping) assigned (forced by node selector set to `node-role.kubernetes.io/infra: "true"`, see `openshift_cluster_autoscaler_node_selector` variable of the `openshift_cluster_autoscaler` role).
In order to run the autoscaler inside the AWS cloud provider,
one has to attach an [IAM role](https://aws.amazon.com/iam/details/manage-roles/) to the `infra` node that has appropriate permissions to scale down nodes in  [Auto Scaling Group](https://docs.aws.amazon.com/autoscaling/ec2/userguide/AutoScalingGroup.html) belonging to the OpenShift compute node group.

In case there is only one `infra` node, the process of attaching the IAM role
reduces to running:
* assuming the AWS credentials are available
* assuming the IAM role with appropriate permissions is named `aws-autoscaler-role`
* region set to `us-east-1`
* with the infra node name set to `CLUSTERID infra group 1`

```sh
# to install the aws executable
sudo pip install awscli
# Get infra node instance ID (assuming there is only one)
INFRA_NODE_INSTANCE_ID=$(aws --region us-east-1 ec2 describe-instances --filters 'Name=tag:openshift-node-group-config,Values=node-config-infra' "Name=tag:Name,Values=CLUSTERID infra group 1" | jq '.Reservations[0].Instances[0].InstanceId' --raw-output)
# Get AssociationId for the AWS instance
ASSOCIATION_ID=$(aws --region us-east-1 ec2 describe-iam-instance-profile-associations --filters "Name=instance-id,Values=${INFRA_NODE_INSTANCE_ID}" | jq '.IamInstanceProfileAssociations[0].AssociationId' --raw-output)
# Attach the IAM role to the infra instance
aws --region us-east-1 ec2 replace-iam-instance-profile-association --association-id ${ASSOCIATION_ID} --iam-instance-profile Name=aws-autoscaler-role
```

## Testing

To exercise the basic autoscaling capability, one can create a workload in the same namespace as the autoscaler lives
and observe creation of new nodes. Followed by deleting the workload and observing nodes being removed from the cluster.

**Workload example**:
```yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: scale-up
  labels:
    app: scale-up
spec:
  replicas: 50
  selector:
    matchLabels:
      app: scale-up
  template:
    metadata:
      labels:
        app: scale-up
    spec:
      containers:
      - name: busybox
        image: docker.io/library/busybox
        resources:
          requests:
            memory: 2Gi
        command:
        - /bin/sh
        - "-c"
        - "echo 'this should be in the logs' && sleep 86400"
      terminationGracePeriodSeconds: 0
```

Assuming the cluster autoscaler is deployed in `openshift-autoscaler` namespace,
the following command will create the workload:

```sh
oc create -n openshift-autoscaler -f scale-up.yaml
```

To observe the cluster autoscaler's behavior closely, one can list its logs by running:

```sh
oc logs cluster-autoscaler-HASH
```

The pod name can be listed by running:

```sh
oc get pods --namespace openshift-autoscaler
```

To observe registration of nodes and their current state, one can run the following command
on the master node:

```sh
watch -n 1 oc get nodes
```
