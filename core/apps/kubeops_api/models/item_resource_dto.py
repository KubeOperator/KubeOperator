class ItemResourceDTO():

    def __init__(self, item_resource, item_name):
        self.id = item_resource.id
        self.resource_id = item_resource.resource_id
        self.item_id = item_resource.item_id
        self.item_name = item_name


class Resource():

    def __init__(self, resource_id, resource_type, data, checked):
        self.resource_id = resource_id
        self.resource_type = resource_type
        self.data = data
        self.checked = checked
