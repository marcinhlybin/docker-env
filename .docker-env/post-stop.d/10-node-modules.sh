#!/usr/bin/env bash

# Exit on error
set -e

# Remove node_modules symlink
[ -L node_modules ] && rm -f node_modules

exit 0
