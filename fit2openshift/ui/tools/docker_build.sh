#!/bin/bash

set -ex
npm i
npm run-script build
docker build -t fitopenshift-portal .
