from ansible_api.models import Group


class Role(Group):
    class Meta:
        proxy = True

    @property
    def nodes(self):
        return self.hosts

    @nodes.setter
    def nodes(self, value):
        self.hosts.set(value)

    def __str__(self):
        return "%s %s" % (self.project, self.name)
