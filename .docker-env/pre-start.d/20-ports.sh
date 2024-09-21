#!/usr/bin/env bash
#
# Check ports 80 and 443
#

# Exit on error
set -e

PORTS="80 443"

function port_is_open {
  local port="$1"
  lsof -i TCP:$port -sTCP:LISTEN -P >/dev/null
}

function port_used_by {
    local port="$1"
    ps -o comm= -p $(lsof -i TCP:$port -sTCP:LISTEN -t | head -1)
}

for PORT in $PORTS; do
    if port_is_open $PORT; then
        comm=$(port_used_by $PORT)
        echo "WARNING! Port $PORT is already in use by $comm" >&2
    fi
done

exit 0
