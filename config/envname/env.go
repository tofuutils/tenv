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

const (
	arch                    = "ARCH"
	autoInstall             = "AUTO_INSTALL"
	DefaultConstraintSuffix = "DEFAULT_CONSTRAINT"
	DefaultVersionSuffix    = "DEFAULT_" + VersionSuffix
	forceRemote             = "FORCE_REMOTE"
	installMode             = "INSTALL_MODE"
	listMode                = "LIST_MODE"
	listURL                 = "LIST_URL"
	log                     = "LOG"
	quiet                   = "QUIET"
	remotePass              = "REMOTE_PASSWORD"
	remoteURL               = "REMOTE"
	remoteUser              = "REMOTE_USER"
	rootPath                = "ROOT"
	VersionSuffix           = "VERSION"

	githubPrefix  = "GITHUB_"
	GithubActions = githubPrefix + "ACTIONS"
	GithubOutput  = githubPrefix + "OUTPUT"
	token         = githubPrefix + "TOKEN"

	AtmosPrefix      = "ATMOS_"
	AtmosInstallMode = AtmosPrefix + installMode
	AtmosListMode    = AtmosPrefix + listMode
	AtmosListURL     = AtmosPrefix + listURL
	AtmosRemotePass  = AtmosPrefix + remotePass
	AtmosRemoteURL   = AtmosPrefix + remoteURL
	AtmosRemoteUser  = AtmosPrefix + remoteUser

	tenvPrefix      = "TENV_"
	TenvArch        = tenvPrefix + arch
	TenvAutoInstall = tenvPrefix + autoInstall
	TenvForceRemote = tenvPrefix + forceRemote
	TenvLog         = tenvPrefix + log
	TenvQuiet       = tenvPrefix + quiet
	TenvRemoteConf  = tenvPrefix + "REMOTE_CONF"
	TenvRootPath    = tenvPrefix + rootPath
	TenvToken       = tenvPrefix + token

	TfenvPrefix          = "TFENV_"
	TfenvTerraformPrefix = TfenvPrefix + "TERRAFORM_"
	TfArch               = TfenvPrefix + arch
	TfAutoInstall        = TfenvPrefix + autoInstall
	TfForceRemote        = TfenvPrefix + forceRemote
	TfHashicorpPGPKey    = TfenvPrefix + "HASHICORP_PGP_KEY"
	TfInstallMode        = TfenvPrefix + installMode
	TfListMode           = TfenvPrefix + listMode
	TfListURL            = TfenvPrefix + listURL
	TfRemotePass         = TfenvPrefix + remotePass
	TfRemoteURL          = TfenvPrefix + remoteURL
	TfRemoteUser         = TfenvPrefix + remoteUser
	TfRootPath           = TfenvPrefix + rootPath

	TgPrefix      = "TG_"
	TgInstallMode = TgPrefix + installMode
	TgListMode    = TgPrefix + listMode
	TgListURL     = TgPrefix + listURL
	TgRemotePass  = TgPrefix + remotePass
	TgRemoteURL   = TgPrefix + remoteURL
	TgRemoteUser  = TgPrefix + remoteUser

	TofuenvPrefix      = "TOFUENV_"
	TofuenvTofuPrefix  = TofuenvPrefix + "TOFU_"
	TofuArch           = TofuenvPrefix + arch
	TofuAutoInstall    = TofuenvPrefix + autoInstall
	TofuForceRemote    = TofuenvPrefix + forceRemote
	TofuInstallMode    = TofuenvPrefix + installMode
	TofuListMode       = TofuenvPrefix + listMode
	TofuListURL        = TofuenvPrefix + listURL
	TofuOpenTofuPGPKey = TofuenvPrefix + "OPENTOFU_PGP_KEY"
	TofuRemotePass     = TofuenvPrefix + remotePass
	TofuRemoteURL      = TofuenvPrefix + remoteURL
	TofuRemoteUser     = TofuenvPrefix + remoteUser
	TofuRootPath       = TofuenvPrefix + rootPath
	TofuToken          = TofuenvPrefix + token
	TofuURLTemplate    = TofuenvPrefix + "URL_TEMPLATE"
)
