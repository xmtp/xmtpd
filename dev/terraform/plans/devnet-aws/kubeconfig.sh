#!/bin/bash
set -eo pipefail
plan_dir="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
set -a; source "${plan_dir}/.env"; set +a

function tf() {
    terraform -chdir="${plan_dir}" "$@"
}

region="$(tf output -json | jq -r '.region.value')"
cluster_id="$(tf output -json | jq -r '.eks_cluster_id.value')"

mkdir -p "${plan_dir}/.xmtp"
aws eks update-kubeconfig --region "${region}" --name "${cluster_id}" --kubeconfig "${plan_dir}/.xmtp/kubeconfig.yaml"
