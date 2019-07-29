# Kustomize Plugins

Collections of [kustomize
plugins](https://github.com/kubernetes-sigs/kustomize/tree/master/docs/plugins).

## Usage

### Unit test

```
$ go test ./...
ok      github.com/LumiqBionics/kustomize-plugins/secretsfromvault      1.192s
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

## ArgoCD

`argocd.Dockerfile` builds [Argo CD image](https://hub.docker.com/r/argoproj/argocd) with these
plugins installed.

Images are available in [lumiqbionics/argocd](https://hub.docker.com/r/lumiqbionics/argocd).
