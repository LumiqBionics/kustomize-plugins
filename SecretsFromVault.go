package main

import (
	"errors"
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
	"sigs.k8s.io/kustomize/v3/pkg/ifc"
	"sigs.k8s.io/kustomize/v3/pkg/resmap"
	"sigs.k8s.io/kustomize/v3/pkg/types"
	"sigs.k8s.io/yaml"
)

// generates kubernetes secret from vault KV secret.
// all key-value pairs in vault secret will be used as `data` part in Opaque kubernetes secret.
//
// vault login is done via AppRole with credentials passed via env variables VAULT_ROLE_ID and
// VAULT_SECRET_ID.
//
// given configuration:
// ---
// apiVersion: lumiq.com/v1
// kind: SecretsFromVault
// metadata:
//   name: something
// name: mySecret
// namespace: production
// kvPath: v1/kv/some/data
// GeneratorOptions: {}
//
// if vault KV secret at `v1/kv/some/data` contains:
// | key  | value |
// | ---- | ----- |
// | key1 | data1 |
// | key2 | data2 |
//
// then the resulting secret will be
// ---
// apiVersion: v1
// data:
//   key1: <base64 encoded of data1>
//   key2: <base64 encoded of data2>
// kind: Secret
// metadata:
//   name: mySecret-<suffix>
//   namespace: production
// type: Opaque
type plugin struct {
	rf  *resmap.Factory
	ldr ifc.Loader
	// name of generated secret
	Name string `json:"name" yaml:"name"`
	// namespace of generated secret
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	// custom approle login path
	ApproleLoginPath string `json:"approleLoginPath,omitempty" yaml:"approleLoginPath,omitempty"`
	// path of KV secret
	KvPath           string                 `json:"kvPath,omitempty" yaml:"kvPath,omitempty"`
	GeneratorOptions types.GeneratorOptions `json:"generatorOptions,omitempty" yaml:"generatorOptions,omitempty"`
}

//noinspection GoUnusedGlobalVariable
//nolint: golint
var KustomizePlugin plugin

func (p *plugin) Config(ldr ifc.Loader, rf *resmap.Factory, c []byte) error {
	p.rf = rf
	p.ldr = ldr
	return yaml.Unmarshal(c, p)
}

const (
	approleLoginPath = "auth/approle/login"
)

// login to vault with AppRole
// https://www.vaultproject.io/api/auth/approle/index.html#login-with-approle
func login(path string) (*vaultapi.Client, error) {
	client, err := vaultapi.NewClient(nil)
	if err != nil {
		return nil, err
	}

	options := map[string]interface{}{
		"role_id":   os.Getenv("VAULT_ROLE_ID"),
		"secret_id": os.Getenv("VAULT_SECRET_ID"),
	}

	// retrieve token for the given approle
	secret, err := client.Logical().Write(path, options)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, errors.New("token is empty")
	}
	client.SetToken(secret.Auth.ClientToken)

	return client, nil
}

// read secret data at a given path
func readSecret(client *vaultapi.Client, path string) (map[string]interface{}, error) {
	data, err := client.Logical().Read(path)
	if err != nil {
		return nil, err
	}
	if data == nil || data.Data == nil {
		return nil, fmt.Errorf("no secret available at '%s'", path)
	}

	return data.Data, nil
}

func (p *plugin) Generate() (resmap.ResMap, error) {
	if p.Name == "" {
		return nil, errors.New("Name cannot be empty")
	}

	loginPath := approleLoginPath
	if p.ApproleLoginPath != "" {
		loginPath = p.ApproleLoginPath
	}
	client, err := login(loginPath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	secret, err := readSecret(client, p.KvPath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// initialize secret to be generated
	args := types.SecretArgs{}
	args.Name = p.Name
	args.Namespace = p.Namespace

	// map each key and value in a vault secret to kubernetes secret
	for k, v := range secret {
		secretData := fmt.Sprintf("%s=%s", k, v)
		args.LiteralSources = append(args.LiteralSources, secretData)
	}
	return p.rf.FromSecretArgs(p.ldr, &p.GeneratorOptions, args)
}
