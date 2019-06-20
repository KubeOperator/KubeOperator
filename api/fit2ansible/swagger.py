from drf_yasg.inspectors import SwaggerAutoSchema


class CustomSwaggerAutoSchema(SwaggerAutoSchema):
    def get_tags(self, operation_keys):
        if operation_keys[0] == 'projects' and len(operation_keys) >= 3:
            if operation_keys[1] == 'inventory' and len(operation_keys) > 3:
                value = operation_keys[2]
            elif operation_keys[1] in ('adhoc', 'playbooks') and len(operation_keys) > 3:
                value = operation_keys[1] + '-' + operation_keys[2]
            else:
                value = operation_keys[1]
            return ['project-' + value]
        elif operation_keys[0] == 'clusters' and len(operation_keys) >= 3:
            return ['cluster-' + operation_keys[1]]
        elif operation_keys[0] == 'inventory' and len(operation_keys) >= 3:
            return [operation_keys[1]]
        elif operation_keys[0] == 'storage' and len(operation_keys) >= 3:
            return ['storage-' + operation_keys[1]]
        return super().get_tags(operation_keys)
