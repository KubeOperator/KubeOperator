#!/bin/bash

branch=$1

echo "`date -u '+%Y-%m-%d %H:%M:%S'` Start to build UI"
rm -rf dist/
docker run -i --rm -v `pwd`:/data -w /data node:9  sh -c "npm i && npm run-script build"

