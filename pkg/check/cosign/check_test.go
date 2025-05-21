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

package cosigncheck

import (
	"os"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"

	"github.com/tofuutils/tenv/v4/pkg/loghelper"
)

func TestCheck(t *testing.T) {
	tests := []struct {
		name           string
		data           []byte
		dataSig        []byte
		dataCert       []byte
		certIdentity   string
		certOidcIssuer string
		expectError    bool
	}{
		{
			name:           "valid signature",
			data:           []byte("test content"),
			dataSig:        []byte("test signature"),
			dataCert:       []byte("test certificate"),
			certIdentity:   "test@example.com",
			certOidcIssuer: "https://accounts.example.com",
			expectError:    false,
		},
		{
			name:           "invalid signature",
			data:           []byte("test content"),
			dataSig:        []byte("invalid signature"),
			dataCert:       []byte("test certificate"),
			certIdentity:   "test@example.com",
			certOidcIssuer: "https://accounts.example.com",
			expectError:    true,
		},
		{
			name:           "empty certificate",
			data:           []byte("test content"),
			dataSig:        []byte("test signature"),
			dataCert:       []byte(""),
			certIdentity:   "test@example.com",
			certOidcIssuer: "https://accounts.example.com",
			expectError:    true,
		},
		{
			name:           "empty signature",
			data:           []byte("test content"),
			dataSig:        []byte(""),
			dataCert:       []byte("test certificate"),
			certIdentity:   "test@example.com",
			certOidcIssuer: "https://accounts.example.com",
			expectError:    true,
		},
		{
			name:           "empty data",
			data:           []byte(""),
			dataSig:        []byte("test signature"),
			dataCert:       []byte("test certificate"),
			certIdentity:   "test@example.com",
			certOidcIssuer: "https://accounts.example.com",
			expectError:    true,
		},
		{
			name:           "invalid certificate format",
			data:           []byte("test content"),
			dataSig:        []byte("test signature"),
			dataCert:       []byte("invalid certificate format"),
			certIdentity:   "test@example.com",
			certOidcIssuer: "https://accounts.example.com",
			expectError:    true,
		},
		{
			name:           "mismatched identity",
			data:           []byte("test content"),
			dataSig:        []byte("test signature"),
			dataCert:       []byte("test certificate"),
			certIdentity:   "wrong@example.com",
			certOidcIssuer: "https://accounts.example.com",
			expectError:    true,
		},
		{
			name:           "mismatched issuer",
			data:           []byte("test content"),
			dataSig:        []byte("test signature"),
			dataCert:       []byte("test certificate"),
			certIdentity:   "test@example.com",
			certOidcIssuer: "https://wrong.example.com",
			expectError:    true,
		},
	}

	// Create test configuration
	logger := hclog.New(&hclog.LoggerOptions{
		Output: os.Stderr,
		Level:  hclog.Info,
	})
	displayer := loghelper.MakeBasicDisplayer(logger, loghelper.StdDisplay)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run check
			err := Check(tt.data, tt.dataSig, tt.dataCert, tt.certIdentity, tt.certOidcIssuer, displayer)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckWithNilDisplayer(t *testing.T) {
	// Test with nil displayer
	err := Check([]byte("test"), []byte("sig"), []byte("cert"), "test@example.com", "https://accounts.example.com", nil)
	assert.Error(t, err)
}

func TestCheckWithInvalidCertificate(t *testing.T) {
	// Test with invalid certificate format
	logger := hclog.New(&hclog.LoggerOptions{
		Output: os.Stderr,
		Level:  hclog.Info,
	})
	displayer := loghelper.MakeBasicDisplayer(logger, loghelper.StdDisplay)

	// Create a certificate with invalid format
	invalidCert := []byte("-----BEGIN CERTIFICATE-----\ninvalid\n-----END CERTIFICATE-----")

	err := Check([]byte("test"), []byte("sig"), invalidCert, "test@example.com", "https://accounts.example.com", displayer)
	assert.Error(t, err)
}
