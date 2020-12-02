#!/usr/bin/env bash

# shellcheck disable=SC1090
source "${E2E_DIR}/lib/defer.bash"

function package_and_upload_chart() {
  local chart=${1}
  local chart_repository=${2}

  gen_dir=$(mktemp -d)
  defer rm -rf "'$gen_dir'"

  # Package
  if [ "$HELM_VERSION" == "v3" ]; then
    helm3 package --destination "$gen_dir" "$chart"
  else
    helm2 package --destination "$gen_dir" "$chart"
  fi

  # Upload
  chart_tarbal=$(find "$gen_dir" -type f -name "*.tgz" | head -n1)
  curl --data-binary "@$chart_tarbal" "$chart_repository/api/charts"
}
