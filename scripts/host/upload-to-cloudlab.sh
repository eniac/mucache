#!/bin/bash

ip=${1?IP of remote machine not given}
user=${2?user in remote machine}

## Optionally the caller can give us a private key for the ssh
key=$3
if [ -z "$key" ]; then
    key_flag=""
else
    key_flag="-i ${key}"
fi

export mucache_dir=${MUCACHE_TOP:-$(git rev-parse --show-toplevel --show-superproject-working-tree)}

## Upload the whole knproto directory to a remote machine
rsync --rsh="ssh -p 22 ${key_flag}" --progress -p -r "${mucache_dir}" "${user}@${ip}:/users/${user}"
