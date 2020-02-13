class ItemResourceDTO():

    def __init__(self, item_resource, resource, checked):
        self.id = item_resource.id
        self.resource_id = item_resource.resource_id
        self.item_id = item_resource.item_id
        self.resource_type = item_resource.resource_type
        self.resource = resource
        self.checked = checked

class Resource():

    def __init__(self,resource_id,resource_type,data,checked):
        self.resource_id = resource_id
        self.resource_type = resource_type
        self.data = data
        self.checked = checked



