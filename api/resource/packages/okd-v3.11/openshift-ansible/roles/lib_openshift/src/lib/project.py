# pylint: skip-file
# flake8: noqa


# pylint: disable=too-many-instance-attributes
class ProjectConfig(OpenShiftCLIConfig):
    ''' project config object '''
    def __init__(self, rname, namespace, kubeconfig, project_options):
        super(ProjectConfig, self).__init__(rname, None, kubeconfig, project_options)


class Project(Yedit):
    ''' Class to wrap the oc command line tools '''
    annotations_path = "metadata.annotations"
    kind = 'Project'
    annotation_prefix = 'openshift.io/'

    def __init__(self, content):
        '''Project constructor'''
        super(Project, self).__init__(content=content)

    def get_annotations(self):
        ''' return the annotations'''
        return self.get(Project.annotations_path) or {}

    def add_annotations(self, inc_annos):
        ''' add an annotation to the other annotations'''
        if not isinstance(inc_annos, list):
            inc_annos = [inc_annos]

        annos = self.get_annotations()
        if not annos:
            self.put(Project.annotations_path, inc_annos)
        else:
            for anno in inc_annos:
                for key, value in anno.items():
                    annos[key] = value

        return True

    def find_annotation(self, key):
        ''' find an annotation'''
        annotations = self.get_annotations()
        for anno in annotations:
            if Project.annotation_prefix + key == anno:
                return annotations[anno]

        return None

    def delete_annotation(self, inc_anno_keys):
        ''' remove an annotation from a project'''
        if not isinstance(inc_anno_keys, list):
            inc_anno_keys = [inc_anno_keys]

        annos = self.get(Project.annotations_path) or {}

        if not annos:
            return True

        removed = False
        for inc_anno in inc_anno_keys:
            anno = self.find_annotation(inc_anno)
            if anno:
                del annos[Project.annotation_prefix + anno]
                removed = True

        return removed

    def update_annotation(self, key, value):
        ''' remove an annotation for a project'''
        annos = self.get(Project.annotations_path) or {}

        if not annos:
            return True

        updated = False
        anno = self.find_annotation(key)
        if anno:
            annos[Project.annotation_prefix + key] = value
            updated = True

        else:
            self.add_annotations({Project.annotation_prefix + key: value})

        return updated
