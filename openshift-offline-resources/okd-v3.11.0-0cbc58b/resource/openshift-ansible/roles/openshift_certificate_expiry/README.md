# OpenShift Certificate Expiration Checker

OpenShift certificate expiration checking. Be warned of certificates
expiring within a configurable window of days, and notified of
certificates which have already expired. Certificates examined
include:

* Master/Node Service Certificates
* Router/Registry Service Certificates from etcd secrets
* Master/Node/Router/Registry/Admin `kubeconfig`s
* Etcd certificates (including embedded)

This role pairs well with the redeploy certificates playbook:

* [Redeploying Certificates Documentation](https://docs.openshift.com/container-platform/latest/install_config/redeploying_certificates.html)

Just like the redeploying certificates playbook, this role is intended
to be used with an inventory that is representative of the
cluster. For best results run `ansible-playbook` with the `-v` option.


# Role Variables

Core variables in this role:

| Name                                                  | Default value                  | Description                                                           |
|-------------------------------------------------------|--------------------------------|-----------------------------------------------------------------------|
| `openshift_certificate_expiry_config_base`            | `/etc/origin`                  | Base openshift config directory                                       |
| `openshift_certificate_expiry_warning_days`           | `30`                           | Flag certificates which will expire in this many days from now        |
| `openshift_certificate_expiry_show_all`               | `no`                           | Include healthy (non-expired and non-warning) certificates in results |

Optional report/result saving variables in this role:

| Name                                                  | Default value                  | Description                                                           |
|-------------------------------------------------------|--------------------------------|-----------------------------------------------------------------------|
| `openshift_certificate_expiry_generate_html_report`   | `no`                           | Generate an HTML report of the expiry check results                   |
| `openshift_certificate_expiry_html_report_path`       | `/tmp/cert-expiry-report.html` | The full path to save the HTML report as                              |
| `openshift_certificate_expiry_save_json_results`      | `no`                           | Save expiry check results as a json file                              |
| `openshift_certificate_expiry_json_results_path`      | `/tmp/cert-expiry-report.json` | The full path to save the json report as                              |


# Using this Role

How to use the Certificate Expiration Checking Role.

> **NOTE:** In the examples shown below, ensure you change **HOSTS**
> to the path of your inventory file.

## Run with ansible-playbook

Run one of the example playbooks using an inventory file
representative of your existing cluster. Some example playbooks are
included in this role, or you can [read on below for more examples](#more-example-playbooks)
to help you craft you own.

```
$ ansible-playbook -v -i HOSTS playbooks/openshift-checks/certificate_expiry/easy-mode.yaml
```

Using the `easy-mode.yaml` playbook will produce:

* Reports including healthy and unhealthy hosts
* A JSON report in `/tmp/`
* A stylized HTML report in `/tmp/`


> **Note:** If you are running from an RPM install use
> `/usr/share/ansible/openshift-ansible/playbooks/openshift-checks/certificate_expiry/easy-mode.yaml`
> instead

## Run from a container

The example playbooks that use this role are packaged in the
[container image for openshift-ansible](../../README_CONTAINER_IMAGE.md), so you
can run any of them by setting the `PLAYBOOK_FILE` environment variable when
running an openshift-ansible container.

There are several [examples](../../examples/README.md) in the `examples` directory that run certificate check playbooks from a container running on OpenShift.

## More Example Playbooks

> **Note:** These Playbooks are available to run directly out of the
> [/playbooks/openshift-checks/certificate_expiry/](../../playbooks/openshift-checks/certificate_expiry/) directory.

### Default behavior

This playbook just invokes the certificate expiration check role with default options:


```yaml
---
- name: Check cert expirys
  hosts: nodes:masters:etcd
  become: yes
  gather_facts: no
  roles:
    - role: openshift_certificate_expiry
```

**From git:**
```
$ ansible-playbook -v -i HOSTS playbooks/openshift-checks/certificate_expiry/default.yaml
```
**From openshift-ansible-playbooks rpm:**
```
$ ansible-playbook -v -i HOSTS /usr/share/ansible/openshift-ansible/playbooks/openshift-checks/certificate_expiry/default.yaml
```

> [View This Playbook](../../playbooks/openshift-checks/certificate_expiry/default.yaml)

### Easy mode

This example playbook is great if you're just wanting to **try the
role out**. This playbook enables HTML and JSON reports. All
certificates (healthy or not) are included in the results:

```yaml
---
- name: Check cert expirys
  hosts: nodes:masters:etcd
  become: yes
  gather_facts: no
  vars:
    openshift_certificate_expiry_save_json_results: yes
    openshift_certificate_expiry_generate_html_report: yes
    openshift_certificate_expiry_show_all: yes
  roles:
    - role: openshift_certificate_expiry
```

**From git:**
```
$ ansible-playbook -v -i HOSTS playbooks/openshift-checks/certificate_expiry/easy-mode.yaml
```
**From openshift-ansible-playbooks rpm:**
```
$ ansible-playbook -v -i HOSTS /usr/share/ansible/openshift-ansible/playbooks/openshift-checks/certificate_expiry/easy-mode.yaml
```

> [View This Playbook](../../playbooks/openshift-checks/certificate_expiry/easy-mode.yaml)

### Easy mode and upload reports to masters

This example builds on top of [easy-mode.yaml](#easy-mode) and additionally
uploads a copy of the generated reports to the masters, with a timestamp in the
file names.

This is specially useful when the playbook runs from within a container, because
the reports are generated inside the container and we need a way to access them.
Uploading a copy of the reports to the masters is one way to make it easy to
access them. Alternatively you can use the
[role variables](#role-variables) that control the path of the generated reports
to point to a container volume (see the [playbook with custom paths](#generate-html-and-json-reports-in-a-custom-path) for an example).

With the container use case in mind, this playbook allows control over some
options via environment variables:

 - `CERT_EXPIRY_WARN_DAYS`: sets `openshift_certificate_expiry_warning_days`, overriding the role's default.
 - `COPY_TO_PATH`: path in the masters where generated reports are uploaded.

```yaml
---
- name: Generate certificate expiration reports
  hosts: nodes:masters:etcd
  gather_facts: no
  vars:
    openshift_certificate_expiry_save_json_results: yes
    openshift_certificate_expiry_generate_html_report: yes
    openshift_certificate_expiry_show_all: yes
    openshift_certificate_expiry_warning_days: "{{ lookup('env', 'CERT_EXPIRY_WARN_DAYS') | default('45', true) }}"
  roles:
    - role: openshift_certificate_expiry

- name: Upload reports to master
  hosts: masters
  gather_facts: no
  vars:
    destination_path: "{{ lookup('env', 'COPY_TO_PATH') | default('/etc/origin/certificate_expiration_report', true) }}"
    timestamp: "{{ lookup('pipe', 'date +%Y%m%d') }}"
  tasks:
    - name: Create directory in masters
      file:
        path: "{{ destination_path }}"
        state: directory
    - name: Copy the reports to the masters
      copy:
        dest: "{{ destination_path }}/{{ timestamp }}-{{ item }}"
        src: "/tmp/{{ item }}"
      with_items:
        - "cert-expiry-report.html"
        - "cert-expiry-report.json"
```

**From git:**
```
$ ansible-playbook -v -i HOSTS playbooks/openshift-checks/certificate_expiry/easy-mode-upload.yaml
```
**From openshift-ansible-playbooks rpm:**
```
$ ansible-playbook -v -i HOSTS /usr/share/ansible/openshift-ansible/playbooks/openshift-checks/certificate_expiry/easy-mode-upload.yaml
```

> [View This Playbook](../../playbooks/openshift-checks/certificate_expiry/easy-mode-upload.yaml)

### Generate HTML and JSON artifacts in their default paths

```yaml
---
- name: Check cert expirys
  hosts: nodes:masters:etcd
  become: yes
  gather_facts: no
  vars:
    openshift_certificate_expiry_generate_html_report: yes
    openshift_certificate_expiry_save_json_results: yes
  roles:
    - role: openshift_certificate_expiry
```

**From git:**
```
$ ansible-playbook -v -i HOSTS playbooks/openshift-checks/certificate_expiry/html_and_json_default_paths.yaml
```
**From openshift-ansible-playbooks rpm:**
```
$ ansible-playbook -v -i HOSTS /usr/share/ansible/openshift-ansible/playbooks/openshift-checks/certificate_expiry/html_and_json_default_paths.yaml
```

> [View This Playbook](../../playbooks/openshift-checks/certificate_expiry/html_and_json_default_paths.yaml)

### Generate HTML and JSON reports in a custom path

This example customizes the report generation path to point to a specific path (`/var/lib/certcheck`) and uses a date timestamp for the generated files. This allows you to reuse a certain location to keep multiple copies of the reports.

```yaml
---
- name: Check cert expirys
  hosts: nodes:masters:etcd
  become: yes
  gather_facts: no
  vars:
    openshift_certificate_expiry_generate_html_report: yes
    openshift_certificate_expiry_save_json_results: yes
    timestamp: "{{ lookup('pipe', 'date +%Y%m%d') }}"
    openshift_certificate_expiry_html_report_path: "/var/lib/certcheck/{{ timestamp }}-cert-expiry-report.html"
    openshift_certificate_expiry_json_results_path: "/var/lib/certcheck/{{ timestamp }}-cert-expiry-report.json"
  roles:
    - role: openshift_certificate_expiry
```

**From git:**
```
$ ansible-playbook -v -i HOSTS playbooks/openshift-checks/certificate_expiry/html_and_json_timestamp.yaml
```
**From openshift-ansible-playbooks rpm:**
```
$ ansible-playbook -v -i HOSTS /usr/share/ansible/openshift-ansible/playbooks/openshift-checks/certificate_expiry/html_and_json_timestamp.yaml
```

> [View This Playbook](../../playbooks/openshift-checks/certificate_expiry/html_and_json_timestamp.yaml)

### Long warning window

Change the expiration warning window to 1500 days (good for testing
the module out):

```yaml
---
- name: Check cert expirys
  hosts: nodes:masters:etcd
  become: yes
  gather_facts: no
  vars:
    openshift_certificate_expiry_warning_days: 1500
  roles:
    - role: openshift_certificate_expiry
```

**From git:**
```
$ ansible-playbook -v -i HOSTS playbooks/openshift-checks/certificate_expiry/longer_warning_period.yaml
```
**From openshift-ansible-playbooks rpm:**
```
$ ansible-playbook -v -i HOSTS /usr/share/ansible/openshift-ansible/playbooks/openshift-checks/certificate_expiry/longer_warning_period.yaml
```

> [View This Playbook](../../playbooks/openshift-checks/certificate_expiry/longer_warning_period.yaml)

### Long warning window and JSON report

Change the expiration warning window to 1500 days (good for testing
the module out) and save the results as a JSON file:

```yaml
---
- name: Check cert expirys
  hosts: nodes:masters:etcd
  become: yes
  gather_facts: no
  vars:
    openshift_certificate_expiry_warning_days: 1500
    openshift_certificate_expiry_save_json_results: yes
  roles:
    - role: openshift_certificate_expiry
```

**From git:**
```
$ ansible-playbook -v -i HOSTS playbooks/openshift-checks/certificate_expiry/longer-warning-period-json-results.yaml
```
**From openshift-ansible-playbooks rpm:**
```
$ ansible-playbook -v -i HOSTS /usr/share/ansible/openshift-ansible/playbooks/openshift-checks/certificate_expiry/longer-warning-period-json-results.yaml
```

> [View This Playbook](../../playbooks/openshift-checks/certificate_expiry/longer-warning-period-json-results.yaml)



# Output Formats

As noted above there are two ways to format your check report. In
`json` format for machine parsing, or as a stylized `html` page for
easy skimming. These options are shown below.

## HTML Report

![HTML Expiration Report](examples/cert-expiry-report-html.png)

For an example of the HTML report you can browse, save
[examples/cert-expiry-report.html](examples/cert-expiry-report.html)
and then open the file in your browser.


## JSON Report

There are two top-level keys in the saved JSON results, `data` and
`summary`.

The `data` key is a hash where the keys are the names of each host
examined and the values are the check results for the certificates
identified on each respective host.

The `summary` key is a hash that summarizes the total number of
certificates:

* examined on the entire cluster
* OK
* expiring within the configured warning window
* already expired

For an example of the full JSON report, see [examples/cert-expiry-report.json](examples/cert-expiry-report.json).

The example below is abbreviated to save space.

```json
{
  "data": {
    "m01.example.com": {
      "etcd": [
        {
          "cert_cn": "CN:172.30.0.1, DNS:kubernetes, DNS:kubernetes.default, DNS:kubernetes.default.svc,...",
          "days_remaining": 722,
          "expiry": "2019-01-09 17:00:03",
          "health": "warning",
          "path": "/etc/origin/master/etcd.server.crt",
          "serial": 7,
          "serial_hex": "0x7"
        }
      ],
      "kubeconfigs": [
        {
          "cert_cn": "O:system:nodes, CN:system:node:m01.example.com",
          "days_remaining": 722,
          "expiry": "2019-01-09 17:03:28",
          "health": "warning",
          "path": "/etc/origin/node/system:node:m01.example.com.kubeconfig",
          "serial": 11,
          "serial_hex": "0xb"
        }
      ],
      "meta": {
        "checked_at_time": "2017-01-17 10:36:25.230920",
        "show_all": "True",
        "warn_before_date": "2021-02-25 10:36:25.230920",
        "warning_days": 1500
      },
      "ocp_certs": [
        {
          "cert_cn": "CN:172.30.0.1, DNS:kubernetes, DNS:kubernetes.default, DNS:kubernetes.default.svc,...",
          "days_remaining": 722,
          "expiry": "2019-01-09 17:00:02",
          "health": "warning",
          "path": "/etc/origin/master/master.server.crt",
          "serial": 4,
          "serial_hex": "0x4"
        }
      ],
      "registry": [
        {
          "cert_cn": "CN:172.30.242.251, DNS:docker-registry-default.router.default.svc.cluster.local,...",
          "days_remaining": 722,
          "expiry": "2019-01-09 17:05:54",
          "health": "warning",
          "path": "/api/v1/namespaces/default/secrets/registry-certificates",
          "serial": 13,
          "serial_hex": "0xd"
        }
      ],
      "router": [
        {
          "cert_cn": "CN:router.default.svc, DNS:router.default.svc, DNS:router.default.svc.cluster.local",
          "days_remaining": 722,
          "expiry": "2019-01-09 17:05:46",
          "health": "warning",
          "path": "/api/v1/namespaces/default/secrets/router-certs",
          "serial": 5050662940948454653,
          "serial_hex": "0x46178f2f6b765cfd"
        }
      ]
    },
    "n01.example.com": {
      "etcd": [],
      "kubeconfigs": [
        {
          "cert_cn": "O:system:nodes, CN:system:node:n01.example.com",
          "days_remaining": 722,
          "expiry": "2019-01-09 17:03:28",
          "health": "warning",
          "path": "/etc/origin/node/system:node:n01.example.com.kubeconfig",
          "serial": 11,
          "serial_hex": "0xb"
        }
      ],
      "meta": {
        "checked_at_time": "2017-01-17 10:36:25.217103",
        "show_all": "True",
        "warn_before_date": "2021-02-25 10:36:25.217103",
        "warning_days": 1500
      },
      "ocp_certs": [
        {
          "cert_cn": "CN:192.168.124.11, DNS:n01.example.com, DNS:192.168.124.11, IP Address:192.168.124.11",
          "days_remaining": 722,
          "expiry": "2019-01-09 17:03:29",
          "health": "warning",
          "path": "/etc/origin/node/server.crt",
          "serial": 12,
          "serial_hex": "0xc"
        }
      ],
      "registry": [],
      "router": []
    }
  },
  "summary": {
    "expired": 0,
    "ok": 3,
    "total": 15,
    "warning": 12
  }
}
```

The `summary` from the json data can be easily checked for
warnings/expirations using a variety of command-line tools.

For exampe, using `grep` we can look for the word `summary` and print
out the 2 lines **after** the match (`-A2`):

```
$ grep -A2 summary /tmp/cert-expiry-report.json
    "summary": {
        "warning": 16,
        "expired": 0
```

If available, the [jq](https://stedolan.github.io/jq/) tool can also
be used to pick out specific values. Example 1 and 2 below show how to
select just one value, either `warning` or `expired`. Example 3 shows
how to select both values at once:

```
$ jq '.summary.warning' /tmp/cert-expiry-report.json
16
$ jq '.summary.expired' /tmp/cert-expiry-report.json
0
$ jq '.summary.warning,.summary.expired' /tmp/cert-expiry-report.json
16
0
```


# Requirements
* None


# Dependencies
* None


# License
Apache License, Version 2.0


# Author Information
Tim Bielawa (tbielawa@redhat.com)
