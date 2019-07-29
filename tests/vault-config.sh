#!/bin/sh

TIMEOUT=${TIMEOUT:-3}
TRIES=${TRIES:-10}
WAIT=${WAIT:-10}

set -uo pipefail

for i in $(seq $TRIES) ; do
    vault status
    result=$?
    if [ $result -eq 0 ] ; then
        if [ $i -eq $TRIES ]; then
            exit 1
        fi
        break
    fi
    echo "vault is not ready, sleeping for $WAIT seconds..."
    sleep $WAIT
done

set -e

vault login myroot
vault secrets enable -path=v1/kv -version=1 kv

echo "generating secrets"

for fixture in $(ls fixtures); do
    for secret in $(ls fixtures/$fixture/secret*.yaml); do
        kv_path=$(cat $secret | awk '/kvPath/ {print $NF}')

        # create dummy secret for all applications
        # since the plugin doesn't care if the secret has
        # all the required KV
        vault write $kv_path dummy=value
    done
done

echo "setting up credentials"
vault auth enable approle
cat <<EOF | vault policy write kustomize -
path "/v1/kv/*" {
policy = "read"
}
EOF
vault write auth/approle/role/kustomize policies=kustomize period=1h

# get approle credentials
vault read auth/approle/role/kustomize/role-id \
    | grep role_id \
    | awk '{print $NF}' \
    | tee /credentials/vault-role-id > /dev/null
vault write -f auth/approle/role/kustomize/secret-id \
    | grep "secret_id " \
    | awk '{print $NF}' \
    | tee /credentials/vault-secret-id > /dev/null
ls -l /credentials
echo "done"
