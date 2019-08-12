import os
import shutil

from jinja2 import FileSystemLoader, Environment


def generate_terraform_file(target_path, cloud_path, vars):
    terraform_path = os.path.join(cloud_path, "terraform")
    lorder = FileSystemLoader(terraform_path)
    env = Environment(loader=lorder)
    _template = env.get_template("terraform.tf.j2")
    result = _template.render(vars)
    if not os.path.exists(target_path):
        os.makedirs(target_path)
    file = os.path.join(target_path, 'main.tf')
    with open(file, 'w') as f:
        f.write(result)
        return os.path.dirname(file)


def init_terraform(target_path, cloud_path):
    if not os.path.exists(os.path.join(target_path, '.terraform', 'plugins')):
        shutil.copytree(os.path.join(os.path.join(cloud_path, "terraform", "plugins")),
                        os.path.join(target_path, '.terraform', 'plugins'))
