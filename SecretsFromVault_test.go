package main_test

import (
	"os"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
	"sigs.k8s.io/kustomize/v3/pkg/kusttest"
	plugins_test "sigs.k8s.io/kustomize/v3/pkg/plugins/test"
)

func init() {
	// Ensure our special envvars are not present
	os.Setenv("VAULT_ADDR", "")
	os.Setenv("VAULT_TOKEN", "")
	os.Setenv("VAULT_ROLE_ID", "")
	os.Setenv("VAULT_SECRET_ID", "")
}

func TestSecretsFromVaultPlugin(t *testing.T) {
	tc := plugins_test.NewEnvForTest(t).Set()
	defer tc.Reset()
	tc.BuildGoPlugin("lumiq.com", "v1", "SecretsFromVault")
	th := kusttest_test.NewKustTestPluginHarness(t, "/app")

	testData := map[string]*vaultapi.Secret{
		"/v1/auth/approle/login": &vaultapi.Secret{
			Auth: &vaultapi.SecretAuth{
				ClientToken: "abcd",
			},
		},
		"/v1/v1/kv/singlekey": &vaultapi.Secret{
			Data: map[string]interface{}{
				"key1": "value1",
			},
		},
		"/v1/v1/kv/multiplekey": &vaultapi.Secret{
			Data: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}

	vaultServer := mockVault(t, testData)
	defer vaultServer.Close()
	os.Setenv("VAULT_ADDR", vaultServer.URL)

	// TODO: test failure case
	// currently failure case can't be tested since KustTestHarness doesn't return error
	// but just failed straight away
	tests := []struct {
		configFile     string
		expectedSecret string
	}{
		// default case
		{
			`
apiVersion: lumiq.com/v1
kind: SecretsFromVault
metadata:
  name: something
name: forbiddenValues
namespace: production
kvPath: v1/kv/multiplekey
`,
			`
apiVersion: v1
data:
  key1: dmFsdWUx
  key2: dmFsdWUy
kind: Secret
metadata:
  name: forbiddenValues
  namespace: production
type: Opaque
`,
		},
		// no namespace
		{
			`
apiVersion: lumiq.com/v1
kind: SecretsFromVault
metadata:
  name: something
name: forbiddenValues
kvPath: v1/kv/singlekey
`,
			`
apiVersion: v1
data:
  key1: dmFsdWUx
kind: Secret
metadata:
  name: forbiddenValues
type: Opaque
`,
		},
	}
	for _, test := range tests {
		generatedSecret := th.LoadAndRunGenerator(test.configFile)
		th.AssertActualEqualsExpected(generatedSecret, test.expectedSecret)
	}
}
