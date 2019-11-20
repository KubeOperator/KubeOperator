class ClusterHealthData():

    def __init__(self, namespace, name, status, ready, age, msg, restart_count):
        self.namespace = namespace
        self.name = name
        self.status = status
        self.ready = ready
        self.age = age
        self.msg = msg
        self.restart_count = restart_count
