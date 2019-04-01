import os
import sys

# extend sys.path so that tests can import openshift_checks and action plugins
# from this role.
openshift_health_checker_path = os.path.dirname(os.path.dirname(__file__))
sys.path[1:1] = [
    openshift_health_checker_path,
    os.path.join(openshift_health_checker_path, 'action_plugins'),
    os.path.join(openshift_health_checker_path, 'callback_plugins'),
    os.path.join(openshift_health_checker_path, 'library'),
]
