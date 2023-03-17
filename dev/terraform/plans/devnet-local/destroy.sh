#!/bin/bash
set -eou pipefail
plan_dir="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

function tf() {
    terraform -chdir="${plan_dir}" "$@"
}

kind_cluster="$(tf output -json | jq -r '.cluster_name.value')"
kind delete cluster --name "${kind_cluster}"

../../clean
