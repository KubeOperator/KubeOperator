class ClusterHealthData():

    def __init__(self, namespace, name, status, ready, age):
        self.namespace = namespace
        self.name = name
        self.status = status
        self.ready = ready
        self.age = age
