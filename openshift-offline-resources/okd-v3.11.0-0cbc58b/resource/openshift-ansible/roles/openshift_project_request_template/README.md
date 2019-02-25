OpenShift Project Request Template
==================================

Configure template used when creating new projects. If enabled only the template is managed. It must still be enabled in the OpenShift master configuration. The base template is created using `oc adm create-bootstrap-project-template` and can be modified by setting `openshift_project_request_template_edits`.


Requirements
------------


Role Variables
--------------

From this role:

| Name                                         | Default value   | Description                                    |
|----------------------------------------------|-----------------|------------------------------------------------|
| openshift_project_request_template_manage    | false           | Whether to manage the project request template |
| openshift_project_request_template_namespace | default         | Namespace for template                         |
| openshift_project_request_template_name      | project-request | Template name                                  |
| openshift_project_request_template_edits     | []              | Changes for template                           |


Dependencies
------------

* lib_utils


License
-------

Apache License Version 2.0
