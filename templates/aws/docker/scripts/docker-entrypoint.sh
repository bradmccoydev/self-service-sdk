#!/bin/bash
set -e

# Get the arguments
CMDLINE="$@"

# Execute
exec ${CMDLINE:-/bin/bash}