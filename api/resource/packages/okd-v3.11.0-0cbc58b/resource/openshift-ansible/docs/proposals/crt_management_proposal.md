# Container Runtime Management

## Description
origin and openshift-ansible support multiple container runtimes.  This proposal
is related to refactoring how we handle those runtimes in openshift-ansible.

### Problems addressed
We currently don't install docker during the install at a point early enough to
not fail health checks, and we don't have a good story around when/how to do it.
This is complicated by logic around containerized and non-containerized installs.

A web of dependencies can cause changes to docker that are unintended and has
resulted in a series of work-around such as 'skip_docker' boolean.

We don't handle docker storage because it's BYO.  By moving docker to a prerequisite
play, we can tackle storage up front and never have to touch it again.

container_runtime logic is currently spread across 3 roles: docker, openshift_docker,
and openshift_docker_facts.  The name 'docker' does not accurately portray what
the role(s) do.

## Rationale
* Refactor docker (and related meta/fact roles) into 'container_runtime' role.
* Strip all meta-depends on container runtime out of other roles and plays.
* Create a 'prerequisites.yml' entry point that will setup various items
such as container storage and container runtime before executing installation.
* All other roles and plays should merely consume container runtime, should not
configure, restart, or change the container runtime as much as feasible.

## Design

The container_runtime role should be comprised of 3 'pseudo-roles' which will be
consumed using import_role; each component area should be enabled/disabled with
a boolean value, defaulting to true.

I call them 'pseudo-roles' because they are more or less independent functional
areas that may share some variables and act on closely related components.  This
is an effort to reuse as much code as possible, limit role-bloat (we already have
an abundance of roles), and make things as modular as possible.

```yaml
# prerequisites.yml
- include: std_include.yml
- include: container_runtime_setup.yml
...
# container_runtime_setup.yml
- hosts: "{{ openshift_runtime_manage_hosts | default('oo_nodes_to_config') }}"
  tasks:
    - import_role:
        name: container_runtime
        tasks_from: install.yml
      when: openshift_container_runtime_install | default(True) | bool
    - import_role:
        name: container_runtime
        tasks_from: storage.yml
      when: openshift_container_runtime_storage | default(True) | bool
    - import_role:
        name: container_runtime
        tasks_from: configure.yml
      when: openshift_container_runtime_configure | default(True) | bool
```

Note the host group on the above play.  No more guessing what hosts to run this
stuff against.  If you want to use an atomic install, specify what hosts will need
us to setup container runtime (such as etcd hosts, loadbalancers, etc);

We should direct users that are using atomic hosts to disable install in the docs,
let's not add a bunch of logic.

Alternatively, we can create a new group.

### Part 1, container runtime install
Install the container runtime components of the desired type.

```yaml
# install.yml
- include: docker.yml
  when: openshift_container_runtime_install_docker | bool

- include: crio.yml
  when: openshift_container_runtime_install_crio | bool

... other container run times...
```

Alternatively to using booleans for each run time, we could use a variable like
"openshift_container_runtime_type".  This would be my preference, as we could
use this information in later roles.

### Part 2, configure/setup container runtime storage
Configure a supported storage solution for containers.

Similar setup to the previous section.  We might need to add some logic for the
different runtimes here, or we maybe create a matrix of possible options.

### Part 3, configure container runtime.
Place config files, environment files, systemd units, etc.  Start/restart
the container runtime as needed.

Similar to Part 1 with how we should do things.

## Checklist
* Strip docker from meta dependencies.
* Combine docker facts and meta roles into container_runtime role.
* Docs

## User Story
As a user of openshift-ansible, I want to be able to manage my container runtime
and related components independent of openshift itself.

## Acceptance Criteria
* Verify that each container runtime installs with this new method.
* Verify that openshift installs with this new method.
