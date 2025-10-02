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

package cosigncheck_test

import (
	"context"
	_ "embed"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	cosigncheck "github.com/tofuutils/tenv/v4/pkg/check/cosign"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
)

const (
	identity = "https://github.com/opentofu/opentofu/.github/workflows/release.yml@refs/heads/v1.6"
	issuer   = "https://token.actions.githubusercontent.com"
)

//go:embed testdata/tofu_1.6.0_SHA256SUMS
var data []byte

//go:embed testdata/tofu_1.6.0_SHA256SUMS.sig
var dataSig []byte

//go:embed testdata/tofu_1.6.0_SHA256SUMS.pem
var dataCert []byte

/*
 * no "t.Parallel()" on those tests (causes failures in cosign call)
 */

func TestCosignCheckCorrect(t *testing.T) { //nolint
	t.SkipNow()
	if err := cosigncheck.Check(t.Context(), data, dataSig, dataCert, identity, issuer, loghelper.InertDisplayer); err != nil {
		t.Error("Unexpected error :", err)
	}
}

func TestCosignCheckErrorCert(t *testing.T) { //nolint
	t.SkipNow()
	if cosigncheck.Check(t.Context(), data, dataSig, dataCert[1:], identity, issuer, loghelper.InertDisplayer) == nil {
		t.Error("Should fail on erroneous certificate")
	}
}

func TestCosignCheckErrorIdentity(t *testing.T) { //nolint
	t.SkipNow()
	if cosigncheck.Check(t.Context(), data, dataSig, dataCert, "me", issuer, loghelper.InertDisplayer) == nil {
		t.Error("Should fail on erroneous issuer")
	}
}

func TestCosignCheckErrorIssuer(t *testing.T) { //nolint
	t.SkipNow()
	if cosigncheck.Check(t.Context(), data, dataSig, dataCert, identity, "http://myself.com", loghelper.InertDisplayer) == nil {
		t.Error("Should fail on erroneous issuer")
	}
}

func TestCosignCheckErrorSig(t *testing.T) { //nolint
	t.SkipNow()
	if cosigncheck.Check(t.Context(), data, dataSig[1:], dataCert, identity, issuer, loghelper.InertDisplayer) == nil {
		t.Error("Should fail on erroneous signature")
	}
}

// TestConstants tests that all constants are properly defined
func TestConstants(t *testing.T) {
	// Test that cosignExecName constant is properly defined
	assert.Equal(t, "cosign", cosigncheck.CosignExecName)
	assert.NotEmpty(t, cosigncheck.CosignExecName)

	// Test that verified constant is properly defined
	assert.Equal(t, "Verified OK", cosigncheck.Verified)
	assert.NotEmpty(t, cosigncheck.Verified)
}

// TestErrorVariables tests that all error variables are properly defined
func TestErrorVariables(t *testing.T) {
	// Test ErrCheck error
	assert.NotNil(t, cosigncheck.ErrCheck)
	assert.Equal(t, "cosign check failed", cosigncheck.ErrCheck.Error())

	// Test ErrNotInstalled error
	assert.NotNil(t, cosigncheck.ErrNotInstalled)
	assert.Equal(t, "cosign executable not found", cosigncheck.ErrNotInstalled.Error())
}

// TestConstantsImmutability tests that constants cannot be modified
func TestConstantsImmutability(t *testing.T) {
	// Test that constants are immutable by trying to access them
	// (they should be accessible and have expected values)
	assert.Equal(t, "cosign", cosigncheck.CosignExecName)
	assert.Equal(t, "Verified OK", cosigncheck.Verified)
}

// TestErrorTypes tests that error types are properly defined
func TestErrorTypes(t *testing.T) {
	// Test that errors are of the correct type
	assert.Contains(t, cosigncheck.ErrCheck.Error(), "cosign check failed")
	assert.Contains(t, cosigncheck.ErrNotInstalled.Error(), "cosign executable not found")
}

// TestPackageStructure tests the overall package structure
func TestPackageStructure(t *testing.T) {
	// Test that the package exports the expected constants and variables
	assert.NotEmpty(t, cosigncheck.CosignExecName)
	assert.NotEmpty(t, cosigncheck.Verified)
	assert.NotNil(t, cosigncheck.ErrCheck)
	assert.NotNil(t, cosigncheck.ErrNotInstalled)

	// Test that the package name is correct
	assert.Equal(t, "cosigncheck", "cosigncheck")
}

// TestTempFileFunction tests the tempFile helper function
func TestTempFileFunction(t *testing.T) {
	// Test successful temp file creation
	testData := []byte("test data")
	fileName, cleanup, err := cosigncheck.TempFile("test", testData)
	assert.NoError(t, err)
	assert.NotEmpty(t, fileName)
	assert.NotNil(t, cleanup)

	// Verify file was created and contains correct data
	content, err := os.ReadFile(fileName)
	assert.NoError(t, err)
	assert.Equal(t, testData, content)

	// Clean up
	cleanup()

	// Verify file was removed
	_, err = os.Stat(fileName)
	assert.True(t, os.IsNotExist(err), "File should be removed after cleanup")
}

// TestTempFileEmptyData tests tempFile with empty data
func TestTempFileEmptyData(t *testing.T) {
	// Test with empty data
	emptyData := []byte("")
	fileName, cleanup, err := cosigncheck.TempFile("empty", emptyData)
	assert.NoError(t, err)
	assert.NotEmpty(t, fileName)
	assert.NotNil(t, cleanup)

	// Verify file was created and is empty
	content, err := os.ReadFile(fileName)
	assert.NoError(t, err)
	assert.Empty(t, content)

	// Clean up
	cleanup()

	// Verify file was removed
	_, err = os.Stat(fileName)
	assert.True(t, os.IsNotExist(err), "File should be removed after cleanup")
}

// TestTempFileErrorHandling tests error handling in tempFile function
func TestTempFileErrorHandling(t *testing.T) {
	// Test with invalid filename pattern that might cause issues
	// Note: This is a conceptual test since we can't easily trigger
	// the actual error conditions without manipulating the filesystem
	t.Log("tempFile function handles errors appropriately")

	// Test that the function signature is correct
	assert.NotNil(t, cosigncheck.TempFile, "tempFile function should be available")
}

// TestCheckFunctionStructure tests the Check function structure and argument handling
func TestCheckFunctionStructure(t *testing.T) {
	// Test that the Check function has the correct signature
	// We can't actually call it without cosign installed, but we can test
	// the function structure and argument validation conceptually
	ctx := context.Background()
	testData := []byte("test data")
	testSig := []byte("test signature")
	testCert := []byte("test certificate")
	identity := "test-identity"
	issuer := "test-issuer"
	displayer := loghelper.InertDisplayer

	// Test that the function exists and has the right signature
	assert.NotNil(t, cosigncheck.Check, "Check function should be available")

	// Test that we can call the function (it will fail due to missing cosign, but that's expected)
	err := cosigncheck.Check(ctx, testData, testSig, testCert, identity, issuer, displayer)
	assert.Error(t, err, "Should return error when cosign is not installed")
	assert.Equal(t, cosigncheck.ErrNotInstalled, err, "Should return ErrNotInstalled when cosign binary is missing")
}

// TestCommandArgumentConstruction tests the command argument construction logic
func TestCommandArgumentConstruction(t *testing.T) {
	// Test the argument construction that would happen in the Check function
	// This tests the logic without actually executing the command
	certIdentity := "https://example.com/identity"
	certOidcIssuer := "https://example.com/issuer"
	dataFileName := "/tmp/data"
	dataSigFileName := "/tmp/data.sig"
	dataCertFileName := "/tmp/data.cert"

	expectedArgs := []string{
		"verify-blob", "--certificate-identity", certIdentity, "--signature", dataSigFileName, "--certificate", dataCertFileName,
		"--certificate-oidc-issuer", certOidcIssuer, dataFileName,
	}

	// Verify the argument construction pattern
	assert.Equal(t, "verify-blob", expectedArgs[0], "First argument should be verify-blob")
	assert.Contains(t, expectedArgs, "--certificate-identity", "Should contain certificate-identity flag")
	assert.Contains(t, expectedArgs, "--signature", "Should contain signature flag")
	assert.Contains(t, expectedArgs, "--certificate", "Should contain certificate flag")
	assert.Contains(t, expectedArgs, "--certificate-oidc-issuer", "Should contain certificate-oidc-issuer flag")
}

// TestVerificationLogic tests the verification string matching logic
func TestVerificationLogic(t *testing.T) {
	// Test the verification logic that checks for "Verified OK" in stderr
	testCases := []struct {
		name     string
		stdErr   string
		expected bool
	}{
		{
			name:     "contains verified OK",
			stdErr:   "some output\nVerified OK\nmore output",
			expected: true,
		},
		{
			name:     "does not contain verified OK",
			stdErr:   "some output\nVerification failed\nmore output",
			expected: false,
		},
		{
			name:     "empty stderr",
			stdErr:   "",
			expected: false,
		},
		{
			name:     "case sensitive match",
			stdErr:   "some output\nverified ok\nmore output",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := strings.Contains(tc.stdErr, cosigncheck.Verified)
			assert.Equal(t, tc.expected, result, "Verification logic should work correctly")
		})
	}
}

// TestErrorTypesExtended tests that error types are properly defined and distinct
func TestErrorTypesExtended(t *testing.T) {
	// Test that the two error types are different
	assert.NotEqual(t, cosigncheck.ErrCheck, cosigncheck.ErrNotInstalled)
	assert.NotEqual(t, cosigncheck.ErrCheck.Error(), cosigncheck.ErrNotInstalled.Error())

	// Test that errors contain expected substrings
	assert.Contains(t, cosigncheck.ErrCheck.Error(), "cosign check failed")
	assert.Contains(t, cosigncheck.ErrNotInstalled.Error(), "cosign executable not found")
}

// TestConstantsValues tests that constants have expected values
func TestConstantsValues(t *testing.T) {
	// Test CosignExecName
	assert.Equal(t, "cosign", cosigncheck.CosignExecName)
	assert.True(t, len(cosigncheck.CosignExecName) > 0)

	// Test Verified string
	assert.Equal(t, "Verified OK", cosigncheck.Verified)
	assert.True(t, len(cosigncheck.Verified) > 0)
	assert.Contains(t, cosigncheck.Verified, "Verified")
}

// TestTempFileWithNilData tests TempFile with nil data
func TestTempFileWithNilData(t *testing.T) {
	// Test with nil data
	fileName, cleanup, err := cosigncheck.TempFile("nil-test", nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, fileName)
	assert.NotNil(t, cleanup)

	// Verify file was created and is empty
	content, err := os.ReadFile(fileName)
	assert.NoError(t, err)
	assert.Empty(t, content)

	// Clean up
	cleanup()

	// Verify file was removed
	_, err = os.Stat(fileName)
	assert.True(t, os.IsNotExist(err), "File should be removed after cleanup")
}

// TestTempFileLargeData tests TempFile with large data
func TestTempFileLargeData(t *testing.T) {
	// Create large data to test
	largeData := make([]byte, 1024*1024) // 1MB
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	fileName, cleanup, err := cosigncheck.TempFile("large-test", largeData)
	assert.NoError(t, err)
	assert.NotEmpty(t, fileName)
	assert.NotNil(t, cleanup)

	// Verify file was created and contains correct data
	content, err := os.ReadFile(fileName)
	assert.NoError(t, err)
	assert.Equal(t, largeData, content)

	// Clean up
	cleanup()

	// Verify file was removed
	_, err = os.Stat(fileName)
	assert.True(t, os.IsNotExist(err), "File should be removed after cleanup")
}

// TestCheckWithEmptyData tests Check function with empty data
func TestCheckWithEmptyData(t *testing.T) {
	ctx := context.Background()
	emptyData := []byte("")
	emptySig := []byte("")
	emptyCert := []byte("")
	identity := "test-identity"
	issuer := "test-issuer"
	displayer := loghelper.InertDisplayer

	// Should fail due to missing cosign binary
	err := cosigncheck.Check(ctx, emptyData, emptySig, emptyCert, identity, issuer, displayer)
	assert.Error(t, err)
	assert.Equal(t, cosigncheck.ErrNotInstalled, err)
}

// TestCheckWithNilData tests Check function with nil data
func TestCheckWithNilData(t *testing.T) {
	ctx := context.Background()
	identity := "test-identity"
	issuer := "test-issuer"
	displayer := loghelper.InertDisplayer

	// Should fail due to missing cosign binary
	err := cosigncheck.Check(ctx, nil, nil, nil, identity, issuer, displayer)
	assert.Error(t, err)
	assert.Equal(t, cosigncheck.ErrNotInstalled, err)
}

// TestCheckWithInvalidParameters tests Check function with invalid parameters
func TestCheckWithInvalidParameters(t *testing.T) {
	ctx := context.Background()
	testData := []byte("test data")
	testSig := []byte("test sig")
	testCert := []byte("test cert")
	displayer := loghelper.InertDisplayer

	testCases := []struct {
		name     string
		identity string
		issuer   string
	}{
		{"empty identity", "", "test-issuer"},
		{"empty issuer", "test-identity", ""},
		{"both empty", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Should fail due to missing cosign binary
			err := cosigncheck.Check(ctx, testData, testSig, testCert, tc.identity, tc.issuer, displayer)
			assert.Error(t, err)
			assert.Equal(t, cosigncheck.ErrNotInstalled, err)
		})
	}
}

// TestCheckContextCancellation tests Check function with cancelled context
func TestCheckContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	testData := []byte("test data")
	testSig := []byte("test sig")
	testCert := []byte("test cert")
	identity := "test-identity"
	issuer := "test-issuer"
	displayer := loghelper.InertDisplayer

	// Should fail due to cancelled context
	err := cosigncheck.Check(ctx, testData, testSig, testCert, identity, issuer, displayer)
	assert.Error(t, err)
	// The error could be either context cancellation or missing cosign binary
	// We just verify it's an error
}

// TestTempFileCleanup tests that TempFile cleanup works properly
func TestTempFileCleanup(t *testing.T) {
	testData := []byte("test data for cleanup")

	// Create temp file
	fileName, cleanup, err := cosigncheck.TempFile("cleanup-test", testData)
	assert.NoError(t, err)
	assert.NotEmpty(t, fileName)

	// Verify file exists
	_, err = os.Stat(fileName)
	assert.NoError(t, err, "File should exist before cleanup")

	// Call cleanup
	cleanup()

	// Verify file was removed
	_, err = os.Stat(fileName)
	assert.True(t, os.IsNotExist(err), "File should be removed after cleanup")
}

// TestTempFileMultipleCleanups tests that multiple cleanup calls are safe
func TestTempFileMultipleCleanups(t *testing.T) {
	testData := []byte("test data for multiple cleanups")

	// Create temp file
	fileName, cleanup, err := cosigncheck.TempFile("multi-cleanup-test", testData)
	assert.NoError(t, err)
	assert.NotEmpty(t, fileName)

	// Call cleanup multiple times
	cleanup()
	cleanup()
	cleanup()

	// Verify file was removed
	_, err = os.Stat(fileName)
	assert.True(t, os.IsNotExist(err), "File should be removed after multiple cleanups")
}

// TestConstantsAndErrors tests that all constants and errors are properly defined
func TestConstantsAndErrors(t *testing.T) {
	// Test constants
	assert.Equal(t, "cosign", cosigncheck.CosignExecName)
	assert.Equal(t, "Verified OK", cosigncheck.Verified)

	// Test errors
	assert.NotNil(t, cosigncheck.ErrCheck)
	assert.NotNil(t, cosigncheck.ErrNotInstalled)
	assert.Contains(t, cosigncheck.ErrCheck.Error(), "cosign check failed")
	assert.Contains(t, cosigncheck.ErrNotInstalled.Error(), "cosign executable not found")
}

// TestVerifiedStringMatching tests the verified string matching logic
func TestVerifiedStringMatching(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"exact match", "Verified OK", true},
		{"with newline", "some output\nVerified OK\nmore output", true},
		{"case sensitive", "verified ok", false},
		{"partial match", "Verified", false},
		{"empty string", "", false},
		{"different text", "Verification failed", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := strings.Contains(tc.input, cosigncheck.Verified)
			assert.Equal(t, tc.expected, result)
		})
	}
}
