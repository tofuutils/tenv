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
package text

const RootLongText = `A compact CLI that manages OpenTofu / Terraform version via tfenv/tofuenv wrappers

Authors:  Alexander Sharov (kvendingoldo@gmail.com), Nikolai Mishin (sanduku.default@gmail.com), Anastasiia Kozlova (anastasiia.kozlova245@gmail.com)
Contributed at https://github.com/tofuutils/tenv
`
const EmptyArgsText = `Error: please use --help|-h to explore tenv commands`

const AdditionalText = `A compact CLI that manages OpenTofu / Terraform version via tfenv/tofuenv wrappers

Authors:  Alexander Sharov (kvendingoldo@gmail.com), Nikolai Mishin (sanduku.default@gmail.com), Anastasiia Kozlova (anastasiia.kozlova245@gmail.com)
Contributed at https://github.com/tfutils/tfenv

Usages:
	tenv tf <tfenv command>
	tenv init

Flags:
	--help	display command help & instructions`

const SubCmdHelpText = `

Explore tenv commands at https://github.com/tfutils/tfenv`

const InitCmdLongText = `This command let you init tenv
Usages:
	tenv init
`

const TfCmdLongText = `This command let you manage Terraform version via tfenv wrapper
Usages:
	tenv tf list-remote
`

const TofuCmdLongText = `This command let you manage OpenTofu version via tofuenv wrapper
Usages:
	tenv tofu list-remote
`
const UninstallDepsCmdLongText = `This command uninstall all tenv dependencies (tfenv and tofuenv)
Usages:
	tenv uninstallDeps
`
const UpgradeDepsCmdLongText = `This command upgrade all tenv dependencies: download tfenv, tofuenv and configure their environment variables
Usages:
	tenv upgradeDeps
`
