package main_test

import (
	"testing"

	"sigs.k8s.io/kustomize/v3/pkg/kusttest"
	plugins_test "sigs.k8s.io/kustomize/v3/pkg/plugins/test"
)

func TestSecretsFromVaultPlugin(t *testing.T) {
	tc := plugins_test.NewEnvForTest(t).Set()
	defer tc.Reset()

	tc.BuildGoPlugin(
		"lumiq.com", "v1", "SecretsFromVault")

	th := kusttest_test.NewKustTestPluginHarness(t, "/app")

	m := th.LoadAndRunGenerator(`
apiVersion: lumiq.com/v1
kind: SecretsFromVault
metadata:
  name: something
name: forbiddenValues
namespace: production
keys:
- ROCKET
- VEGETABLE
`)
	th.AssertActualEqualsExpected(m, `
apiVersion: v1
data:
  ROCKET: U2F0dXJuVg==
  VEGETABLE: Y2Fycm90
kind: Secret
metadata:
  name: forbiddenValues
  namespace: production
type: Opaque
`)
}
