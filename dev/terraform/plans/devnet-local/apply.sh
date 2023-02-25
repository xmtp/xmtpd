#!/bin/bash
set -eo pipefail
plan_dir="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

function tf() {
    terraform -chdir="${plan_dir}" "$@"
}

# Build docker images.
(
    cd "${plan_dir}/../../../../"
    export BUILDKIT_PROGRESS=plain
    pids=()
    dev/docker/xmtpd/build &
    pids+=("$!")
    for pid in ${pids[@]+"${pids[@]}"}; do
        wait "${pid}"
    done
)
CONTAINER_IMAGE_ID="$(docker inspect --format='{{.Id}}' "${CONTAINER_IMAGE}")"
CONTAINER_IMAGE_FULL="${CONTAINER_IMAGE}-${CONTAINER_IMAGE_ID##*:}"
docker tag "${CONTAINER_IMAGE}" "${CONTAINER_IMAGE_FULL}"

export TF_VAR_node_container_image="${CONTAINER_IMAGE_FULL}"
export TF_VAR_kubeconfig_path="${plan_dir}/.xmtp/kubeconfig.yaml"

# Initialize terraform.
tf init -upgrade

# Create clusters.
tf apply -auto-approve -target=module.cluster
kind_cluster="$(tf output -json | jq -r '.cluster_name.value')"
echo

# Load local docker images into kind cluster.
export KUBECONFIG="${TF_VAR_kubeconfig_path}"
nodes="$(kubectl get nodes -l "node-pool=xmtp-nodes" --no-headers -o custom-columns=":metadata.name" | paste -s -d, -)"
kind load docker-image "${CONTAINER_IMAGE_FULL}" --name "${kind_cluster}" --nodes "${nodes}"

# Apply the rest.
tf apply -auto-approve

# Sync local xmtp-node helm chart as node repo source in argocd.
argocd_namespace="$(tf output -json | jq -r '.argocd_namespace.value')"
argocd_hostname="$(tf output -json | jq -r '.argocd_hostnames.value[0]')"
argocd_username="$(tf output -json | jq -r '.argocd_username.value')"
argocd_password="$(tf output -json | jq -r '.argocd_password.value')"
chronic argocd login "${argocd_hostname}" --port-forward --port-forward-namespace "${argocd_namespace}" --plaintext --username "${argocd_username}" --password "${argocd_password}"

nodes="$(tf output -json | jq -r '.nodes | .value')"
while IFS= read -r node; do
    name="$(echo "${node}" | jq -r '.name')"
    echo "Syncing ${name} argo app from local"
    chronic argocd --port-forward --port-forward-namespace="${argocd_namespace}" app sync "${name}" --local "${plan_dir}/../../../helm/xmtp-node"
done < <(printf '%s\n' "$(echo "${nodes}" | jq -rc '.[]')")
