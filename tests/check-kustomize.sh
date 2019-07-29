#!/bin/bash

TRIES=${TRIES:-5}
WAIT=${WAIT:-5}

set -uo pipefail

for i in $(seq $TRIES) ; do
    cat /credentials/vault-role-id && cat /credentials/vault-secret-id
    result=$?
    if [ $result -eq 0 ] ; then
        if [ $i -eq $TRIES ]; then
            echo "timeout exceeded!"
            exit 1
        fi
        break
    fi
    echo "vault is not ready, sleeping for $WAIT seconds..."
    sleep $WAIT
done

set -e

export VAULT_ROLE_ID=$(cat /credentials/vault-role-id)
export VAULT_SECRET_ID=$(cat /credentials/vault-secret-id)

echo "verifying kustomize files"
for fixture in $(ls fixtures); do
    base_dir=fixtures/$fixture
    echo "building kustomize file: $base_dir"
    /kustomize build --enable_alpha_plugins $base_dir > $base_dir/result.yaml
    diff $base_dir/result.yaml $base_dir/expected.yaml
done
echo "done"
