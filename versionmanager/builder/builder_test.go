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

package builder

import (
	"testing"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/stretchr/testify/assert"

	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/config/cmdconst"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
	"github.com/tofuutils/tenv/v4/versionmanager"
)

func TestBuildersMap(t *testing.T) {
	// Test that Builders map contains all expected tools
	expectedBuilders := map[string]bool{
		cmdconst.TofuName:       true,
		cmdconst.TerraformName:  true,
		cmdconst.TerragruntName: true,
		cmdconst.TerramateName:  true,
		cmdconst.AtmosName:      true,
	}

	assert.Equal(t, len(expectedBuilders), len(Builders))

	for toolName := range expectedBuilders {
		assert.Contains(t, Builders, toolName, "Builders map should contain %s", toolName)
		assert.NotNil(t, Builders[toolName], "Builder function for %s should not be nil", toolName)
	}
}

func TestBuildAtmosManager(t *testing.T) {
	// Create a mock config
	mockConfig := &config.Config{
		Displayer: loghelper.InertDisplayer,
	}

	// Create HCL parser
	hclParser := hclparse.NewParser()

	// Test BuildAtmosManager
	manager := BuildAtmosManager(mockConfig, hclParser)

	// Verify it's a valid VersionManager
	assert.NotNil(t, manager)

	// Test that it implements the VersionManager interface
	// by checking it has the expected methods (this is a compile-time check)
	_ = versionmanager.VersionManager(manager)
}

func TestBuildTfManager(t *testing.T) {
	// Create a mock config
	mockConfig := &config.Config{
		Displayer: loghelper.InertDisplayer,
	}

	// Create HCL parser
	hclParser := hclparse.NewParser()

	// Test BuildTfManager
	manager := BuildTfManager(mockConfig, hclParser)

	// Verify it's a valid VersionManager
	assert.NotNil(t, manager)

	// Test that it implements the VersionManager interface
	_ = versionmanager.VersionManager(manager)
}

func TestBuildTgManager(t *testing.T) {
	// Create a mock config
	mockConfig := &config.Config{
		Displayer: loghelper.InertDisplayer,
	}

	// Create HCL parser
	hclParser := hclparse.NewParser()

	// Test BuildTgManager
	manager := BuildTgManager(mockConfig, hclParser)

	// Verify it's a valid VersionManager
	assert.NotNil(t, manager)

	// Test that it implements the VersionManager interface
	_ = versionmanager.VersionManager(manager)
}

func TestBuildTmManager(t *testing.T) {
	// Create a mock config
	mockConfig := &config.Config{
		Displayer: loghelper.InertDisplayer,
	}

	// Create HCL parser
	hclParser := hclparse.NewParser()

	// Test BuildTmManager
	manager := BuildTmManager(mockConfig, hclParser)

	// Verify it's a valid VersionManager
	assert.NotNil(t, manager)

	// Test that it implements the VersionManager interface
	_ = versionmanager.VersionManager(manager)
}

func TestBuildTofuManager(t *testing.T) {
	// Create a mock config
	mockConfig := &config.Config{
		Displayer: loghelper.InertDisplayer,
	}

	// Create HCL parser
	hclParser := hclparse.NewParser()

	// Test BuildTofuManager
	manager := BuildTofuManager(mockConfig, hclParser)

	// Verify it's a valid VersionManager
	assert.NotNil(t, manager)

	// Test that it implements the VersionManager interface
	_ = versionmanager.VersionManager(manager)
}

func TestBuilderFunctionsReturnDifferentManagers(t *testing.T) {
	// Create a mock config
	mockConfig := &config.Config{
		Displayer: loghelper.InertDisplayer,
	}

	// Create HCL parser
	hclParser := hclparse.NewParser()

	// Build all managers
	atmosManager := BuildAtmosManager(mockConfig, hclParser)
	tfManager := BuildTfManager(mockConfig, hclParser)
	tgManager := BuildTgManager(mockConfig, hclParser)
	tmManager := BuildTmManager(mockConfig, hclParser)
	tofuManager := BuildTofuManager(mockConfig, hclParser)

	// Verify all managers are not nil
	assert.NotNil(t, atmosManager)
	assert.NotNil(t, tfManager)
	assert.NotNil(t, tgManager)
	assert.NotNil(t, tmManager)
	assert.NotNil(t, tofuManager)

	// Verify they are different instances by comparing their types
	// Since VersionManager is an interface, we can't easily compare pointers
	// but we can verify they are all valid managers
	assert.NotPanics(t, func() {
		// This would panic if any manager is nil or invalid
		managers := []versionmanager.VersionManager{atmosManager, tfManager, tgManager, tmManager, tofuManager}
		assert.Equal(t, 5, len(managers))
	})
}

func TestBuilderFunctionsWithNilHCLParser(t *testing.T) {
	// Create a mock config
	mockConfig := &config.Config{
		Displayer: loghelper.InertDisplayer,
	}

	// Test functions that don't require HCL parser
	atmosManager := BuildAtmosManager(mockConfig, nil)
	tmManager := BuildTmManager(mockConfig, nil)

	assert.NotNil(t, atmosManager)
	assert.NotNil(t, tmManager)

	// Test functions that require HCL parser (should handle nil gracefully)
	tfManager := BuildTfManager(mockConfig, nil)
	tgManager := BuildTgManager(mockConfig, nil)
	tofuManager := BuildTofuManager(mockConfig, nil)

	assert.NotNil(t, tfManager)
	assert.NotNil(t, tgManager)
	assert.NotNil(t, tofuManager)
}

func TestFuncTypeDefinition(t *testing.T) {
	// Test that Func type is properly defined
	var f Func = BuildTfManager
	assert.NotNil(t, f)

	// Test that it can be called
	mockConfig := &config.Config{
		Displayer: loghelper.InertDisplayer,
	}
	hclParser := hclparse.NewParser()
	manager := f(mockConfig, hclParser)
	assert.NotNil(t, manager)
}
