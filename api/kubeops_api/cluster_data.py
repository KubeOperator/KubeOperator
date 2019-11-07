import json
from uuid import UUID


class ClusterData():

    def __init__(self, cluster, token, pods, nodes, namespaces, deployments, cpu_usage, cpu_total, mem_usage,
                 mem_total, restart_pods, warn_containers,error_loki_containers):
        self.id = str(cluster.id)
        self.name = cluster.name
        self.pods = pods
        self.nodes = nodes
        self.token = token
        self.namespaces = namespaces
        self.deployments = deployments
        self.cpu_usage = cpu_usage
        self.cpu_total = cpu_total
        self.mem_usage = mem_usage
        self.mem_total = mem_total
        self.restart_pods = restart_pods
        self.warn_containers = warn_containers
        self.error_loki_containers = error_loki_containers


class Pod():

    def __init__(self, name, cluster_name, restart_count, status, namespace, host_ip, pod_ip, host_name, containers):
        self.name = name
        self.cluster_name = cluster_name
        self.restart_count = restart_count
        self.status = status
        self.namespace = namespace
        self.host_ip = host_ip
        self.pod_ip = pod_ip
        self.host_name = host_name
        self.containers = containers


class NameSpace():

    def __init__(self, name, status):
        self.name = name
        self.status = status


class Node():

    def __init__(self, name, status, cpu, mem, cpu_usage, mem_usage):
        self.name = name
        self.status = status
        self.cpu = cpu
        self.mem = mem
        self.cpu_usage = cpu_usage
        self.mem_usage = mem_usage


class Container():

    def __init__(self, name, ready, restart_count, pod_name):
        self.name = name
        self.ready = ready
        self.restart_count = restart_count
        self.pod_name = pod_name


class Deployment():

    def __init__(self, name, ready_replicas, replicas, namespace):
        self.name = name
        self.ready_replicas = ready_replicas
        self.replicas = replicas
        self.namespace = namespace


class LokiContainer():

    def __init__(self, name, error_count,cluster_name):
        self.name = name
        self.error_count = error_count
        self.cluster_name = cluster_name
