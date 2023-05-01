#!/bin/bash
set -eou pipefail
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
    dev/docker/xmtpd-e2e/build &
    pids+=("$!")
    for pid in ${pids[@]+"${pids[@]}"}; do
        wait "${pid}"
    done
)
node_container_image_id="$(docker inspect --format='{{.Id}}' "${CONTAINER_IMAGE}")"
node_container_image_full="${CONTAINER_IMAGE}-${node_container_image_id##*:}"
docker tag "${CONTAINER_IMAGE}" "${node_container_image_full}"

e2e_container_image_id="$(docker inspect --format='{{.Id}}' "${E2E_CONTAINER_IMAGE}")"
e2e_container_image_full="${E2E_CONTAINER_IMAGE}-${e2e_container_image_id##*:}"
docker tag "${E2E_CONTAINER_IMAGE}" "${e2e_container_image_full}"

export TF_VAR_node_container_image="${node_container_image_full}"
export TF_VAR_e2e_container_image="${e2e_container_image_full}"
export TF_VAR_kubeconfig_path="${plan_dir}/.xmtp/kubeconfig.yaml"

# Initialize terraform.
tf init -upgrade

# Create clusters.
tf apply -auto-approve -target=module.cluster.module.k8s
kind_cluster="$(tf output -json | jq -r '.k8s_cluster_name.value')"
echo

# Load local docker images into kind cluster.
export KUBECONFIG="${TF_VAR_kubeconfig_path}"
nodes="$(kubectl get nodes -l "node-pool=xmtp-nodes" --no-headers -o custom-columns=":metadata.name" | paste -s -d, -)"
kind load docker-image "${node_container_image_full}" --name "${kind_cluster}" --nodes "${nodes}"
nodes="$(kubectl get nodes -l "node-pool=xmtp-tools" --no-headers -o custom-columns=":metadata.name" | paste -s -d, -)"
kind load docker-image "${e2e_container_image_full}" --name "${kind_cluster}" --nodes "${nodes}"

# Apply the rest.
tf apply -auto-approve "$@"
