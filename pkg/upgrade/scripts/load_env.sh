#!/bin/bash

set -eu

# Get the directory where the script is located
SCRIPT_DIR=$(dirname "$(realpath "$0")")

# Navigate to the top-level directory
TOP_LEVEL_DIR=$(realpath "$SCRIPT_DIR/../../..")

cd $TOP_LEVEL_DIR

. ./dev/local.env

env | grep XMTPD