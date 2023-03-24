#!/bin/bash
#
# List series for metric matching $1
# Reduce the series to those matching label conditions in $2
#
# Examples:
#   ./series.sh api 'host_name="node1-0",le="+Inf",grpc_method="Publish"'
#
source "$(dirname "${BASH_SOURCE[0]}")"/util.sh
set_metric "$1"
if [[ $2 ]]; then
condition="{$2}"
fi

# shellcheck disable=SC2154
curl --silent http://${PROMETHEUS:-prometheus.localhost}/api/v1/series \
    --data-urlencode "match[]=$metric$condition" \
| jq -r '.data'
