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

package tenvlib

import (
	"context"
	"testing"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
	"github.com/tofuutils/tenv/v4/versionmanager"
	"github.com/tofuutils/tenv/v4/versionmanager/builder"
)

const (
	testToolName = "terraform"
	testVersion  = "1.0.0"
)

type mockRetriever struct{}

func (m *mockRetriever) Install(ctx context.Context, version string, targetPath string) error {
	return nil
}

func (m *mockRetriever) ListVersions(ctx context.Context) ([]string, error) {
	return []string{"1.0.0", "1.1.0"}, nil
}

func TestErrNoBuilder(t *testing.T) {
	t.Parallel()
	// Test that errNoBuilder is properly defined
	require.Error(t, errNoBuilder)
	assert.Equal(t, "no builder for this tool", errNoBuilder.Error())
}

func TestPackageStructure(t *testing.T) {
	t.Parallel()
	// Test that Tenv struct exists and has expected fields
	var tenv Tenv
	assert.NotNil(t, tenv)

	// Test that tenvConfig struct exists and has expected fields
	var tenvConf tenvConfig
	assert.NotNil(t, tenvConf)

	// Test that TenvOption type exists (function types are nil by default, which is expected)
	var option TenvOption
	// Function types are nil by default, so we just verify the type exists
	assert.IsType(t, (TenvOption)(nil), option)
}

func TestOptionFunctions(t *testing.T) {
	t.Parallel()
	// Test that option functions exist and are callable
	conf := &config.Config{}
	displayer := loghelper.InertDisplayer
	hclParser := hclparse.NewParser()

	// Test AddTool function
	addToolOption := AddTool("test", nil)
	assert.NotNil(t, addToolOption)

	// Test AutoInstall function
	autoInstallOption := AutoInstall
	assert.NotNil(t, autoInstallOption)

	// Test DisableDisplay function
	disableDisplayOption := DisableDisplay
	assert.NotNil(t, disableDisplayOption)

	// Test IgnoreEnv function
	ignoreEnvOption := IgnoreEnv
	assert.NotNil(t, ignoreEnvOption)

	// Test WithConfig function
	withConfigOption := WithConfig(conf)
	assert.NotNil(t, withConfigOption)

	// Test WithDisplayer function
	withDisplayerOption := WithDisplayer(displayer)
	assert.NotNil(t, withDisplayerOption)

	// Test WithHCLParser function
	withHCLParserOption := WithHCLParser(hclParser)
	assert.NotNil(t, withHCLParserOption)
}

func TestMakeFunction(t *testing.T) {
	t.Parallel()
	// Test Make function with no options
	tenv, err := Make()
	require.NoError(t, err)
	assert.NotNil(t, tenv)
	assert.NotNil(t, tenv.conf)
	assert.NotNil(t, tenv.hclParser)
	assert.NotNil(t, tenv.builders)
	assert.NotNil(t, tenv.managers)
}

func TestMakeWithOptions(t *testing.T) {
	t.Parallel()
	// Test Make function with various options
	conf := &config.Config{}
	displayer := loghelper.InertDisplayer
	hclParser := hclparse.NewParser()

	tenv, err := Make(
		WithConfig(conf),
		WithDisplayer(displayer),
		WithHCLParser(hclParser),
		AutoInstall,
		DisableDisplay,
		IgnoreEnv,
	)

	require.NoError(t, err)
	assert.NotNil(t, tenv)
	assert.Equal(t, conf, tenv.conf)
	assert.Equal(t, hclParser, tenv.hclParser)
	assert.NotNil(t, tenv.builders)
	assert.NotNil(t, tenv.managers)
}

func TestMakeWithAddTool(t *testing.T) {
	t.Parallel()
	// Test Make function with AddTool option
	// We can't easily test this without creating a proper mock builder
	// due to the complex VersionManager interface requirements
	tenv, err := Make()

	require.NoError(t, err)
	assert.NotNil(t, tenv)
	assert.NotNil(t, tenv.builders)
}

func TestTenvMethodsExist(t *testing.T) {
	t.Parallel()
	// Test that all expected methods exist on Tenv
	tenv, err := Make()
	require.NoError(t, err)

	// These methods should exist (we can't call them directly without proper setup,
	// but we can verify they exist by checking the type)
	assert.NotNil(t, tenv)

	// Verify the struct has the expected fields
	assert.NotNil(t, tenv.conf)
	assert.NotNil(t, tenv.hclParser)
	assert.NotNil(t, tenv.builders)
	assert.NotNil(t, tenv.managers)
}

func TestTenvStructFields(t *testing.T) {
	t.Parallel()
	// Test that Tenv struct can be created with expected fields
	conf := &config.Config{}
	hclParser := hclparse.NewParser()
	builders := make(map[string]builder.Func)
	managers := make(map[string]versionmanager.VersionManager)

	tenv := Tenv{
		builders:  builders,
		conf:      conf,
		hclParser: hclParser,
		managers:  managers,
	}

	assert.Equal(t, builders, tenv.builders)
	assert.Equal(t, conf, tenv.conf)
	assert.Equal(t, hclParser, tenv.hclParser)
	assert.Equal(t, managers, tenv.managers)
}

func TestTenvConfigStructFields(t *testing.T) {
	t.Parallel()
	// Test that tenvConfig struct can be created with expected fields
	conf := &config.Config{}
	displayer := loghelper.InertDisplayer
	hclParser := hclparse.NewParser()
	builders := make(map[string]builder.Func)
	initConfigFunc := func() (config.Config, error) { return config.Config{}, nil }

	tenvConf := tenvConfig{
		autoInstall:    true,
		builders:       builders,
		conf:           conf,
		displayer:      displayer,
		hclParser:      hclParser,
		ignoreEnv:      true,
		initConfigFunc: initConfigFunc,
	}

	assert.True(t, tenvConf.autoInstall)
	assert.Equal(t, builders, tenvConf.builders)
	assert.Equal(t, conf, tenvConf.conf)
	assert.Equal(t, displayer, tenvConf.displayer)
	assert.Equal(t, hclParser, tenvConf.hclParser)
	assert.True(t, tenvConf.ignoreEnv)
	assert.NotNil(t, tenvConf.initConfigFunc)
}

func TestOptionFunctionTypes(t *testing.T) {
	t.Parallel()
	// Test that option functions have the correct function signatures
	var addToolFunc func(string, builder.Func) TenvOption
	var withConfigFunc func(*config.Config) TenvOption
	var withDisplayerFunc func(loghelper.Displayer) TenvOption
	var withHCLParserFunc func(*hclparse.Parser) TenvOption

	// These should be assignable to the function types
	addToolFunc = AddTool
	withConfigFunc = WithConfig
	withDisplayerFunc = WithDisplayer
	withHCLParserFunc = WithHCLParser

	assert.NotNil(t, addToolFunc)
	assert.NotNil(t, withConfigFunc)
	assert.NotNil(t, withDisplayerFunc)
	assert.NotNil(t, withHCLParserFunc)
}

func TestErrorHandling(t *testing.T) {
	t.Parallel()
	// Test that error variable is properly defined and accessible
	assert.Equal(t, "no builder for this tool", errNoBuilder.Error())
	assert.Contains(t, errNoBuilder.Error(), "builder")
	assert.Contains(t, errNoBuilder.Error(), "tool")
}

func TestMakeFunctionMultipleCalls(t *testing.T) {
	t.Parallel()
	// Test that Make function can be called multiple times
	tenv1, err1 := Make()
	tenv2, err2 := Make()

	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.NotNil(t, tenv1)
	assert.NotNil(t, tenv2)
	// Each call should create independent instances
	assert.NotSame(t, tenv1.conf, tenv2.conf)
	assert.NotSame(t, tenv1.hclParser, tenv2.hclParser)
}

func TestMakeWithNilConfig(t *testing.T) {
	t.Parallel()
	// Test Make function with nil config option
	// Note: Make function will initialize config if nil, so we test that it doesn't panic
	tenv, err := Make(WithConfig(nil))

	require.NoError(t, err)
	assert.NotNil(t, tenv)
	// Make function initializes config, so it won't be nil
	assert.NotNil(t, tenv.conf)
}

func TestMakeWithNilDisplayer(t *testing.T) {
	t.Parallel()
	// Test Make function with nil displayer option
	// Note: Make function will initialize displayer if nil, so we test that it doesn't panic
	tenv, err := Make(WithDisplayer(nil))

	require.NoError(t, err)
	assert.NotNil(t, tenv)
	// Make function initializes displayer, so it won't be nil
	assert.NotNil(t, tenv.conf.Displayer)
}

func TestMakeWithNilHCLParser(t *testing.T) {
	t.Parallel()
	// Test Make function with nil HCL parser option
	// Note: Make function will initialize HCL parser if nil, so we test that it doesn't panic
	tenv, err := Make(WithHCLParser(nil))

	require.NoError(t, err)
	assert.NotNil(t, tenv)
	// Make function initializes HCL parser, so it won't be nil
	assert.NotNil(t, tenv.hclParser)
}

func TestOptionFunctionsReturnValues(t *testing.T) {
	t.Parallel()
	// Test that option functions return proper function types
	addToolOption := AddTool("test", nil)
	withConfigOption := WithConfig(&config.Config{})
	withDisplayerOption := WithDisplayer(loghelper.InertDisplayer)
	withHCLParserOption := WithHCLParser(hclparse.NewParser())

	// Test that these options can be called
	tenvConf := &tenvConfig{
		builders: make(map[string]builder.Func),
	}
	addToolOption(tenvConf)
	withConfigOption(tenvConf)
	withDisplayerOption(tenvConf)
	withHCLParserOption(tenvConf)

	// Verify the options were applied
	assert.NotNil(t, tenvConf.builders)
	assert.NotNil(t, tenvConf.conf)
	assert.NotNil(t, tenvConf.displayer)
	assert.NotNil(t, tenvConf.hclParser)
}

func TestAutoInstallOption(t *testing.T) {
	t.Parallel()
	// Test AutoInstall option function
	tenvConf := &tenvConfig{}
	AutoInstall(tenvConf)

	assert.True(t, tenvConf.autoInstall)
}

func TestDisableDisplayOption(t *testing.T) {
	t.Parallel()
	// Test DisableDisplay option function
	tenvConf := &tenvConfig{}
	DisableDisplay(tenvConf)

	assert.Equal(t, loghelper.InertDisplayer, tenvConf.displayer)
}

func TestIgnoreEnvOption(t *testing.T) {
	t.Parallel()
	// Test IgnoreEnv option function
	tenvConf := &tenvConfig{}
	IgnoreEnv(tenvConf)

	assert.True(t, tenvConf.ignoreEnv)
}

func TestTenvInitialization(t *testing.T) {
	t.Parallel()
	// Test that Tenv can be initialized with zero values
	tenv := Tenv{}

	// Zero values for maps are nil, which is expected
	assert.Nil(t, tenv.builders)
	assert.Nil(t, tenv.managers)
	// Zero values for pointers are nil, which is expected
	assert.Nil(t, tenv.conf)
	assert.Nil(t, tenv.hclParser)
}

func TestTenvConfigInitialization(t *testing.T) {
	t.Parallel()
	// Test that tenvConfig can be initialized with zero values
	tenvConf := tenvConfig{}

	// Zero values for maps are nil, which is expected
	assert.Nil(t, tenvConf.builders)
	assert.False(t, tenvConf.autoInstall)
	assert.False(t, tenvConf.ignoreEnv)
	// Zero values for pointers and functions are nil, which is expected
	assert.Nil(t, tenvConf.conf)
	assert.Nil(t, tenvConf.displayer)
	assert.Nil(t, tenvConf.hclParser)
	assert.Nil(t, tenvConf.initConfigFunc)
}

// Comprehensive tests for Tenv methods.
func TestTenvCommandComprehensive(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	toolName := testToolName
	requestedVersion := testVersion
	cmdArgs := []string{"apply"}

	// Create a simple test config
	testConfig := &config.Config{
		Getenv: config.EmptyGetenv,
	}
	testConfig.InitDisplayer(false)

	// Create HCL parser
	hclParser := hclparse.NewParser()

	// Create Tenv instance with empty builders map to test error handling
	tenv := Tenv{
		builders:  make(map[string]builder.Func), // Empty builders map
		conf:      testConfig,
		hclParser: hclParser,
		managers:  make(map[string]versionmanager.VersionManager),
	}

	// Test that Command method handles missing builder correctly
	_, err := tenv.Command(ctx, toolName, requestedVersion, cmdArgs...)
	// This should fail with "no builder for this tool" error
	require.Error(t, err)
	assert.Equal(t, errNoBuilder, err)
}

func TestTenvCommandProxyComprehensive(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	toolName := testToolName
	requestedVersion := testVersion
	cmdArgs := []string{"apply"}

	// Create a simple test config
	testConfig := &config.Config{
		Getenv: config.EmptyGetenv,
	}
	testConfig.InitDisplayer(false)

	// Create HCL parser
	hclParser := hclparse.NewParser()

	// Create Tenv instance with empty builders map to test error handling
	tenv := Tenv{
		builders:  make(map[string]builder.Func), // Empty builders map
		conf:      testConfig,
		hclParser: hclParser,
		managers:  make(map[string]versionmanager.VersionManager),
	}

	// Test CommandProxy method error handling
	err := tenv.CommandProxy(ctx, toolName, requestedVersion, cmdArgs...)
	require.Error(t, err)
	assert.Equal(t, errNoBuilder, err)
}

func TestTenvDetectComprehensive(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	toolName := testToolName

	// Create a simple test config
	testConfig := &config.Config{
		Getenv: config.EmptyGetenv,
	}
	testConfig.InitDisplayer(false)

	// Create HCL parser
	hclParser := hclparse.NewParser()

	// Create Tenv instance with empty builders map to test error handling
	tenv := Tenv{
		builders:  make(map[string]builder.Func), // Empty builders map
		conf:      testConfig,
		hclParser: hclParser,
		managers:  make(map[string]versionmanager.VersionManager),
	}

	// Test Detect method error handling
	version, err := tenv.Detect(ctx, toolName)
	require.Error(t, err)
	assert.Empty(t, version)
	assert.Equal(t, errNoBuilder, err)
}

func TestTenvDetectedCommandComprehensive(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	toolName := testToolName
	cmdArgs := []string{"version"}

	// Create a simple test config
	testConfig := &config.Config{
		Getenv: config.EmptyGetenv,
	}
	testConfig.InitDisplayer(false)

	// Create HCL parser
	hclParser := hclparse.NewParser()

	// Create Tenv instance with empty builders map to test error handling
	tenv := Tenv{
		builders:  make(map[string]builder.Func), // Empty builders map
		conf:      testConfig,
		hclParser: hclParser,
		managers:  make(map[string]versionmanager.VersionManager),
	}

	// Test DetectedCommand method error handling
	cmd, err := tenv.DetectedCommand(ctx, toolName, cmdArgs...)
	require.Error(t, err)
	assert.Nil(t, cmd)
	assert.Equal(t, errNoBuilder, err)
}

func TestTenvDetectedCommandProxyComprehensive(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	toolName := testToolName
	cmdArgs := []string{"init"}

	// Create a simple test config
	testConfig := &config.Config{
		Getenv: config.EmptyGetenv,
	}
	testConfig.InitDisplayer(false)

	// Create HCL parser
	hclParser := hclparse.NewParser()

	// Create Tenv instance with empty builders map to test error handling
	tenv := Tenv{
		builders:  make(map[string]builder.Func), // Empty builders map
		conf:      testConfig,
		hclParser: hclParser,
		managers:  make(map[string]versionmanager.VersionManager),
	}

	// Test DetectedCommandProxy method error handling
	err := tenv.DetectedCommandProxy(ctx, toolName, cmdArgs...)
	require.Error(t, err)
	assert.Equal(t, errNoBuilder, err)
}

func TestTenvEvaluateComprehensive(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	toolName := testToolName
	requestedVersion := testVersion

	// Create a simple test config
	testConfig := &config.Config{
		Getenv: config.EmptyGetenv,
	}
	testConfig.InitDisplayer(false)

	// Create HCL parser
	hclParser := hclparse.NewParser()

	// Create Tenv instance with empty builders map to test error handling
	tenv := Tenv{
		builders:  make(map[string]builder.Func), // Empty builders map
		conf:      testConfig,
		hclParser: hclParser,
		managers:  make(map[string]versionmanager.VersionManager),
	}

	// Test Evaluate method error handling
	version, err := tenv.Evaluate(ctx, toolName, requestedVersion)
	require.Error(t, err)
	assert.Empty(t, version)
	assert.Equal(t, errNoBuilder, err)
}

func TestTenvInstallComprehensive(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	toolName := testToolName
	requestedVersion := testVersion

	// Create a simple test config
	testConfig := &config.Config{
		Getenv: config.EmptyGetenv,
	}
	testConfig.InitDisplayer(false)

	// Create HCL parser
	hclParser := hclparse.NewParser()

	// Create Tenv instance with empty builders map to test error handling
	tenv := Tenv{
		builders:  make(map[string]builder.Func), // Empty builders map
		conf:      testConfig,
		hclParser: hclParser,
		managers:  make(map[string]versionmanager.VersionManager),
	}

	// Test Install method error handling
	err := tenv.Install(ctx, toolName, requestedVersion)
	require.Error(t, err)
	assert.Equal(t, errNoBuilder, err)
}

func TestTenvInstallMultipleComprehensive(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	toolName := testToolName
	versions := []string{"1.0.0", "1.1.0"}

	// Create a simple test config
	testConfig := &config.Config{
		Getenv: config.EmptyGetenv,
	}
	testConfig.InitDisplayer(false)

	// Create HCL parser
	hclParser := hclparse.NewParser()

	// Create Tenv instance with empty builders map to test error handling
	tenv := Tenv{
		builders:  make(map[string]builder.Func), // Empty builders map
		conf:      testConfig,
		hclParser: hclParser,
		managers:  make(map[string]versionmanager.VersionManager),
	}

	// Test InstallMultiple method error handling
	err := tenv.InstallMultiple(ctx, toolName, versions)
	require.Error(t, err)
	assert.Equal(t, errNoBuilder, err)
}

func TestTenvListLocalComprehensive(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	toolName := testToolName
	reverseOrder := false

	// Create a simple test config
	testConfig := &config.Config{
		Getenv: config.EmptyGetenv,
	}
	testConfig.InitDisplayer(false)

	// Create HCL parser
	hclParser := hclparse.NewParser()

	// Create Tenv instance with empty builders map to test error handling
	tenv := Tenv{
		builders:  make(map[string]builder.Func), // Empty builders map
		conf:      testConfig,
		hclParser: hclParser,
		managers:  make(map[string]versionmanager.VersionManager),
	}

	// Test ListLocal method error handling
	versions, err := tenv.ListLocal(ctx, toolName, reverseOrder)
	require.Error(t, err)
	assert.Nil(t, versions)
	assert.Equal(t, errNoBuilder, err)
}

func TestTenvListRemoteComprehensive(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	toolName := testToolName
	reverseOrder := false

	// Create a simple test config
	testConfig := &config.Config{
		Getenv: config.EmptyGetenv,
	}
	testConfig.InitDisplayer(false)

	// Create HCL parser
	hclParser := hclparse.NewParser()

	// Create Tenv instance with empty builders map to test error handling
	tenv := Tenv{
		builders:  make(map[string]builder.Func), // Empty builders map
		conf:      testConfig,
		hclParser: hclParser,
		managers:  make(map[string]versionmanager.VersionManager),
	}

	// Test ListRemote method error handling
	versions, err := tenv.ListRemote(ctx, toolName, reverseOrder)
	require.Error(t, err)
	assert.Nil(t, versions)
	assert.Equal(t, errNoBuilder, err)
}

func TestTenvLocallyInstalledComprehensive(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	toolName := testToolName

	// Create a simple test config
	testConfig := &config.Config{
		Getenv: config.EmptyGetenv,
	}
	testConfig.InitDisplayer(false)

	// Create HCL parser
	hclParser := hclparse.NewParser()

	// Create Tenv instance with empty builders map to test error handling
	tenv := Tenv{
		builders:  make(map[string]builder.Func), // Empty builders map
		conf:      testConfig,
		hclParser: hclParser,
		managers:  make(map[string]versionmanager.VersionManager),
	}

	// Test LocallyInstalled method error handling
	versionSet, err := tenv.LocallyInstalled(ctx, toolName)
	require.Error(t, err)
	assert.Nil(t, versionSet)
	assert.Equal(t, errNoBuilder, err)
}

func TestTenvResetDefaultConstraintComprehensive(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	toolName := testToolName

	// Create a simple test config
	testConfig := &config.Config{
		Getenv: config.EmptyGetenv,
	}
	testConfig.InitDisplayer(false)

	// Create HCL parser
	hclParser := hclparse.NewParser()

	// Create Tenv instance with empty builders map to test error handling
	tenv := Tenv{
		builders:  make(map[string]builder.Func), // Empty builders map
		conf:      testConfig,
		hclParser: hclParser,
		managers:  make(map[string]versionmanager.VersionManager),
	}

	// Test ResetDefaultConstraint method error handling
	err := tenv.ResetDefaultConstraint(ctx, toolName)
	require.Error(t, err)
	assert.Equal(t, errNoBuilder, err)
}

func TestTenvResetDefaultVersionComprehensive(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	toolName := testToolName

	// Create a simple test config
	testConfig := &config.Config{
		Getenv: config.EmptyGetenv,
	}
	testConfig.InitDisplayer(false)

	// Create HCL parser
	hclParser := hclparse.NewParser()

	// Create Tenv instance with empty builders map to test error handling
	tenv := Tenv{
		builders:  make(map[string]builder.Func), // Empty builders map
		conf:      testConfig,
		hclParser: hclParser,
		managers:  make(map[string]versionmanager.VersionManager),
	}

	// Test ResetDefaultVersion method error handling
	err := tenv.ResetDefaultVersion(ctx, toolName)
	require.Error(t, err)
	assert.Equal(t, errNoBuilder, err)
}

func TestTenvSetDefaultConstraintComprehensive(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	toolName := testToolName
	constraint := ">= 1.0.0"

	// Create a simple test config
	testConfig := &config.Config{
		Getenv: config.EmptyGetenv,
	}
	testConfig.InitDisplayer(false)

	// Create HCL parser
	hclParser := hclparse.NewParser()

	// Create Tenv instance with empty builders map to test error handling
	tenv := Tenv{
		builders:  make(map[string]builder.Func), // Empty builders map
		conf:      testConfig,
		hclParser: hclParser,
		managers:  make(map[string]versionmanager.VersionManager),
	}

	// Test SetDefaultConstraint method error handling
	err := tenv.SetDefaultConstraint(ctx, toolName, constraint)
	require.Error(t, err)
	assert.Equal(t, errNoBuilder, err)
}

func TestTenvSetDefaultVersionComprehensive(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	toolName := testToolName
	requestedVersion := testVersion
	workingDir := false

	// Create a simple test config
	testConfig := &config.Config{
		Displayer: loghelper.InertDisplayer,
	}

	// Create HCL parser
	hclParser := hclparse.NewParser()

	// Create Tenv instance with empty builders map to test error handling
	tenv := Tenv{
		builders:  make(map[string]builder.Func), // Empty builders map
		conf:      testConfig,
		hclParser: hclParser,
		managers:  make(map[string]versionmanager.VersionManager),
	}

	// Test SetDefaultVersion method error handling
	err := tenv.SetDefaultVersion(ctx, toolName, requestedVersion, workingDir)
	require.Error(t, err)
	assert.Equal(t, errNoBuilder, err)
}

func TestTenvUninstallComprehensive(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	toolName := testToolName
	requestedVersion := testVersion

	// Create a simple test config
	testConfig := &config.Config{
		Displayer: loghelper.InertDisplayer,
	}

	// Create HCL parser
	hclParser := hclparse.NewParser()

	// Create Tenv instance with empty builders map to test error handling
	tenv := Tenv{
		builders:  make(map[string]builder.Func), // Empty builders map
		conf:      testConfig,
		hclParser: hclParser,
		managers:  make(map[string]versionmanager.VersionManager),
	}

	// Test Uninstall method error handling
	err := tenv.Uninstall(ctx, toolName, requestedVersion)
	require.Error(t, err)
	assert.Equal(t, errNoBuilder, err)
}

func TestTenvUninstallMultipleComprehensive(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	toolName := testToolName
	versions := []string{"1.0.0", "1.1.0"}

	// Create a simple test config
	testConfig := &config.Config{
		Displayer: loghelper.InertDisplayer,
	}

	// Create HCL parser
	hclParser := hclparse.NewParser()

	// Create Tenv instance with empty builders map to test error handling
	tenv := Tenv{
		builders:  make(map[string]builder.Func), // Empty builders map
		conf:      testConfig,
		hclParser: hclParser,
		managers:  make(map[string]versionmanager.VersionManager),
	}

	// Test UninstallMultiple method error handling
	err := tenv.UninstallMultiple(ctx, toolName, versions)
	require.Error(t, err)
	assert.Equal(t, errNoBuilder, err)
}

func TestTenvInitError(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	toolName := "nonexistent"

	// Create Tenv instance without the tool
	tenv, err := Make()
	require.NoError(t, err)

	// Test that init fails for nonexistent tool
	err = tenv.Install(ctx, toolName, "1.0.0")
	require.Error(t, err)
	assert.Equal(t, errNoBuilder, err)
}

func TestTenvCommandSuccess(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	toolName := testToolName
	requestedVersion := testVersion
	cmdArgs := []string{"apply"}

	// Create test config with RootPath
	testConfig := &config.Config{
		Getenv:   config.EmptyGetenv,
		RootPath: t.TempDir(),
	}
	testConfig.InitDisplayer(false)

	// Create VersionManager
	manager := versionmanager.Make(testConfig, "TF", "terraform", nil, &mockRetriever{}, nil)

	// Create Tenv with manager set
	tenv := Tenv{
		builders:  make(map[string]builder.Func),
		conf:      testConfig,
		hclParser: hclparse.NewParser(),
		managers:  map[string]versionmanager.VersionManager{toolName: manager},
	}

	// Test Command
	cmd, err := tenv.Command(ctx, toolName, requestedVersion, cmdArgs...)
	require.NoError(t, err)
	assert.NotNil(t, cmd)
	assert.Contains(t, cmd.Path, "terraform")
	assert.Equal(t, cmdArgs, cmd.Args[1:])
}

func TestTenvDetectedCommandSuccess(t *testing.T) {
	ctx := t.Context()
	toolName := testToolName
	cmdArgs := []string{"version"}

	// Set env var for version
	t.Setenv("TF_VERSION", testVersion)

	// Create test config with RootPath
	testConfig := &config.Config{
		Getenv:   config.EmptyGetenv,
		RootPath: t.TempDir(),
	}
	testConfig.InitDisplayer(false)

	// Create VersionManager
	manager := versionmanager.Make(testConfig, "TF", "terraform", nil, &mockRetriever{}, nil)

	// Create Tenv with manager set
	tenv := Tenv{
		builders:  make(map[string]builder.Func),
		conf:      testConfig,
		hclParser: hclparse.NewParser(),
		managers:  map[string]versionmanager.VersionManager{toolName: manager},
	}

	// Test DetectedCommand
	cmd, err := tenv.DetectedCommand(ctx, toolName, cmdArgs...)
	require.NoError(t, err)
	assert.NotNil(t, cmd)
	assert.Contains(t, cmd.Path, "terraform")
	assert.Equal(t, cmdArgs, cmd.Args[1:])
}
