#!/bin/bash
#
# List labels for specified metric/series or all metrics.
#
# Examples:
#   ./labels.sh api 'host_name="node1-0",le="+Inf",grpc_method="Publish"'
#
args=()
source "$(dirname "${BASH_SOURCE[0]}")"/util.sh
if [[ $2 ]]; then condition="{$2}"; fi
if [[ $1 ]]; then
    set_metric "$1"
    args+=("--data-urlencode match[]=$metric$condition")
fi

curl --silent http://"${PROMETHEUS:-prometheus.localhost}"/api/v1/labels "${args[@]}" \
| jq -r '.data[]'
