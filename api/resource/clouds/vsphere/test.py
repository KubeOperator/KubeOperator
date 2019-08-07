import json
import os

from jinja2 import FileSystemLoader, Environment


def test():
    json_obj = json.load(open("test.json", 'r'))
    file_name = "/Users/shenchenyang/PycharmProjects/new/KubeOperator/api/resource/clouds/vshpere/terraform.tf.j2"
    lorder = FileSystemLoader(os.path.dirname(file_name))
    env = Environment(loader=lorder)
    template = env.get_template(os.path.basename(file_name))
    result = template.render(json_obj)
    f = open('result.tf', 'w')
    f.write(result)
    f.close()


test()
