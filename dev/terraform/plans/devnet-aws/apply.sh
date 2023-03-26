#!/bin/bash
set -eo pipefail
plan_dir="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
set -a; source "${plan_dir}/.env"; set +a

function tf() {
    terraform -chdir="${plan_dir}" "$@"
}

# Initialize terraform.
terraform init -upgrade \
     -backend-config="bucket=${TFSTATE_AWS_BUCKET}" \
     -backend-config="key=${PLAN}/terraform.tfstate" \
     -backend-config="encrypt=true" \
     -backend-config="dynamodb_table=${TFSTATE_AWS_BUCKET}-locking" \
     -backend-config="region=${TFSTATE_AWS_REGION}"

# Create cluster and container registry.
tf apply -auto-approve \
    -target=module.cluster \
    -target=module.ecr_node_repo
tf apply -auto-approve \
    -target=module.system.kubernetes_namespace.system \
    -target=module.system.helm_release.argocd
node_repo_url="$(tf output -json | jq -r '.ecr_node_repo_url.value')"
aws_region="$(tf output -json | jq -r '.region.value')"
ecr_repo_id="$(tf output -json | jq -r '.ecr_node_repo_id.value')"
echo

# Login to registry with docker.
aws ecr get-login-password --region "${aws_region}" | docker login --username AWS --password-stdin "${ecr_repo_id}.dkr.ecr.${aws_region}.amazonaws.com"

# Build docker images.
CONTAINER_IMAGE="${node_repo_url}:dev"
(
    cd "${plan_dir}/../../../../"
    export BUILD_CONTAINER_IMAGE="${CONTAINER_IMAGE}"
    export BUILDKIT_PROGRESS=plain
    pids=()
    dev/docker/xmtpd/buildx &
    pids+=("$!")
    for pid in ${pids[@]+"${pids[@]}"}; do
        wait "${pid}"
    done
)
CONTAINER_IMAGE_ID="$(docker inspect --format='{{.Id}}' "${CONTAINER_IMAGE}")"
CONTAINER_IMAGE_FULL="${CONTAINER_IMAGE}-${CONTAINER_IMAGE_ID##*:}"
docker tag "${CONTAINER_IMAGE}" "${CONTAINER_IMAGE_FULL}"
docker push "${CONTAINER_IMAGE_FULL}"
export TF_VAR_node_container_image="${CONTAINER_IMAGE_FULL}"

# Apply the rest.
tf apply -auto-approve
