# find xmptd metric matching $1
# set $metrict_meta (name, type, help, unit), $metric_type, $metric
function set_metric() {
    metric_meta=$(
        curl --silent http://"${PROMETHEUS:-prometheus.localhost}"/api/v1/metadata \
        | jq -rM ".data | to_entries[]
            | select(.key | test(\"xmtpd.*\"))
            | select(.key | test(\".*$1.*\"))
            | .value[0] + {name: .key}"
    )
    metric_type=$(jq -nr --argjson meta "$metric_meta" '$meta.type')
    metric=$(jq -nr --argjson meta "$metric_meta" '$meta.name')
    # add _bucket suffix for histogram metrics
    if [[ $metric_type == "histogram" ]]; then
        metric=${metric}_bucket;
    fi 
}
