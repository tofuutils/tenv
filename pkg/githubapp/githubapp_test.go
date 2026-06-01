/*
 *
 * Copyright 2024 tofuutils authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package githubapp

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testAppID = int64(1)

func generateTestPEM(t *testing.T) []byte {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal("failed to generate RSA key:", err)
	}

	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
}

func TestInstallationToken_InvalidAppID(t *testing.T) {
	t.Parallel()

	_, err := InstallationToken(context.Background(), "not-a-number", "", "", "")
	if err == nil {
		t.Fatal("expected error for invalid app ID, got nil")
	}
}

func TestInstallationToken_NoPEM(t *testing.T) {
	t.Parallel()

	_, err := InstallationToken(context.Background(), "1", "", "", "")
	if err == nil {
		t.Fatal("expected error when neither PEM nor PEM file is set, got nil")
	}
}

func TestInstallationToken_PEMFileNotFound(t *testing.T) {
	t.Parallel()

	_, err := InstallationToken(context.Background(), "1", "", "", "/nonexistent/path/key.pem")
	if err == nil {
		t.Fatal("expected error for non-existent PEM file, got nil")
	}
}

func TestInstallationToken_InvalidInstallationID(t *testing.T) {
	t.Parallel()

	keyPEM := generateTestPEM(t)

	_, err := InstallationToken(context.Background(), "1", "not-a-number", string(keyPEM), "")
	if err == nil {
		t.Fatal("expected error for invalid installation ID, got nil")
	}
}

func TestResolveInstallationID_EmptyInstallations(t *testing.T) {
	t.Parallel()

	keyPEM := generateTestPEM(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]any{}) //nolint
	}))
	defer srv.Close()

	_, err := resolveInstallationID(context.Background(), testAppID, "", keyPEM, srv.URL+"/")
	if err == nil {
		t.Fatal("expected error for empty installations list, got nil")
	}
}

func TestResolveInstallationID_APIError(t *testing.T) {
	t.Parallel()

	keyPEM := generateTestPEM(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
	defer srv.Close()

	_, err := resolveInstallationID(context.Background(), testAppID, "", keyPEM, srv.URL+"/")
	if err == nil {
		t.Fatal("expected error for API error response, got nil")
	}
}
