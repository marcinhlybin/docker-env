#!/usr/bin/env bash
#
# Exit on error
set -e

# Run post-start scripts
for f in .docker-env/post-start.d/*; do
  if [ -x "$f" ]; then
    echo "(post-start) Running $f with args $@"
    "$f" "$@"
  fi
done

