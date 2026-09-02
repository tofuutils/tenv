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

package terragruntparser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
	terragruntparser "github.com/tofuutils/tenv/v4/versionmanager/semantic/parser/terragrunt"
)

func testConf() *config.Config {
	return &config.Config{Displayer: loghelper.InertDisplayer}
}

func TestRetrieveTerragruntVersionConstraintFromHCL(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, terragruntparser.HCLName)
	require.NoError(t, os.WriteFile(filePath, []byte(`terragrunt_version_constraint = "0.69.1"`), 0o600))

	p := terragruntparser.Make(hclparse.NewParser())
	version, err := p.RetrieveTerragruntVersionConstraintFromHCL(filePath, testConf())

	require.NoError(t, err)
	assert.Equal(t, "0.69.1", version)
}

func TestRetrieveTerraformVersionConstraintFromHCL(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, terragruntparser.HCLNameLegacy)
	require.NoError(t, os.WriteFile(filePath, []byte(`terraform_version_constraint = ">= 1.5.0"`), 0o600))

	p := terragruntparser.Make(hclparse.NewParser())
	version, err := p.RetrieveTerraformVersionConstraintFromHCL(filePath, testConf())

	require.NoError(t, err)
	assert.Equal(t, ">= 1.5.0", version)
}

func TestRetrieveTerragruntVersionConstraintFromJSON(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, terragruntparser.JSONName)
	require.NoError(t, os.WriteFile(filePath, []byte(`{"terragrunt_version_constraint": "1.2.3"}`), 0o600))

	p := terragruntparser.Make(hclparse.NewParser())
	version, err := p.RetrieveTerragruntVersionConstraintFromJSON(filePath, testConf())

	require.NoError(t, err)
	assert.Equal(t, "1.2.3", version)
}

func TestRetrieveTerraformVersionConstraintFromJSON(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, terragruntparser.JSONNameLegacy)
	require.NoError(t, os.WriteFile(filePath, []byte(`{"terraform_version_constraint": "1.6.0"}`), 0o600))

	p := terragruntparser.Make(hclparse.NewParser())
	version, err := p.RetrieveTerraformVersionConstraintFromJSON(filePath, testConf())

	require.NoError(t, err)
	assert.Equal(t, "1.6.0", version)
}

func TestRetrieveVersionConstraint_MissingFile(t *testing.T) {
	t.Parallel()

	filePath := filepath.Join(t.TempDir(), "does-not-exist.hcl")

	p := terragruntparser.Make(hclparse.NewParser())
	version, err := p.RetrieveTerragruntVersionConstraintFromHCL(filePath, testConf())

	require.NoError(t, err, "a missing file must be treated as absent constraint, not an error")
	assert.Empty(t, version)
}

func TestRetrieveVersionConstraint_AttributeAbsent(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, terragruntparser.HCLName)
	require.NoError(t, os.WriteFile(filePath, []byte(`some_other_attribute = "value"`), 0o600))

	p := terragruntparser.Make(hclparse.NewParser())
	version, err := p.RetrieveTerragruntVersionConstraintFromHCL(filePath, testConf())

	require.NoError(t, err)
	assert.Empty(t, version)
}

func TestRetrieveVersionConstraint_InvalidHCLSyntax(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, terragruntparser.HCLName)
	require.NoError(t, os.WriteFile(filePath, []byte(`terragrunt_version_constraint = `), 0o600))

	p := terragruntparser.Make(hclparse.NewParser())
	version, err := p.RetrieveTerragruntVersionConstraintFromHCL(filePath, testConf())

	require.Error(t, err)
	assert.Empty(t, version)
}

func TestRetrieveVersionConstraint_NullValue(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, terragruntparser.HCLName)
	require.NoError(t, os.WriteFile(filePath, []byte(`terragrunt_version_constraint = null`), 0o600))

	p := terragruntparser.Make(hclparse.NewParser())
	version, err := p.RetrieveTerragruntVersionConstraintFromHCL(filePath, testConf())

	require.NoError(t, err)
	assert.Empty(t, version)
}

func TestRetrieveVersionConstraint_NonConvertibleType(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, terragruntparser.HCLName)
	require.NoError(t, os.WriteFile(filePath, []byte(`terragrunt_version_constraint = ["not", "a", "string"]`), 0o600))

	p := terragruntparser.Make(hclparse.NewParser())
	version, err := p.RetrieveTerragruntVersionConstraintFromHCL(filePath, testConf())

	require.NoError(t, err, "conversion failures are logged and treated as absent, not returned as an error")
	assert.Empty(t, version)
}
