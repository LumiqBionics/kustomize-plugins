package main_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func writeJsonSecret(w http.ResponseWriter, secret *vaultapi.Secret) error {
	data, err := json.Marshal(secret)
	if err != nil {
		return err
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
	return nil
}

func mockVault(t *testing.T, data map[string]*vaultapi.Secret) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path
		secret, ok := data[r.URL.Path]

		if !ok {
			http.Error(w, fmt.Sprintf("path not found: %s", requestPath), http.StatusNotFound)
			return
		}
		if err := writeJsonSecret(w, secret); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}))

	return server
}
