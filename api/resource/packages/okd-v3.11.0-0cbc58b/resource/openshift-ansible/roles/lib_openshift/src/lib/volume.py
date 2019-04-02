# pylint: skip-file
# flake8: noqa

class Volume(object):
    ''' Class to represent an openshift volume object'''
    volume_mounts_path = {"pod": "spec.containers[0].volumeMounts",
                          "dc":  "spec.template.spec.containers[0].volumeMounts",
                          "rc":  "spec.template.spec.containers[0].volumeMounts",
                         }
    volumes_path = {"pod": "spec.volumes",
                    "dc":  "spec.template.spec.volumes",
                    "rc":  "spec.template.spec.volumes",
                   }

    @staticmethod
    def create_volume_structure(volume_info):
        ''' return a properly structured volume '''
        volume_mount = None
        volume = {'name': volume_info['name']}
        volume_type = volume_info['type'].lower()
        if volume_type == 'secret':
            volume['secret'] = {}
            volume[volume_info['type']] = {'secretName': volume_info['secret_name']}
            volume_mount = {'mountPath': volume_info['path'],
                            'name': volume_info['name']}
        elif volume_type == 'emptydir':
            volume['emptyDir'] = {}
            volume_mount = {'mountPath': volume_info['path'],
                            'name': volume_info['name']}
        elif volume_type == 'pvc' or volume_type == 'persistentvolumeclaim':
            volume['persistentVolumeClaim'] = {}
            volume['persistentVolumeClaim']['claimName'] = volume_info['claimName']
            volume['persistentVolumeClaim']['claimSize'] = volume_info['claimSize']
        elif volume_type == 'hostpath':
            volume['hostPath'] = {}
            volume['hostPath']['path'] = volume_info['path']
        elif volume_type == 'configmap':
            volume['configMap'] = {}
            volume['configMap']['name'] = volume_info['configmap_name']
            volume_mount = {'mountPath': volume_info['path'],
                            'name': volume_info['name']}

        return (volume, volume_mount)
