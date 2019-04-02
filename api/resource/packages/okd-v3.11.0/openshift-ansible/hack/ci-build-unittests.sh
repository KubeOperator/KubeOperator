#!/bin/bash

# This script installs tox dependencies in the test container
yum install -y gcc libffi-devel python-devel openssl-devel python-pip
pip install tox
chmod uga+w /etc/passwd
