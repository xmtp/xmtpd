#!/bin/bash
#
# List metadata for all xmptd metrics
#
curl --silent http://${PROMETHEUS:-prometheus.localhost}/api/v1/metadata \
| jq -r '.data | with_entries(select(.key | test("xmtpd.*")))'
