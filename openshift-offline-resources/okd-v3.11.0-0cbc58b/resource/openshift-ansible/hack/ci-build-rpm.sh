#!/bin/bash

# This script builds RPMs for Prow CI
tito tag --offline --accept-auto-changelog --use-release '9999%{?dist}'
tito build --output="_output/local/releases" --rpm --test --offline --quiet

mkdir _output/local/releases/rpms
mv _output/local/releases/noarch/* _output/local/releases/rpms
createrepo _output/local/releases/rpms
