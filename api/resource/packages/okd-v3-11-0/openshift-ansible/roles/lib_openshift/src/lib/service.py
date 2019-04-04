# pylint: skip-file
# flake8: noqa


# pylint: disable=too-many-instance-attributes
class ServiceConfig(object):
    ''' Handle service options '''
    # pylint: disable=too-many-arguments
    def __init__(self,
                 sname,
                 namespace,
                 ports,
                 annotations=None,
                 selector=None,
                 labels=None,
                 cluster_ip=None,
                 portal_ip=None,
                 session_affinity=None,
                 service_type=None,
                 external_ips=None):
        ''' constructor for handling service options '''
        self.name = sname
        self.namespace = namespace
        self.ports = ports
        self.annotations = annotations
        self.selector = selector
        self.labels = labels
        self.cluster_ip = cluster_ip
        self.portal_ip = portal_ip
        self.session_affinity = session_affinity
        self.service_type = service_type
        self.external_ips = external_ips
        self.data = {}

        self.create_dict()

    def create_dict(self):
        ''' instantiates a service dict '''
        self.data['apiVersion'] = 'v1'
        self.data['kind'] = 'Service'
        self.data['metadata'] = {}
        self.data['metadata']['name'] = self.name
        self.data['metadata']['namespace'] = self.namespace
        if self.labels:
            self.data['metadata']['labels'] = {}
            for lab, lab_value in self.labels.items():
                self.data['metadata']['labels'][lab] = lab_value
        if self.annotations:
            self.data['metadata']['annotations'] = self.annotations

        self.data['spec'] = {}

        if self.ports:
            self.data['spec']['ports'] = self.ports
        else:
            self.data['spec']['ports'] = []

        if self.selector:
            self.data['spec']['selector'] = self.selector

        self.data['spec']['sessionAffinity'] = self.session_affinity or 'None'

        if self.cluster_ip:
            self.data['spec']['clusterIP'] = self.cluster_ip

        if self.portal_ip:
            self.data['spec']['portalIP'] = self.portal_ip

        if self.service_type:
            self.data['spec']['type'] = self.service_type

        if self.external_ips:
            self.data['spec']['externalIPs'] = self.external_ips


# pylint: disable=too-many-instance-attributes,too-many-public-methods
class Service(Yedit):
    ''' Class to model the oc service object '''
    port_path = "spec.ports"
    portal_ip = "spec.portalIP"
    cluster_ip = "spec.clusterIP"
    selector_path = 'spec.selector'
    kind = 'Service'
    external_ips = "spec.externalIPs"

    def __init__(self, content):
        '''Service constructor'''
        super(Service, self).__init__(content=content)

    def get_ports(self):
        ''' get a list of ports '''
        return self.get(Service.port_path) or []

    def get_selector(self):
        ''' get the service selector'''
        return self.get(Service.selector_path) or {}

    def add_ports(self, inc_ports):
        ''' add a port object to the ports list '''
        if not isinstance(inc_ports, list):
            inc_ports = [inc_ports]

        ports = self.get_ports()
        if not ports:
            self.put(Service.port_path, inc_ports)
        else:
            ports.extend(inc_ports)

        return True

    def find_ports(self, inc_port):
        ''' find a specific port '''
        for port in self.get_ports():
            if port['port'] == inc_port['port']:
                return port

        return None

    def delete_ports(self, inc_ports):
        ''' remove a port from a service '''
        if not isinstance(inc_ports, list):
            inc_ports = [inc_ports]

        ports = self.get(Service.port_path) or []

        if not ports:
            return True

        removed = False
        for inc_port in inc_ports:
            port = self.find_ports(inc_port)
            if port:
                ports.remove(port)
                removed = True

        return removed

    def add_cluster_ip(self, sip):
        '''add cluster ip'''
        self.put(Service.cluster_ip, sip)

    def add_portal_ip(self, pip):
        '''add cluster ip'''
        self.put(Service.portal_ip, pip)

    def get_external_ips(self):
        ''' get a list of external_ips '''
        return self.get(Service.external_ips) or []

    def add_external_ips(self, inc_external_ips):
        ''' add an external_ip to the external_ips list '''
        if not isinstance(inc_external_ips, list):
            inc_external_ips = [inc_external_ips]

        external_ips = self.get_external_ips()
        if not external_ips:
            self.put(Service.external_ips, inc_external_ips)
        else:
            external_ips.extend(inc_external_ips)

        return True

    def find_external_ips(self, inc_external_ip):
        ''' find a specific external IP '''
        val = None
        try:
            idx = self.get_external_ips().index(inc_external_ip)
            val = self.get_external_ips()[idx]
        except ValueError:
            pass

        return val

    def delete_external_ips(self, inc_external_ips):
        ''' remove an external IP from a service '''
        if not isinstance(inc_external_ips, list):
            inc_external_ips = [inc_external_ips]

        external_ips = self.get(Service.external_ips) or []

        if not external_ips:
            return True

        removed = False
        for inc_external_ip in inc_external_ips:
            external_ip = self.find_external_ips(inc_external_ip)
            if external_ip:
                external_ips.remove(external_ip)
                removed = True

        return removed
