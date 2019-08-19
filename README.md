# Kustomize Plugins

Collections of [kustomize
plugins](https://github.com/kubernetes-sigs/kustomize/tree/master/docs/plugins).

## Usage

### Unit test

```
$ go test ./...
ok      github.com/LumiqBionics/kustomize-plugins/secretsfromvault      1.192s
```

### Acceptance test

```
$ cd tests/
$ docker-compose run test
Creating network "tests_default" with the default driver
Creating volume "tests_credentials" with default driver
Creating tests_vaultconfig_1 ... done
Creating tests_vault_1       ... done
cat: /credentials/vault-role-id: No such file or directory
vault is not ready, sleeping for 5 seconds...
cat: /credentials/vault-role-id: No such file or directory
vault is not ready, sleeping for 5 seconds...
f2e6cbb2-312d-17fa-c015-2256fe256112
13b2617e-8a06-9153-2edc-e08fc85a1cf6
verifying kustomize files
building kustomize file: fixtures/base
building kustomize file: fixtures/no-suffix
done
```

### Building and using the plugin locally

See [Go Plugin Guided Example for Linux](
https://github.com/kubernetes-sigs/kustomize/blob/master/docs/plugins/goPluginGuidedExample.md) for
more information.

TLDR:
```sh
# set variables
export GOPATH=$(mktemp -d)
go install sigs.k8s.io/kustomize/v3/cmd/kustomize
export XDG_CONFIG_HOME=$GOPATH/xdg
export API_VERSION=lumiq.com/v1

# build plugins
for plugin in SecretsFromVault; do
    lowername=$(echo $plugin | awk '{print tolower($0)}')
    plugin_home=$XDG_CONFIG_HOME/kustomize/plugin/$API_VERSION/$lowername
    mkdir -p $plugin_home
    go build -buildmode plugin -o $plugin.so $plugin.go
    mv $plugin.so $plugin_home/
done

# test with kustomize
cat <<EOF > kustomization.yaml
generators:
- secret.yaml
EOF

cat <<EOF > secret.yaml
apiVersion: lumiq.com/v1
kind: SecretsFromVault
metadata:
  name: something
name: my-secret
kvPath: some/existing/vault/kv/secret
EOF

export VAULT_ADDR="<vault address>"
export VAULT_ROLE_ID="<approle role id>"
export VAULT_SECRET_ID="<approle secret id>"
$GOPATH/bin/kustomize build --enable_alpha_plugins
```

### Docker image

This is available as docker image in [lumiqbionics/kustomize](
https://hub.docker.com/r/lumiqbionics/kustomize). It can be run with:

```sh
export VAULT_ADDR=<vault address>
# login with your method of choice
vault login <options>
docker run -e VAULT_ADDR=$VAULT_ADDR \
    -e VAULT_TOKEN=$(cat ~/.vault-token) \
    --rm -v ${PWD}:/code -ti lumiqbionics/kustomize \
    build --enable_alpha_plugins /code/<path to kustomization file>
```

## ArgoCD

`argocd.Dockerfile` builds [Argo CD image](https://hub.docker.com/r/argoproj/argocd) with these
plugins installed.

Images are available in [lumiqbionics/argocd](https://hub.docker.com/r/lumiqbionics/argocd).
