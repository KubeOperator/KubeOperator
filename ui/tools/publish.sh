#!/bin/bash
# coding: utf-8
# Copyright (c) 2017
# Gmail:liuzheng712
#

set -ex


git checkout github_dev && \
  git pull github dev --rebase && \
  git merge dev -m "publish" && \
  git reset --soft HEAD^ && \
  git commit -m "publish" && \
  git push github github_dev:dev && \
  echo "success"
git checkout dev
git pull github dev --commit && git push origin dev
