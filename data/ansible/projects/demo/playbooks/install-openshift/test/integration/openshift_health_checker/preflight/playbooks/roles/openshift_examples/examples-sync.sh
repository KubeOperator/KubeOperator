#!/bin/bash

# Utility script to update the ansible repo with the latest templates and image
# streams from several github repos
#
# This script should be run from openshift-ansible/roles/openshift_examples

XPAAS_VERSION=ose-v1.4.14
RHDM70_VERSION=ose-v1.4.8-1
RHPAM70_VERSION=7.0.0.GA
DG_72_VERSION=1.1.1
ORIGIN_VERSION=${1:-v3.10}
ORIGIN_BRANCH=${2:-release-3.10}
RHAMP_TAG=2.0.0.GA
EXAMPLES_BASE=$(pwd)/files/examples/${ORIGIN_VERSION}
find ${EXAMPLES_BASE} -name '*.json' -delete
TEMP=`mktemp -d`
pushd $TEMP

if [ ! -d "${EXAMPLES_BASE}" ]; then
  mkdir -p ${EXAMPLES_BASE}
fi
wget https://github.com/openshift/origin/archive/${ORIGIN_BRANCH}.zip -O origin.zip
wget https://github.com/jboss-fuse/application-templates/archive/GA.zip -O fis-GA.zip
wget https://github.com/jboss-openshift/application-templates/archive/${XPAAS_VERSION}.zip -O application-templates-master.zip
wget https://github.com/jboss-container-images/rhdm-7-openshift-image/archive/${RHDM70_VERSION}.zip -O rhdm-application-templates.zip
wget https://github.com/jboss-container-images/rhpam-7-openshift-image/archive/${RHPAM70_VERSION}.zip -O rhpam-application-templates.zip
wget https://github.com/3scale/rhamp-openshift-templates/archive/${RHAMP_TAG}.zip -O amp.zip
wget https://github.com/jboss-container-images/jboss-datagrid-7-openshift-image/archive/${DG_72_VERSION}.zip -O dg-application-templates.zip
unzip origin.zip
unzip application-templates-master.zip
unzip rhdm-application-templates.zip
unzip rhpam-application-templates.zip
unzip fis-GA.zip
unzip amp.zip
unzip dg-application-templates.zip
mv origin-${ORIGIN_BRANCH}/examples/db-templates/* ${EXAMPLES_BASE}/db-templates/
mv origin-${ORIGIN_BRANCH}/examples/quickstarts/* ${EXAMPLES_BASE}/quickstart-templates/
mv origin-${ORIGIN_BRANCH}/examples/jenkins/jenkins-*template.json ${EXAMPLES_BASE}/quickstart-templates/
mv origin-${ORIGIN_BRANCH}/examples/image-streams/* ${EXAMPLES_BASE}/image-streams/
mv application-templates-${XPAAS_VERSION}/jboss-image-streams.json ${EXAMPLES_BASE}/xpaas-streams/
mv rhdm-7-openshift-image-${RHDM70_VERSION}/rhdm70-image-streams.yaml ${EXAMPLES_BASE}/xpaas-streams/
mv rhpam-7-openshift-image-${RHPAM70_VERSION}/rhpam70-image-streams.yaml ${EXAMPLES_BASE}/xpaas-streams/
mv jboss-datagrid-7-openshift-image-${DG_72_VERSION}/templates/datagrid72-image-stream.json ${EXAMPLES_BASE}/xpaas-streams/
# fis content from jboss-fuse/application-templates-GA would collide with jboss-openshift/application-templates
# as soon as they use the same branch/tag names
mv application-templates-GA/fis-image-streams.json ${EXAMPLES_BASE}/xpaas-streams/fis-image-streams.json
mv application-templates-GA/quickstarts/* ${EXAMPLES_BASE}/xpaas-templates/
mv application-templates-GA/fis-console-namespace-template.json application-templates-GA/fis-console-cluster-template.json ${EXAMPLES_BASE}/xpaas-templates/
find application-templates-${XPAAS_VERSION}/ -name '*.json' ! -wholename '*secret*' ! -wholename '*demo*' ! -wholename '*image-stream.json' -exec mv {} ${EXAMPLES_BASE}/xpaas-templates/ \;
find application-templates-${XPAAS_VERSION}/ -name '*image-stream.json' -exec mv {} ${EXAMPLES_BASE}/xpaas-streams/ \;
find rhdm-7-openshift-image-${RHDM70_VERSION}/templates -name '*.yaml' -exec mv {} ${EXAMPLES_BASE}/xpaas-templates/ \;
find rhpam-7-openshift-image-${RHPAM70_VERSION}/templates -name '*.yaml' -exec mv {} ${EXAMPLES_BASE}/xpaas-templates/ \;
find 3scale-amp-openshift-templates-${RHAMP_TAG}/ -name '*.yml' -exec mv {} ${EXAMPLES_BASE}/quickstart-templates/ \;
find jboss-datagrid-7-openshift-image-${DG_72_VERSION}/templates/ -name '*.json' -exec mv {} ${EXAMPLES_BASE}/xpaas-templates/ \;
popd

wget https://raw.githubusercontent.com/redhat-developer/s2i-dotnetcore/master/dotnet_imagestreams.json         -O ${EXAMPLES_BASE}/image-streams/dotnet_imagestreams.json
wget https://raw.githubusercontent.com/redhat-developer/s2i-dotnetcore/master/dotnet_imagestreams_centos.json         -O ${EXAMPLES_BASE}/image-streams/dotnet_imagestreams_centos.json
wget https://raw.githubusercontent.com/redhat-developer/s2i-dotnetcore/master/templates/dotnet-example.json           -O ${EXAMPLES_BASE}/quickstart-templates/dotnet-example.json
wget https://raw.githubusercontent.com/redhat-developer/s2i-dotnetcore/master/templates/dotnet-pgsql-persistent.json    -O ${EXAMPLES_BASE}/quickstart-templates/dotnet-pgsql-persistent.json
wget https://raw.githubusercontent.com/redhat-developer/s2i-dotnetcore/master/templates/dotnet-runtime-example.json    -O ${EXAMPLES_BASE}/quickstart-templates/dotnet-runtime-example.json

git diff files/examples
