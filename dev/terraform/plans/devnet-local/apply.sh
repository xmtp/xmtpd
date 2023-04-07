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
tf apply -auto-approve -target=module.cluster.module.k8s
kind_cluster="$(tf output -json | jq -r '.k8s_cluster_name.value')"
echo

# Load local docker images into kind cluster.
export KUBECONFIG="${TF_VAR_kubeconfig_path}"
nodes="$(kubectl get nodes -l "node-pool=xmtp-nodes" --no-headers -o custom-columns=":metadata.name" | paste -s -d, -)"
kind load docker-image "${CONTAINER_IMAGE_FULL}" --name "${kind_cluster}" --nodes "${nodes}"

# Apply the rest.
tf apply -auto-approve
