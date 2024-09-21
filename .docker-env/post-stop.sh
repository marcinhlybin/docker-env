#!/usr/bin/env bash
#
# Exit on error
set -e

# Run post-stop scripts
for f in .docker-env/post-stop.d/*; do
  if [ -x "$f" ]; then
    echo "(post-stop) Running $f"
    "$f"
  fi
done

