#!/usr/bin/env bash
#
# Exit on error
set -e

# Run pre-start scripts
for f in .docker-env/pre-start.d/*; do
  if [ -x "$f" ]; then
    echo "(pre-start) Running $f with args $@"
    "$f" "$@"
  fi
done

