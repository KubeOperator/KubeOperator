#!/bin/bash

# This script runs tox tests in test container
echo "${USER:-default}:x:$(id -u):$(id -g):Default User:${HOME:-/tmp}:/sbin/nologin" >> /etc/passwd
tox 2>&1 | tee /tmp/artifacts/output.log
