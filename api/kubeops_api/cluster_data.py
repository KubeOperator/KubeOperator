
import json
from uuid import UUID


class ClusterData():

    def __init__(self,cluster,token,pods,nodes,name_spaces):
        self.id = str(cluster.id)
        self.name = cluster.name
        self.pods = pods
        self.nodes = nodes
        self.token = token
        self.name_spaces =name_spaces


class Pod():

    def __init__(self,name,cluster_name,restart_count,status,name_space,host_ip,pod_ip,host_name):
        self.name = name
        self.cluster_name = cluster_name
        self.restart_count = restart_count
        self.status = status
        self.name_space = name_space
        self.host_ip = host_ip
        self.pod_ip = pod_ip
        self.host_name = host_name

class NameSpace():

    def __init__(self,name,status):
        self.name = name
        self.status = status


class Node():

    def __init__(self,name,status):
        self.name = name
        self.status = status
