package text

const RootLongText = `A compact CLI that manages OpenTofu / Terraform version via tfenv/tofuenv wrappers

Authors:  Alexander Sharov (kvendingoldo@gmail.com), Nikolai Mishin (sanduku.default@gmail.com), Anastasiia Kozlova (anastasiia.kozlova245@gmail.com)
Contributed at https://github.com/tfutils/tfenv
`
const EmptyArgsText = `Error: please use --help|-h to explore tenv commands`

const AdditionalText = `A compact CLI that manages OpenTofu / Terraform version via tfenv/tofuenv wrappers

Authors:  Alexander Sharov (kvendingoldo@gmail.com), Nikolai Mishin (sanduku.default@gmail.com(, Anastasiia Kozlova (anastasiia.kozlova245@gmail.com)
Contributed at https://github.com/tfutils/tfenv

Usages:
	tenv TODO
	tenv TODO

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
const UninstallDepsCmdLongText = `This command uninstall all tenv dependencies
Usages:
	tenv uninstallDeps
`
const UpgradeDepsCmdLongText = `This command upgrade all tenv dependencies
Usages:
	tenv upgradeDeps
`
