'''
---
module: openshift_logging_facts
version_added: ""
short_description: Gather facts about the OpenShift logging stack
description:
  - Determine the current facts about the OpenShift logging stack (e.g. cluster size)
options:
author: Red Hat, Inc
'''

import copy
import json

# pylint: disable=redefined-builtin, unused-wildcard-import, wildcard-import
from subprocess import *   # noqa: F402,F403

# ignore pylint errors related to the module_utils import
# pylint: disable=redefined-builtin, unused-wildcard-import, wildcard-import
from ansible.module_utils.basic import *   # noqa: F402,F403

import yaml

EXAMPLES = """
- action: opneshift_logging_facts
"""

RETURN = """
"""

DEFAULT_OC_OPTIONS = ["-o", "json"]

# constants used for various labels and selectors
COMPONENT_KEY = "component"
LOGGING_INFRA_KEY = "logging-infra"

# selectors for filtering resources
DS_FLUENTD_SELECTOR = LOGGING_INFRA_KEY + "=" + "fluentd"
LOGGING_SELECTOR = LOGGING_INFRA_KEY + "=" + "support"
ROUTE_SELECTOR = "component=support,logging-infra=support,provider=openshift"
# pylint: disable=line-too-long
COMPONENTS = ["kibana", "curator", "elasticsearch", "fluentd", "kibana_ops", "curator_ops", "elasticsearch_ops", "mux", "eventrouter"]


class OCBaseCommand(object):
    ''' The base class used to query openshift '''

    def __init__(self, binary, kubeconfig, namespace):
        ''' the init method of OCBaseCommand class '''
        self.binary = binary
        self.kubeconfig = kubeconfig
        self.user = self.get_system_admin(self.kubeconfig)
        self.namespace = namespace

    # pylint: disable=no-self-use
    def get_system_admin(self, kubeconfig):
        ''' Retrieves the system admin '''
        with open(kubeconfig, 'r') as kubeconfig_file:
            config = yaml.load(kubeconfig_file)
            for user in config["users"]:
                if user["name"].startswith("system:admin"):
                    return user["name"]
        raise Exception("Unable to find system:admin in: " + kubeconfig)

    # pylint: disable=too-many-arguments, dangerous-default-value
    def oc_command(self, sub, kind, namespace=None, name=None, add_options=None):
        ''' Wrapper method for the "oc" command '''
        cmd = [self.binary, sub, kind]
        if name is not None:
            cmd = cmd + [name]
        if namespace is not None:
            cmd = cmd + ["-n", namespace]
        if add_options is None:
            add_options = []
        cmd = cmd + ["--user=" + self.user, "--config=" + self.kubeconfig] + DEFAULT_OC_OPTIONS + add_options
        try:
            process = Popen(cmd, stdout=PIPE, stderr=PIPE)   # noqa: F405
            out, err = process.communicate(cmd)
            err = err.decode(encoding='utf8', errors='replace')
            if len(err) > 0:
                if 'not found' in err:
                    return {'items': []}
                if 'No resources found' in err:
                    return {'items': []}
                raise Exception(err)
        except Exception as excp:
            err = "There was an exception trying to run the command '" + " ".join(cmd) + "' " + str(excp)
            raise Exception(err)

        return json.loads(out)


class OpenshiftLoggingFacts(OCBaseCommand):
    ''' The class structure for holding the OpenshiftLogging Facts'''
    name = "facts"

    def __init__(self, logger, binary, kubeconfig, namespace):
        ''' The init method for OpenshiftLoggingFacts '''
        super(OpenshiftLoggingFacts, self).__init__(binary, kubeconfig, namespace)
        self.logger = logger
        self.facts = dict()

    def default_keys_for(self, kind):
        ''' Sets the default key values for kind '''
        for comp in COMPONENTS:
            self.add_facts_for(comp, kind)

    def add_facts_for(self, comp, kind, name=None, facts=None):
        ''' Add facts for the provided kind '''
        if comp not in self.facts:
            self.facts[comp] = dict()
        if kind not in self.facts[comp]:
            self.facts[comp][kind] = dict()
        if name:
            self.facts[comp][kind][name] = facts

    def facts_for_routes(self, namespace):
        ''' Gathers facts for Routes in logging namespace '''
        self.default_keys_for("routes")
        route_list = self.oc_command("get", "routes", namespace=namespace, add_options=["-l", ROUTE_SELECTOR])
        if len(route_list["items"]) == 0:
            return None
        for route in route_list["items"]:
            name = route["metadata"]["name"]
            comp = self.comp(name)
            if comp is not None:
                self.add_facts_for(comp, "routes", name, dict(host=route["spec"]["host"]))
        self.facts["agl_namespace"] = namespace

    def facts_for_daemonsets(self, namespace):
        ''' Gathers facts for Daemonsets in logging namespace '''
        self.default_keys_for("daemonsets")
        ds_list = self.oc_command("get", "daemonsets", namespace=namespace,
                                  add_options=["-l", LOGGING_INFRA_KEY + "=fluentd"])
        if len(ds_list["items"]) == 0:
            return
        for ds_item in ds_list["items"]:
            name = ds_item["metadata"]["name"]
            comp = self.comp(name)
            spec = ds_item["spec"]["template"]["spec"]
            result = dict(
                selector=ds_item["spec"]["selector"],
                containers=dict(),
                nodeSelector=spec["nodeSelector"],
                serviceAccount=spec["serviceAccount"],
                terminationGracePeriodSeconds=spec["terminationGracePeriodSeconds"]
            )
            for container in spec["containers"]:
                result["containers"][container["name"]] = container
            self.add_facts_for(comp, "daemonsets", name, result)

    def facts_for_pvcs(self, namespace):
        ''' Gathers facts for PVCS in logging namespace'''
        self.default_keys_for("pvcs")
        pvclist = self.oc_command("get", "pvc", namespace=namespace, add_options=["-l", LOGGING_INFRA_KEY])
        if len(pvclist["items"]) == 0:
            return
        for pvc in pvclist["items"]:
            name = pvc["metadata"]["name"]
            comp = self.comp(name)
            self.add_facts_for(comp, "pvcs", name, dict())

    def facts_for_deploymentconfigs(self, namespace):
        ''' Gathers facts for DeploymentConfigs in logging namespace '''
        self.default_keys_for("deploymentconfigs")
        dclist = self.oc_command("get", "deploymentconfigs", namespace=namespace, add_options=["-l", LOGGING_INFRA_KEY])
        if len(dclist["items"]) == 0:
            return
        dcs = dclist["items"]
        for dc_item in dcs:
            name = dc_item["metadata"]["name"]
            comp = self.comp(name)
            if comp is not None:
                spec = dc_item["spec"]["template"]["spec"]
                facts = dict(
                    name=name,
                    selector=dc_item["spec"]["selector"],
                    replicas=dc_item["spec"]["replicas"],
                    serviceAccount=spec["serviceAccount"],
                    containers=dict(),
                    volumes=dict()
                )
                if "nodeSelector" in spec:
                    facts["nodeSelector"] = spec["nodeSelector"]
                if "supplementalGroups" in spec["securityContext"]:
                    facts["storageGroups"] = spec["securityContext"]["supplementalGroups"]
                facts["spec"] = spec
                if "volumes" in spec:
                    for vol in spec["volumes"]:
                        clone = copy.deepcopy(vol)
                        clone.pop("name", None)
                        facts["volumes"][vol["name"]] = clone
                for container in spec["containers"]:
                    facts["containers"][container["name"]] = container
                self.add_facts_for(comp, "deploymentconfigs", name, facts)

    def facts_for_services(self, namespace):
        ''' Gathers facts for services in logging namespace '''
        self.default_keys_for("services")
        servicelist = self.oc_command("get", "services", namespace=namespace, add_options=["-l", LOGGING_SELECTOR])
        if len(servicelist["items"]) == 0:
            return
        for service in servicelist["items"]:
            name = service["metadata"]["name"]
            comp = self.comp(name)
            if comp is not None:
                self.add_facts_for(comp, "services", name, dict())

    # pylint: disable=too-many-arguments
    def facts_from_configmap(self, comp, kind, name, config_key, yaml_file=None):
        '''Extracts facts in logging namespace from configmap'''
        if yaml_file is not None:
            if config_key.endswith(".yml") or config_key.endswith(".yaml"):
                config_facts = yaml.load(yaml_file)
                self.facts[comp][kind][name][config_key] = config_facts
                self.facts[comp][kind][name][config_key]["raw"] = yaml_file

    def facts_for_configmaps(self, namespace):
        ''' Gathers facts for configmaps in logging namespace '''
        self.default_keys_for("configmaps")
        a_list = self.oc_command("get", "configmaps", namespace=namespace)
        if len(a_list["items"]) == 0:
            return
        for item in a_list["items"]:
            name = item["metadata"]["name"]
            comp = self.comp(name)
            if comp is not None:
                self.add_facts_for(comp, "configmaps", name, dict(item["data"]))
                if comp in ["elasticsearch", "elasticsearch_ops"]:
                    for config_key in item["data"]:
                        self.facts_from_configmap(comp, "configmaps", name, config_key, item["data"][config_key])

    def facts_for_oauthclients(self, namespace):
        ''' Gathers facts for oauthclients used with logging '''
        self.default_keys_for("oauthclients")
        a_list = self.oc_command("get", "oauthclients", namespace=namespace, add_options=["-l", LOGGING_SELECTOR])
        if len(a_list["items"]) == 0:
            return
        for item in a_list["items"]:
            name = item["metadata"]["name"]
            comp = self.comp(name)
            if comp is not None:
                result = dict(
                    redirectURIs=item["redirectURIs"]
                )
                self.add_facts_for(comp, "oauthclients", name, result)

    def facts_for_secrets(self, namespace):
        ''' Gathers facts for secrets in the logging namespace '''
        self.default_keys_for("secrets")
        a_list = self.oc_command("get", "secrets", namespace=namespace)
        if len(a_list["items"]) == 0:
            return
        for item in a_list["items"]:
            name = item["metadata"]["name"]
            comp = self.comp(name)
            if comp is not None and item["type"] == "Opaque":
                result = dict(
                    keys=item["data"].keys()
                )
                self.add_facts_for(comp, "secrets", name, result)

    def facts_for_sccs(self):
        ''' Gathers facts for SCCs used with logging '''
        self.default_keys_for("sccs")
        scc = self.oc_command("get", "securitycontextconstraints.v1.security.openshift.io", name="privileged")
        if len(scc["users"]) == 0:
            return
        for item in scc["users"]:
            comp = self.comp(item)
            if comp is not None:
                self.add_facts_for(comp, "sccs", "privileged", dict())

    def facts_for_clusterrolebindings(self, namespace):
        ''' Gathers ClusterRoleBindings used with logging '''
        self.default_keys_for("clusterrolebindings")
        role = self.oc_command("get", "clusterrolebindings", name="cluster-readers")
        if "subjects" not in role or len(role["subjects"]) == 0:
            return
        for item in role["subjects"]:
            comp = self.comp(item["name"])
            if comp is not None and namespace == item.get("namespace"):
                self.add_facts_for(comp, "clusterrolebindings", "cluster-readers", dict())

# this needs to end up nested under the service account...
    def facts_for_rolebindings(self, namespace):
        ''' Gathers facts for RoleBindings used with logging '''
        self.default_keys_for("rolebindings")
        role = self.oc_command("get", "rolebindings", namespace=namespace, name="logging-elasticsearch-view-role")
        if "subjects" not in role or len(role["subjects"]) == 0:
            return
        for item in role["subjects"]:
            comp = self.comp(item["name"])
            if comp is not None and namespace == item.get("namespace"):
                self.add_facts_for(comp, "rolebindings", "logging-elasticsearch-view-role", dict())

    # pylint: disable=no-self-use, too-many-return-statements
    def comp(self, name):
        ''' Does a comparison to evaluate the logging component '''
        if name.startswith("logging-curator-ops"):
            return "curator_ops"
        elif name.startswith("logging-kibana-ops") or name.startswith("kibana-ops"):
            return "kibana_ops"
        elif name.startswith("logging-es-ops") or name.startswith("logging-elasticsearch-ops"):
            return "elasticsearch_ops"
        elif name.startswith("logging-curator"):
            return "curator"
        elif name.startswith("logging-kibana") or name.startswith("kibana"):
            return "kibana"
        elif name.startswith("logging-es") or name.startswith("logging-elasticsearch"):
            return "elasticsearch"
        elif name.startswith("logging-fluentd") or name.endswith("aggregated-logging-fluentd"):
            return "fluentd"
        elif name.startswith("logging-mux"):
            return "mux"
        elif name.startswith("logging-eventrouter"):
            return "eventrouter"
        else:
            return None

    def build_facts(self):
        ''' Builds the logging facts and returns them '''
        self.facts_for_routes(self.namespace)
        self.facts_for_daemonsets(self.namespace)
        self.facts_for_deploymentconfigs(self.namespace)
        self.facts_for_services(self.namespace)
        self.facts_for_configmaps(self.namespace)
        self.facts_for_sccs()
        self.facts_for_oauthclients(self.namespace)
        self.facts_for_clusterrolebindings(self.namespace)
        self.facts_for_rolebindings(self.namespace)
        self.facts_for_secrets(self.namespace)
        self.facts_for_pvcs(self.namespace)

        return self.facts


def main():
    ''' The main method '''
    module = AnsibleModule(   # noqa: F405
        argument_spec=dict(
            admin_kubeconfig={"default": "/etc/origin/master/admin.kubeconfig", "type": "str"},
            oc_bin={"required": True, "type": "str"},
            openshift_logging_namespace={"required": True, "type": "str"}
        ),
        supports_check_mode=False
    )
    try:
        cmd = OpenshiftLoggingFacts(module, module.params['oc_bin'], module.params['admin_kubeconfig'],
                                    module.params['openshift_logging_namespace'])
        module.exit_json(
            ansible_facts={"openshift_logging_facts": cmd.build_facts()}
        )
    # ignore broad-except error to avoid stack trace to ansible user
    # pylint: disable=broad-except
    except Exception as error:
        module.fail_json(msg=str(error))


if __name__ == '__main__':
    main()
