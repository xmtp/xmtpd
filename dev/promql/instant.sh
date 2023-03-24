#!/bin/bash
#
# Run an instant query $1, provide the list of group by labels in $BY (default 'host_name').
# The time at which to execute the query can be specified as an offset (date -v) in $AT (default now).
# If $1 is absent run the default sum query for $METRIC, grouped by $BY,
# optionally filtered by list of label conditions in $WHERE.
# The results are displayed in a table with column for each group by label
# and result value in the last column.
#
# Examples:
#   AT=-5H METRIC=api BY='host_name' WHERE='le="+Inf",grpc_method="Publish"' ./instant.sh
#   
#   # return count of series by xmptd metric
#   BY='__name__' ./instant.sh 'count by (__name__) ({__name__=~"xmtpd.*"})'

source "$(dirname "${BASH_SOURCE[0]}")"/util.sh
if [[ $METRIC ]]; then set_metric "$METRIC"; fi
group_by=${BY:-host_name}
if [[ $WHERE ]]; then labels="{$WHERE}"; fi

query=${1:-sum($metric$labels) by ($group_by)}
offset=${AT:-'+0H'}
start=$(date -v"$offset" +%s)

echo "$(date -r "$start" +%H:%M:%S)" "$query"


curl --silent http://${PROMETHEUS:-prometheus.localhost}/api/v1/query \
    --data-urlencode "query=$query" \
    --data-urlencode "time=$start" \
| jq -r --arg group_by "$group_by"  '
    if .status == "error" then . else
    ($group_by | split(",")) as $labels
    |   $labels + [ "value" ],
        ( .data.result[]
            | . as $metric
            |   ($labels
                    | reduce .[] as $l ([];
                        . + [ $metric.metric[$l] ]
                    )
                ) + [ $metric.value[1] ]
        )
    | @tsv
    end' \
| column -t
