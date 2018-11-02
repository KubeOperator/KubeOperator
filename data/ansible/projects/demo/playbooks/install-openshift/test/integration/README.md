# Integration tests

Integration tests exercise the OpenShift Ansible playbooks by running them
against an inventory with Docker containers as hosts.

## Requirements

The tests assume that:

* docker is running on localhost and the present user has access to use it.
* golang is installed and the go binary is in PATH.
* python and tox are installed.

## Building images

The tests rely on images built in the local docker index. You can build them
from the repository root with:

```
./test/integration/build-images.sh
```

Use the `--help` option to view available options.

## Running the tests

From the repository root, run the integration tests with:

```
./test/integration/run-tests.sh
```

Use the `--help` option to view available options.

You can also run tests more directly, for example to run a specific check:

```
go test ./test/integration/... -run TestPackageUpdateDepMissing
```
