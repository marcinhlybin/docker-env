#!/usr/bin/env bash
#
# Check if ssh-agent is running

# Exit on error
set -e

function ssh_found_in_ps {
    ps aux | grep [s]sh-agent > /dev/null 2>&1
}

function ssh_auth_sock_defined {
    [ -n "$SSH_AUTH_SOCK" ]
}

if ! ssh_auth_sock_defined || ! ssh_found_in_ps; then
    echo "Ssh-agent is not running"
    echo "Start ssh-agent with 'eval \$(ssh-agent -s)' and add your key with 'ssh-add'"
    exit 1
fi

exit 0
