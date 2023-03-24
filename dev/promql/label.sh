#!/bin/bash
#
# List values for specified label
# Optionally select metric/series to read the values from.
#
# Examples:
#   ./label.sh grpc_method api 'host_name="node1-0",le="+Inf"'
#   ./label.sh __name__   # all values of __name__ label => names of all metrics
#
args=()
source "$(dirname "${BASH_SOURCE[0]}")"/util.sh
if [[ $3 ]]; then condition="{$3}"; fi
if [[ $2 ]]; then
    set_metric "$2"
    args+=("--data-urlencode match[]=$metric$condition")
fi

curl --silent http://${PROMETHEUS:-prometheus.localhost}/api/v1/label/"$1"/values -X GET "${args[@]}" \
| jq -r '.data[]'
