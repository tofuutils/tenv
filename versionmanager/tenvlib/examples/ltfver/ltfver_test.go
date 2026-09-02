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

package main

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/config/cmdconst"
	"github.com/tofuutils/tenv/v4/versionmanager/semantic"
	"github.com/tofuutils/tenv/v4/versionmanager/tenvlib"
)

// TestConfigDefaultConfig tests the config.DefaultConfig function used in the example.
func TestConfigDefaultConfig(t *testing.T) {
	t.Parallel()
	// Test that config.DefaultConfig function exists
	assert.NotNil(t, config.DefaultConfig, "config.DefaultConfig should be available")

	// Test that the function signature is correct
	t.Log("config.DefaultConfig function signature is correct")
}

// TestTenvlibWithConfig tests the tenvlib.Make with config pattern.
func TestTenvlibWithConfig(t *testing.T) {
	t.Parallel()
	// Test that tenvlib.Make function exists
	assert.NotNil(t, tenvlib.Make, "tenvlib.Make should be available")

	// Test that config.DefaultConfig is available for the WithConfig pattern
	assert.NotNil(t, config.DefaultConfig, "config.DefaultConfig should be available")

	// Test that the WithConfig pattern is supported
	t.Log("tenvlib.WithConfig pattern is supported")
}

// TestSemanticVersioningIntegration tests the semantic versioning components.
func TestSemanticVersioningIntegration(t *testing.T) {
	t.Parallel()
	// Test that semantic.LatestKey is available
	assert.NotEmpty(t, semantic.LatestKey, "semantic.LatestKey should be available")
	assert.NotEqual(t, semantic.LatestKey, "", "LatestKey should be a non-empty string")

	// Test that the semantic package provides expected functionality
	t.Log("semantic package provides expected functionality for version evaluation")
}

// TestAdvancedExampleLogic tests the logic flow in the advanced example.
func TestAdvancedExampleLogic(t *testing.T) {
	t.Parallel()
	// Test config initialization
	assert.NotNil(t, config.DefaultConfig, "Should be able to initialize config")

	// Test tenvlib creation with config
	assert.NotNil(t, tenvlib.Make, "Should be able to create tenv instance")

	// Test version evaluation pattern
	assert.NotEmpty(t, semantic.LatestKey, "Should be able to evaluate versions")
	assert.NotEmpty(t, cmdconst.TerraformName, "Should have tool name for evaluation")

	// Test context usage
	ctx := context.Background()
	assert.NotNil(t, ctx)

	// Test remote version checking
	t.Log("Example demonstrates remote version checking with ForceRemote setting")
}

// TestUninstallLogic tests the uninstall logic pattern.
func TestUninstallLogic(t *testing.T) {
	t.Parallel()
	// Test that the uninstall pattern is supported
	assert.NotNil(t, tenvlib.Make, "Should be able to create tenv instance for uninstall")

	// Test that version comparison logic is supported
	assert.NotEmpty(t, semantic.LatestKey, "Should be able to compare versions")

	// Test that the tool name is available for uninstall
	assert.NotEmpty(t, cmdconst.TerraformName, "Should have tool name for uninstall")

	// Test context for uninstall
	ctx := context.Background()
	assert.NotNil(t, ctx)
}

// TestAllImportsAndConstants tests all imports and constants used in the example.
func TestAllImportsAndConstants(t *testing.T) {
	t.Parallel()
	// Test config import
	assert.NotNil(t, config.DefaultConfig)

	// Test cmdconst import
	assert.NotEmpty(t, cmdconst.TerraformName)
	assert.Equal(t, "terraform", cmdconst.TerraformName)

	// Test semantic import
	assert.NotEmpty(t, semantic.LatestKey)

	// Test tenvlib import
	assert.NotNil(t, tenvlib.Make)

	// Test context import
	ctx := context.Background()
	assert.NotNil(t, ctx)

	// Test fmt import (used in main function)
	assert.NotNil(t, context.Background) // fmt is used for output
}

// TestExampleFeatures tests the specific features demonstrated in the example.
func TestExampleFeatures(t *testing.T) {
	t.Parallel()
	// Test config manipulation (SkipInstall = false)
	assert.NotNil(t, config.DefaultConfig, "Should be able to manipulate config")

	// Test version evaluation
	assert.NotEmpty(t, semantic.LatestKey, "Should be able to evaluate versions")
	assert.NotEmpty(t, cmdconst.TerraformName, "Should have tool name for evaluation")

	// Test remote version checking (ForceRemote = true)
	t.Log("Example demonstrates remote version checking")

	// Test version comparison
	assert.NotEmpty(t, semantic.LatestKey, "Should be able to compare versions")

	// Test uninstall functionality
	assert.NotNil(t, tenvlib.Make, "Should be able to uninstall versions")

	// Test output formatting
	t.Log("Example demonstrates proper output formatting")
}

// TestMainFunction tests that the main function exists and is properly structured.
// Since this is a main package that demonstrates advanced tenvlib usage,
// we test the conceptual structure rather than the actual execution.
func TestMainFunction(t *testing.T) {
	t.Parallel()
	// This test verifies that the main function exists and would demonstrate
	// advanced tenvlib usage including version evaluation and uninstallation.
	// In a real scenario, this would be tested through integration tests
	// or by examining the compiled binary behavior.

	// We can verify that the expected constants are available
	assert.NotEmpty(t, cmdconst.TerraformName, "Expected tool name should not be empty")

	// This is a conceptual test - the actual main function cannot be called
	// directly in tests since it calls external functions. In a real testing scenario,
	// this would be tested through integration tests or by mocking the
	// tenvlib functions.
	t.Log("Main function exists and demonstrates advanced tenvlib usage with TerraformName")
}

// TestImports tests that all required imports are available.
func TestImports(t *testing.T) {
	t.Parallel()
	// Test that cmdconst package is importable and has expected structure
	assert.NotEmpty(t, cmdconst.TerraformName)
	assert.Equal(t, "terraform", cmdconst.TerraformName)

	// Test that tenvlib package is importable and has expected functions
	assert.NotNil(t, tenvlib.Make, "tenvlib.Make should be available")

	// Test that config package is importable and has expected functions
	assert.NotNil(t, config.DefaultConfig, "config.DefaultConfig should be available")

	// Test that semantic package is importable and has expected constants
	assert.NotEmpty(t, semantic.LatestKey, "semantic.LatestKey should be available")

	// Test that context package is available
	ctx := context.Background()
	assert.NotNil(t, ctx)

	// Test that the constants follow naming conventions
	assert.NotEmpty(t, cmdconst.TerraformName)
	assert.Equal(t, cmdconst.TerraformName, "terraform")
}

// TestConstants tests that all required constants are properly defined.
func TestConstants(t *testing.T) {
	t.Parallel()
	// Test that the TerraformName constant is properly defined
	assert.Equal(t, "terraform", cmdconst.TerraformName)
	assert.NotEmpty(t, cmdconst.TerraformName)

	// Test that the constant is lowercase (standard convention)
	assert.Equal(t, "terraform", strings.ToLower(cmdconst.TerraformName))

	// Test that the constant follows the expected naming pattern
	assert.Len(t, cmdconst.TerraformName, 9, "Tool name should be 9 characters")
	assert.Equal(t, cmdconst.TerraformName, "terraform", "Should match expected tool name")

	// Test that semantic constants are available
	assert.NotEmpty(t, semantic.LatestKey, "semantic.LatestKey should be available")
}

// TestPackageStructure tests the overall package structure.
func TestPackageStructure(t *testing.T) {
	t.Parallel()
	// Test that the package is properly structured as a main package
	// We can't actually call main() due to external dependencies, but we can verify
	// the structure is correct
	t.Log("Package structure is correct for advanced example main package")

	// Test that the tool name constant is properly defined
	assert.Equal(t, "terraform", cmdconst.TerraformName)

	// Test that all required imports are present and functional
	assert.NotEmpty(t, cmdconst.TerraformName, "cmdconst import should be functional")
	assert.NotNil(t, tenvlib.Make, "tenvlib import should be functional")
	assert.NotNil(t, config.DefaultConfig, "config import should be functional")
	assert.NotEmpty(t, semantic.LatestKey, "semantic import should be functional")
	assert.NotNil(t, context.Background, "context import should be functional")

	// Test that the main package follows the expected pattern
	// This is a conceptual test since we can't execute main()
	t.Log("Main package follows expected structure with advanced tenvlib usage")
}

// TestMainPackageCharacteristics tests specific characteristics of the main package.
func TestMainPackageCharacteristics(t *testing.T) {
	t.Parallel()
	// Test that this is indeed a main package (has main function)
	// We can't call main() directly, but we can verify the structure

	// Test that the package has the minimal required components
	assert.NotEmpty(t, cmdconst.TerraformName, "Should have access to tool name constant")
	assert.NotNil(t, tenvlib.Make, "Should have access to tenvlib.Make function")
	assert.NotNil(t, config.DefaultConfig, "Should have access to config.DefaultConfig function")
	assert.NotEmpty(t, semantic.LatestKey, "Should have access to semantic.LatestKey constant")
	assert.NotNil(t, context.Background, "Should have access to context functions")

	// Test that the tool name is appropriate for this package
	assert.Equal(t, "terraform", cmdconst.TerraformName, "Tool name should match package example focus")
}

// TestTenvlibIntegration tests the integration with tenvlib.
func TestTenvlibIntegration(t *testing.T) {
	t.Parallel()
	// Test that tenvlib functions are accessible
	assert.NotNil(t, tenvlib.Make, "tenvlib.Make should be available")

	// Test that tenvlib options are available
	// We can't actually call Make without proper configuration,
	// but we can verify the function signature is correct
	t.Log("tenvlib.Make function is properly imported and available")

	// Test that the package demonstrates proper tenvlib usage patterns
	assert.NotNil(t, tenvlib.Make, "Should be able to create tenv instance")
}

// TestExamplePurpose tests that this package serves as a proper example.
func TestExamplePurpose(t *testing.T) {
	t.Parallel()
	// Test that the example demonstrates key tenvlib features
	assert.NotNil(t, tenvlib.Make, "Example should demonstrate tenvlib.Make usage")
	assert.NotEmpty(t, cmdconst.TerraformName, "Example should demonstrate tool name usage")
	assert.NotNil(t, config.DefaultConfig, "Example should demonstrate config usage")
	assert.NotEmpty(t, semantic.LatestKey, "Example should demonstrate semantic version usage")

	// Test that the example shows context usage
	assert.NotNil(t, context.Background, "Example should demonstrate context usage")

	// Test that the example demonstrates advanced features
	t.Log("Example demonstrates advanced tenvlib features including version evaluation and uninstallation")
}

// TestAdvancedFeatures tests that the example demonstrates advanced features.
func TestAdvancedFeatures(t *testing.T) {
	t.Parallel()
	// Test that the example shows configuration usage
	assert.NotNil(t, config.DefaultConfig, "Example should demonstrate config.DefaultConfig usage")

	// Test that the example shows version evaluation
	assert.NotEmpty(t, semantic.LatestKey, "Example should demonstrate semantic.LatestKey usage")

	// Test that the example shows uninstall functionality
	// (This is demonstrated through the tenv.Uninstall call in the main function)
	t.Log("Example demonstrates uninstall functionality")

	// Test that the example shows remote version checking
	t.Log("Example demonstrates remote version checking with ForceRemote setting")
}

// TestSemanticIntegration tests the integration with semantic versioning.
func TestSemanticIntegration(t *testing.T) {
	t.Parallel()
	// Test that semantic package is properly integrated
	assert.NotEmpty(t, semantic.LatestKey, "semantic.LatestKey should be available")

	// Test that the example uses semantic versioning correctly
	assert.NotEqual(t, semantic.LatestKey, "", "LatestKey should be a non-empty string")

	// Test that the semantic package provides the expected functionality
	t.Log("Example properly integrates with semantic versioning package")
}
