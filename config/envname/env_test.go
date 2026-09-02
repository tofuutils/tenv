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

package envname

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type environmentConstantTest struct {
	name     string
	constant string
	expected string
}

func runEnvironmentConstantsTest(t *testing.T, tests []environmentConstantTest) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, tt.constant)
		})
	}
}

func baseEnvironmentTests() []environmentConstantTest {
	// This is the base environment constants
	consts := []environmentConstantTest{
		{
			name:     "agnosticProxy constant",
			constant: agnosticProxy,
			expected: "AGNOSTIC_PROXY",
		},
		{
			name:     "arch constant",
			constant: arch,
			expected: "ARCH",
		},
		{
			name:     "autoInstall constant",
			constant: autoInstall,
			expected: "AUTO_INSTALL",
		},
		{
			name:     "DefaultConstraintSuffix constant",
			constant: DefaultConstraintSuffix,
			expected: "DEFAULT_CONSTRAINT",
		},
		{
			name:     "DefaultVersionSuffix constant",
			constant: DefaultVersionSuffix,
			expected: "DEFAULT_VERSION",
		},
		{
			name:     "forceRemote constant",
			constant: forceRemote,
			expected: "FORCE_REMOTE",
		},
		{
			name:     "installMode constant",
			constant: installMode,
			expected: "INSTALL_MODE",
		},
		{
			name:     "listMode constant",
			constant: listMode,
			expected: "LIST_MODE",
		},
		{
			name:     "listURL constant",
			constant: listURL,
			expected: "LIST_URL",
		},
		{
			name:     "log constant",
			constant: log,
			expected: "LOG",
		},
		{
			name:     "quiet constant",
			constant: quiet,
			expected: "QUIET",
		},
		{
			name:     "remotePass constant",
			constant: remotePass,
			expected: "REMOTE_PASSWORD",
		},
		{
			name:     "remoteURL constant",
			constant: remoteURL,
			expected: "REMOTE",
		},
		{
			name:     "remoteUser constant",
			constant: remoteUser,
			expected: "REMOTE_USER",
		},
		{
			name:     "rootPath constant",
			constant: rootPath,
			expected: "ROOT",
		},
		{
			name:     "VersionSuffix constant",
			constant: VersionSuffix,
			expected: "VERSION",
		},
	}

	return consts
}

func TestBaseEnvironmentConstants(t *testing.T) {
	t.Parallel()
	runEnvironmentConstantsTest(t, baseEnvironmentTests())
}

func githubEnvironmentTests() []environmentConstantTest {
	return []environmentConstantTest{
		{
			name:     "githubPrefix constant",
			constant: githubPrefix,
			expected: "GITHUB_",
		},
		{
			name:     "GithubActions constant",
			constant: GithubActions,
			expected: "GITHUB_ACTIONS",
		},
		{
			name:     "GithubOutput constant",
			constant: GithubOutput,
			expected: "GITHUB_OUTPUT",
		},
		{
			name:     "token constant",
			constant: token,
			expected: "GITHUB_TOKEN",
		},
	}
}

func TestGitHubEnvironmentConstants(t *testing.T) {
	t.Parallel()
	runEnvironmentConstantsTest(t, githubEnvironmentTests())
}

func atmosEnvironmentTests() []environmentConstantTest {
	return []environmentConstantTest{
		{
			name:     "AtmosPrefix constant",
			constant: AtmosPrefix,
			expected: "ATMOS_",
		},
		{
			name:     "AtmosInstallMode constant",
			constant: AtmosInstallMode,
			expected: "ATMOS_INSTALL_MODE",
		},
		{
			name:     "AtmosListMode constant",
			constant: AtmosListMode,
			expected: "ATMOS_LIST_MODE",
		},
		{
			name:     "AtmosListURL constant",
			constant: AtmosListURL,
			expected: "ATMOS_LIST_URL",
		},
		{
			name:     "AtmosRemotePass constant",
			constant: AtmosRemotePass,
			expected: "ATMOS_REMOTE_PASSWORD",
		},
		{
			name:     "AtmosRemoteURL constant",
			constant: AtmosRemoteURL,
			expected: "ATMOS_REMOTE",
		},
		{
			name:     "AtmosRemoteUser constant",
			constant: AtmosRemoteUser,
			expected: "ATMOS_REMOTE_USER",
		},
	}
}

func TestAtmosEnvironmentConstants(t *testing.T) {
	t.Parallel()
	runEnvironmentConstantsTest(t, atmosEnvironmentTests())
}

func tenvEnvironmentTests() []environmentConstantTest {
	return []environmentConstantTest{
		{
			name:     "tenvPrefix constant",
			constant: tenvPrefix,
			expected: "TENV_",
		},
		{
			name:     "TenvArch constant",
			constant: TenvArch,
			expected: "TENV_ARCH",
		},
		{
			name:     "TenvAutoInstall constant",
			constant: TenvAutoInstall,
			expected: "TENV_AUTO_INSTALL",
		},
		{
			name:     "TenvForceRemote constant",
			constant: TenvForceRemote,
			expected: "TENV_FORCE_REMOTE",
		},
		{
			name:     "TenvLog constant",
			constant: TenvLog,
			expected: "TENV_LOG",
		},
		{
			name:     "TenvQuiet constant",
			constant: TenvQuiet,
			expected: "TENV_QUIET",
		},
		{
			name:     "TenvRemoteConf constant",
			constant: TenvRemoteConf,
			expected: "TENV_REMOTE_CONF",
		},
		{
			name:     "TenvRootPath constant",
			constant: TenvRootPath,
			expected: "TENV_ROOT",
		},
		{
			name:     "TenvSkipLastUse constant",
			constant: TenvSkipLastUse,
			expected: "TENV_SKIP_LAST_USE",
		},
		{
			name:     "TenvToken constant",
			constant: TenvToken,
			expected: "TENV_GITHUB_TOKEN",
		},
		{
			name:     "TenvValidation constant",
			constant: TenvValidation,
			expected: "TENV_VALIDATION",
		},
	}
}

func TestTenvEnvironmentConstants(t *testing.T) {
	t.Parallel()
	runEnvironmentConstantsTest(t, tenvEnvironmentTests())
}

func terraformEnvironmentTests() []environmentConstantTest {
	return []environmentConstantTest{
		{
			name:     "TfenvPrefix constant",
			constant: TfenvPrefix,
			expected: "TFENV_",
		},
		{
			name:     "TfenvTerraformPrefix constant",
			constant: TfenvTerraformPrefix,
			expected: "TFENV_TERRAFORM_",
		},
		{
			name:     "TfAgnostic constant",
			constant: TfAgnostic,
			expected: "TFENV_AGNOSTIC_PROXY",
		},
		{
			name:     "TfArch constant",
			constant: TfArch,
			expected: "TFENV_ARCH",
		},
		{
			name:     "TfAutoInstall constant",
			constant: TfAutoInstall,
			expected: "TFENV_AUTO_INSTALL",
		},
		{
			name:     "TfForceRemote constant",
			constant: TfForceRemote,
			expected: "TFENV_FORCE_REMOTE",
		},
		{
			name:     "TfHashicorpPGPKey constant",
			constant: TfHashicorpPGPKey,
			expected: "TFENV_HASHICORP_PGP_KEY",
		},
		{
			name:     "TfInstallMode constant",
			constant: TfInstallMode,
			expected: "TFENV_INSTALL_MODE",
		},
		{
			name:     "TfListMode constant",
			constant: TfListMode,
			expected: "TFENV_LIST_MODE",
		},
		{
			name:     "TfListURL constant",
			constant: TfListURL,
			expected: "TFENV_LIST_URL",
		},
		{
			name:     "TfRemotePass constant",
			constant: TfRemotePass,
			expected: "TFENV_REMOTE_PASSWORD",
		},
		{
			name:     "TfRemoteURL constant",
			constant: TfRemoteURL,
			expected: "TFENV_REMOTE",
		},
		{
			name:     "TfRemoteUser constant",
			constant: TfRemoteUser,
			expected: "TFENV_REMOTE_USER",
		},
		{
			name:     "TfRootPath constant",
			constant: TfRootPath,
			expected: "TFENV_ROOT",
		},
	}
}

func TestTerraformEnvironmentConstants(t *testing.T) {
	t.Parallel()
	runEnvironmentConstantsTest(t, terraformEnvironmentTests())
}

func terragruntEnvironmentTests() []environmentConstantTest {
	return []environmentConstantTest{
		{
			name:     "TgPrefix constant",
			constant: TgPrefix,
			expected: "TG_",
		},
		{
			name:     "TgInstallMode constant",
			constant: TgInstallMode,
			expected: "TG_INSTALL_MODE",
		},
		{
			name:     "TgListMode constant",
			constant: TgListMode,
			expected: "TG_LIST_MODE",
		},
		{
			name:     "TgListURL constant",
			constant: TgListURL,
			expected: "TG_LIST_URL",
		},
		{
			name:     "TgRemotePass constant",
			constant: TgRemotePass,
			expected: "TG_REMOTE_PASSWORD",
		},
		{
			name:     "TgRemoteURL constant",
			constant: TgRemoteURL,
			expected: "TG_REMOTE",
		},
		{
			name:     "TgRemoteUser constant",
			constant: TgRemoteUser,
			expected: "TG_REMOTE_USER",
		},
	}
}

func TestTerragruntEnvironmentConstants(t *testing.T) {
	t.Parallel()
	runEnvironmentConstantsTest(t, terragruntEnvironmentTests())
}

func terramateEnvironmentTests() []environmentConstantTest {
	return []environmentConstantTest{
		{
			name:     "TmPrefix constant",
			constant: TmPrefix,
			expected: "TM_",
		},
		{
			name:     "TmInstallMode constant",
			constant: TmInstallMode,
			expected: "TM_INSTALL_MODE",
		},
		{
			name:     "TmListMode constant",
			constant: TmListMode,
			expected: "TM_LIST_MODE",
		},
		{
			name:     "TmListURL constant",
			constant: TmListURL,
			expected: "TM_LIST_URL",
		},
		{
			name:     "TmRemotePass constant",
			constant: TmRemotePass,
			expected: "TM_REMOTE_PASSWORD",
		},
		{
			name:     "TmRemoteURL constant",
			constant: TmRemoteURL,
			expected: "TM_REMOTE",
		},
		{
			name:     "TmRemoteUser constant",
			constant: TmRemoteUser,
			expected: "TM_REMOTE_USER",
		},
	}
}

func TestTerramateEnvironmentConstants(t *testing.T) {
	t.Parallel()
	runEnvironmentConstantsTest(t, terramateEnvironmentTests())
}

func tofuEnvironmentTests() []environmentConstantTest {
	var tests []environmentConstantTest
	tests = append(tests, environmentConstantTest{
		name:     "TofuenvPrefix constant",
		constant: TofuenvPrefix,
		expected: "TOFUENV_",
	})
	tests = append(tests, environmentConstantTest{
		name:     "TofuenvTofuPrefix constant",
		constant: TofuenvTofuPrefix,
		expected: "TOFUENV_TOFU_",
	})
	tests = append(tests, environmentConstantTest{
		name:     "TofuAgnostic constant",
		constant: TofuAgnostic,
		expected: "TOFUENV_AGNOSTIC_PROXY",
	})
	tests = append(tests, environmentConstantTest{
		name:     "TofuArch constant",
		constant: TofuArch,
		expected: "TOFUENV_ARCH",
	})
	tests = append(tests, environmentConstantTest{
		name:     "TofuAutoInstall constant",
		constant: TofuAutoInstall,
		expected: "TOFUENV_AUTO_INSTALL",
	})
	tests = append(tests, environmentConstantTest{
		name:     "TofuForceRemote constant",
		constant: TofuForceRemote,
		expected: "TOFUENV_FORCE_REMOTE",
	})
	tests = append(tests, environmentConstantTest{
		name:     "TofuInstallMode constant",
		constant: TofuInstallMode,
		expected: "TOFUENV_INSTALL_MODE",
	})
	tests = append(tests, environmentConstantTest{
		name:     "TofuListMode constant",
		constant: TofuListMode,
		expected: "TOFUENV_LIST_MODE",
	})
	tests = append(tests, environmentConstantTest{
		name:     "TofuListURL constant",
		constant: TofuListURL,
		expected: "TOFUENV_LIST_URL",
	})
	tests = append(tests, environmentConstantTest{
		name:     "TofuOpenTofuPGPKey constant",
		constant: TofuOpenTofuPGPKey,
		expected: "TOFUENV_OPENTOFU_PGP_KEY",
	})
	tests = append(tests, environmentConstantTest{
		name:     "TofuRemotePass constant",
		constant: TofuRemotePass,
		expected: "TOFUENV_REMOTE_PASSWORD",
	})
	tests = append(tests, environmentConstantTest{
		name:     "TofuRemoteURL constant",
		constant: TofuRemoteURL,
		expected: "TOFUENV_REMOTE",
	})
	tests = append(tests, environmentConstantTest{
		name:     "TofuRemoteUser constant",
		constant: TofuRemoteUser,
		expected: "TOFUENV_REMOTE_USER",
	})
	tests = append(tests, environmentConstantTest{
		name:     "TofuRootPath constant",
		constant: TofuRootPath,
		expected: "TOFUENV_ROOT",
	})
	tests = append(tests, environmentConstantTest{
		name:     "TofuToken constant",
		constant: TofuToken,
		expected: "TOFUENV_GITHUB_TOKEN",
	})
	tests = append(tests, environmentConstantTest{
		name:     "TofuURLTemplate constant",
		constant: TofuURLTemplate,
		expected: "TOFUENV_URL_TEMPLATE",
	})

	return tests
}

func TestTofuEnvironmentConstants(t *testing.T) {
	t.Parallel()
	runEnvironmentConstantsTest(t, tofuEnvironmentTests())
}
