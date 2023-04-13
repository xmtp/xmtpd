<!-- markdownlint-disable MD034 MD041-->

This is a set of helper scripts to query local prometheus directly through its [HTTP API](https://prometheus.io/docs/prometheus/latest/querying/api/).

To understand promql queries read

* https://promlabs.com/blog/2020/06/18/the-anatomy-of-a-promql-query/
* https://promlabs.com/blog/2020/09/25/metric-types-in-prometheus-and-promql/
* https://promlabs.com/blog/2020/07/02/selecting-data-in-promql/
* https://promlabs.com/promql-cheat-sheet/

In general there is a separate script for each endpoint:

* instant.sh - instant query (single set of results)
* range.sh   - range query (results over a period of time)
* series.sh  - list individual timeseries of a metric (just properties, not values)
* labels.sh  - list labels of a metric or all metrics
* label.sh   - list values of a label
* meta.sh    - list all xmtpd metrics
* util.sh    - helper to translate a part of metric name into its full name,
               e.g. METRIC=api -> xmtpd_api_request_duration_ms_bucket
