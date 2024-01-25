<details><summary><b>TFENV_AUTO_INSTALL (alias TFENV_AUTO_INSTALL)</b></summary><br>
String (Default: true)

If set to true tenv will automatically install a missing Terraform version needed (fallback to latest-allowed strategy when no [`.terraform-version`](#terraform-version-file) files are found).

`tenv` subcommands `detect` and `use` support a `--no-install`, `-n` disabling flag version.

Example: Use Terraform version 1.6.0-rc1 that is not installed, and auto installation is disabled. (-v flag is equivalent to `TFENV_VERBOSE=true`)

```console
$ TFENV_AUTO_INSTALL=false tenv use -v 1.6.0-rc1
Write 1.6.0-rc1 in /home/dvaumoron/.tenv/.terraform-version
```

Example: Use Terraform version 1.6.0-rc1 that is not installed, and auto installation stay enabled.

```console
$ tenv use -v 1.6.0-rc1
Installation of Terraform 1.6.0-rc1
Write 1.6.0-rc1 in /home/dvaumoron/.tenv/.terraform-version
```
</details>


<details><summary><b>TFENV_FORCE_REMOTE (alias TFENV_FORCE_REMOTE)</b></summary><br>
String (Default: false)

If set to true tenv detection of needed version will skip local check and verify compatibility on remote list.

`tenv` subcommands `detect` and `use` support a `--force-remote`, `-f` flag version.
</details>


<details><summary><b>TFENV_GITHUB_TOKEN</b></summary><br>
String (Default: "")

Allow to specify a GitHub token to increase [GitHub Rate limits for the REST API](https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api). Useful because Terraform binares are downloaded from the Terraform GitHub repository.

`tenv` subcommands `detect`, `install`, `list-remote` and `use` support a `--github-token`, `-t` flag version.
</details>


<details><summary><b>TFENV_terraform_PGP_KEY</b></summary><br>
String (Default: "")

Allow to specify a local file path to Terraform PGP public key, if not present download https://get.terraform.org/terraform.asc.

`tenv` subcommands `detect`, `ìnstall` and `use` support a `--key-file`, `-k` flag version.
</details>


<details><summary><b>TFENV_HASHICORP_PGP_KEY</b></summary><br>
String (Default: "")

Allow to specify a local file path to Hashicorp PGP public key, if not present download https://www.hashicorp.com/.well-known/pgp-key.txt.

`tenv tf` subcommands `detect`, `ìnstall` and `use` support a `--key-file`, `-k` flag version.
</details>


<details><summary><b>TFENV_REMOTE</b></summary><br>
String (Default: https://api.github.com/repos/terraform/terraform/releases)

To install Terraform from a remote other than the default (must comply with [Github REST API](https://docs.github.com/en/rest?apiVersion=2022-11-28))

`tenv` subcommands `detect`, `install`, `list-remote` and `use` support a `--remote-url`, `-u` flag version.
</details>


<details><summary><b>TFENV_REMOTE</b></summary><br>
String (Default: https://releases.hashicorp.com/terraform)

To install Terraform from a remote other than the default (must comply with [Hashicorp Release API](https://releases.hashicorp.com/docs/api/v1))

`tenv tf` subcommands `detect`, `install`, `list-remote` and `use` support a `--remote-url`, `-u` flag version.
</details>


<details><summary><b>TFENV_ROOT (alias TFENV_ROOT)</b></summary><br>
Path (Default: `$HOME/.tenv`)

The path to a directory where the local Terraform versions, Terraform versions and tenv configuration files exist.

`tenv` support a `--root-path`, `-r` flag version.
</details>


<details><summary><b>TFENV_TOFU_VERSION</b></summary><br>
String (Default: "")

If not empty string, this variable overrides Terraform version, specified in [`.terraform-version`](#terraform-version-file) files.
`tenv` subcommands `install` and `detect` also respects this variable.

e.g. with :

```console
$ tofu version
Terraform v1.6.0
on linux_amd64
```

then :

```console
$ TFENV_TOFU_VERSION=1.6.0-rc1 tofu version
Terraform v1.6.0-rc1
on linux_amd64
```
</details>


<details><summary><b>TFENV_TERRAFORM_VERSION</b></summary><br>
String (Default: "")

If not empty string, this variable overrides Terraform version, specified in `.terraform-version` files.

`tenv tf` subcommands `install` and `detect` also respects this variable.
</details>


<details><summary><b>TFENV_VERBOSE</b></summary><br>
String (Default: false)

Active the verbose display of tenv.

`tenv` support a `--verbose`, `-v` flag version.
</details>