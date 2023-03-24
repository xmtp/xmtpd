#!/bin/bash
#
# Run a range query $1, provide the list of group by labels in $BY (default 'host_name').
# The range over which to execute the query is specified as an offsets (date -v)
# in $FROM (default -1H) and $TO (default now).
# $STEP defines the range query step (default 30s).
# If $1 is absent run the default sum/increase query for $METRIC, grouped by $BY,
# optionally filtered by list of label conditions in $WHERE.
# The resulting timeseries are then displayed in a table
# showing a column for each series with timestamps in the first column.
#
# Examples:
#   FROM=-5H TO=-4H METRIC=api BY=host_name WHERE='le="+Inf",grpc_method="Publish"' ./range.sh
#   METRIC=fetch BY=le WHERE='host_name="node1-0",grpc_method="Publish"' ./range.sh
#

source "$(dirname "${BASH_SOURCE[0]}")"/util.sh
if [[ $METRIC ]]; then set_metric "$METRIC"; fi
group_by=${BY:-host_name}
if [[ $WHERE ]]; then labels="{$WHERE}"; fi

query=${1:-sum(increase($metric${labels}[30s])) by ($group_by)}

fromOffset=${FROM:-'-1H'}   # date -v offset
toOffset=${TO:-'+0H'}
step=${STEP:-'30s'}           # golang duration value
startt=$(date -v "$fromOffset" +%s)
endt=$(date -v "$toOffset" +%s)

echo "$(date -r "$startt" +%H:%M:%S)" "$(date -r "$endt" +%H:%M:%S)" "$query"

curl --silent http://${PROMETHEUS:-prometheus.localhost}/api/v1/query_range \
    --data-urlencode "query=$query" \
    --data-urlencode "start=$startt" \
    --data-urlencode "end=$endt" \
    --data-urlencode "step=$step" \
| jq -r --arg group_by "$group_by"  '
    if .status == "error" then . else
    (   ($group_by | split(",")) as $labels 
        | .data.result as $result
        | $labels[]
        | . as $l
        |   reduce $result[].metric as $metric ([$l];
                . +  [ $metric[$l] ]
            )
    ),
    (   reduce .data.result[].values[] as $sample ({};
            .[ ($sample[0] | strftime("%H:%M:%S")) ] |= . + [ $sample[1] ])
        | to_entries[] | [.key] + .value
    )
    | @tsv
    end' \
| column -t
