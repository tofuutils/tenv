<!-- BADGES -->
[![Github release](https://img.shields.io/github/v/release/tofuutils/tenv?style=for-the-badge)](https://github.com/tofuutils/tenv/releases) [![Contributors](https://img.shields.io/github/contributors/tofuutils/tenv?style=for-the-badge)](https://github.com/tofuutils/tenv/graphs/contributors) ![maintenance status](https://img.shields.io/maintenance/yes/2024.svg?style=for-the-badge)


<!-- LOGO -->
<br />
<div align="center">
  <a>
    <img src="assets/logo.png" alt="Logo" width="200" height="200">
  </a>

<h3 align="center">tenv</h3>

  <p align="center">
    Terraform and OpenTofu version manager, written in Go.
    <br />
    ·
    <a href="https://github.com/othneildrew/Best-README-Template/issues">Report Bug</a>
    ·
    <a href="https://github.com/othneildrew/Best-README-Template/issues">Request Feature</a>
  </p>
</div>


## About The Project

Welcome to tenv, a versatile version manager for [OpenTofu](https://opentofu.org) and [Terraform](https://www.terraform.io/), written in Go. Our tool simplifies the complexity of handling different versions of these powerful tools, ensuring developers and DevOps professionals can focus on what matters most - building and deploying efficiently.
tenv is inspired by [tofuenv](https://github.com/tofuutils/tofuenv) and tfenv.

### Key Features

- Versatile Version Management: Easily switch between different versions of Terraform and OpenTofu.
- [Semver 2.0.0](https://semver.org/) Compatibility: Utilizes [go-version](https://github.com/hashicorp/go-version) for semantic versioning and use the [HCL](https://github.com/hashicorp/hcl) parser to extract required version constraint from OpenTofu/Terraform files.
- Signature Verification: Supports [cosign](https://github.com/sigstore/cosign) (if present on your machine) and PGP (via [gopenpgp](https://github.com/ProtonMail/gopenpgp)) for verifying OpenTofu signatures. However, unstable OpenTofu versions are signed only with cosign (in this case, if cosign is not found tenv will display a warning).
- Intuitive Installation: Simple installation process with Homebrew and manual options.


## Table of Contents
<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#key-features">Key Features</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>


## Getting Started

### Prerequisites
If you need to enable cosign checks, install `cosign` via the following command:

```sh
brew install cosign
  ```

### Installation

#### Automatic Installation
<details><summary><b>MacOS (Homebrew)</b></summary><br>

```console
brew tap tofuutils/tap
brew install tenv
```
</details>

<details><summary><b>Ubuntu</b></summary><br>
TODO
</details>


#### Manual Installation
Get the most recent packaged binaries in either .zip or .tar.gz format by visiting the [release page](https://github.com/tofuutils/tenv/releases). After downloading, unzip the folder and seamlessly integrate it into your system's PATH.

#### Docker Installation
You can use dockerized version of tenv via the following way:

```shell
docker run -d -e HOME=abc --name=tenv --entrypoint=/bin/sh tofuutils/tenv
```

## Usage

<details><summary><b>tenv (tool) install [version]</b></summary><br>

Install a requested version of OpenTofu (into TOFUENV_ROOT directory from TOFUENV_REMOTE url).

Without a parameter, the version to use is resolved automatically via TOFUENV_TOFU_VERSION or [`.opentofu-version`](#opentofu-version-file) files
(searched in the working directory, user home directory, and TOFUENV_ROOT directory).
Use "latest-stable" when none are found.

If a parameter is passed, available options include:

- an exact [Semver 2.0.0](https://semver.org/) version string to install
- a [version constraint](https://opentofu.org/docs/language/expressions/version-constraints) string (checked against versions available at TOFUENV_REMOTE url)
- `latest` or `latest-stable` (checked against versions available at TOFUENV_REMOTE url)
- `latest-allowed` or `min-required` to scan your OpenTofu files to detect which version is maximally allowed or minimally required.

See [required_version](https://opentofu.org/docs/language/settings#specifying-a-required-opentofu-version) docs.

```console
tenv install 1.6.0-beta5
tenv install "~> 1.6.0"
tenv install latest
tenv install latest-stable
tenv install latest-allowed
tenv install min-required
```
</details>




<details><summary><b>tenv (tool) use</b></summary><br>

Switch the default OpenTofu version to use (set in [`.opentofu-version`](#opentofu-version-file) file in TOFUENV_ROOT).

`tenv use` has a `--working-dir`, `-w` flag to write [`.opentofu-version`](#opentofu-version-file) file in working directory.

Available parameter options:

- an exact [Semver 2.0.0](https://semver.org/) version string to use
- a [version constraint](https://opentofu.org/docs/language/expressions/version-constraints) string (checked against available in TOFUENV_ROOT directory)
- latest or latest-stable (checked against available in TOFUENV_ROOT directory)
- latest-allowed or min-required to scan your OpenTofu files to detect which version is maximally allowed or minimally required.

See [required_version](https://opentofu.org/docs/language/settings#specifying-a-required-opentofu-version) docs.

```console
tenv use min-required
tenv use v1.6.0-beta5
tenv use latest
tenv use latest-allowed
```
</details>

<details><summary><b>tenv (tool) detect</b></summary><br>

Detect the used version of OpenTofu for the working directory.

```console
$ tenv detect
OpenTofu 1.6.0 will be run from this directory.
```
</details>

<details><summary><b>tenv (tool) reset</b></summary><br>
Reset used version of OpenTofu (remove .opentofu-version file from TOFUENV_ROOT).

```console
tenv reset
```
</details>


<details><summary><b>tenv (tool) uninstall [version]</b></summary><br>
Uninstall a specific version of OpenTofu (remove it from TOFUENV_ROOT directory without interpretation).

```console
tenv uninstall v1.6.0-alpha4
```
</details>

<details><summary><b>tenv (tool) list</b></summary><br>

List installed OpenTofu versions (located in TOFUENV_ROOT directory), sorted in ascending version order.

`tenv list` has a `--descending`, `-d` flag to sort in descending order.

```console
$ tenv list
  1.6.0-rc1 
* 1.6.0 (set by /home/dvaumoron/.tenv/.opentofu-version)
```
</details>


<details><summary><b>tenv (tool) list-remote</b></summary><br>
List installable OpenTofu versions (from TOFUENV_REMOTE url), sorted in ascending version order.

`tenv list-remote` has a `--descending`, `-d` flag to sort in descending order.

`tenv list-remote` has a `--stable`, `-s` flag to display only stable version.

```console
$ tenv list-remote
1.6.0-alpha1
1.6.0-alpha2
1.6.0-alpha3
1.6.0-alpha4
1.6.0-alpha5
1.6.0-beta1
1.6.0-beta2
1.6.0-beta3
1.6.0-beta4
1.6.0-beta5
1.6.0-rc1 (installed)
1.6.0 (installed)
```
</details>


<details><summary><b>tenv help [command]</b></summary><br>
Help about any command.

You can use `--help` `-h` flag instead.

```console
$ tenv help tf detect
Display Terraform current version.

Usage:
  tenv tf detect [flags]

Flags:
  -f, --force-remote        force search on versions available at TFENV_REMOTE url
  -h, --help                help for detect
  -k, --key-file string     local path to PGP public key file (replace check against remote one)
  -n, --no-install          disable installation of missing version
  -u, --remote-url string   remote url to install from (default "https://releases.hashicorp.com/terraform")

Global Flags:
  -r, --root-path string   local path to install versions of OpenTofu and Terraform (default "/home/dvaumoron/.tenv")
  -v, --verbose            verbose output
```

```console
$ tenv use -h
Switch the default OpenTofu version to use (set in .opentofu-version file in TOFUENV_ROOT)

Available parameter options:
- an exact Semver 2.0.0 version string to use
- a version constraint string (checked against versions available in TOFUENV_ROOT directory)
- `latest` or `latest-stable` (checked against versions available in TOFUENV_ROOT directory)
- `latest-allowed` or `min-required` to scan your OpenTofu files to detect which version is maximally allowed or minimally required.

Usage:
  tenv use version [flags]

Flags:
  -f, --force-remote          force search on versions available at TOFUENV_REMOTE url
  -t, --github-token string   GitHub token (increases GitHub REST API rate limits)
  -h, --help                  help for use
  -k, --key-file string       local path to PGP public key file (replace check against remote one)
  -n, --no-install            disable installation of missing version
  -u, --remote-url string     remote url to install from (default "https://api.github.com/repos/opentofu/opentofu/releases")
  -w, --working-dir           create .opentofu-version file in working directory

Global Flags:
  -r, --root-path string   local path to install versions of OpenTofu and Terraform (default "/home/dvaumoron/.tenv")
  -v, --verbose            verbose output
```
</details>


## Environment Variables

tenv commands support the following environment variables.

<details><summary><b>TOFUENV_AUTO_INSTALL (alias TFENV_AUTO_INSTALL)</b></summary><br>
String (Default: true)

If set to true tenv will automatically install a missing OpenTofu version needed (fallback to latest-allowed strategy when no [`.opentofu-version`](#opentofu-version-file) files are found).

`tenv` subcommands `detect` and `use` support a `--no-install`, `-n` disabling flag version.

Example: Use OpenTofu version 1.6.0-rc1 that is not installed, and auto installation is disabled. (-v flag is equivalent to `TOFUENV_VERBOSE=true`)

```console
$ TOFUENV_AUTO_INSTALL=false tenv use -v 1.6.0-rc1
Write 1.6.0-rc1 in /home/dvaumoron/.tenv/.opentofu-version
```

Example: Use OpenTofu version 1.6.0-rc1 that is not installed, and auto installation stay enabled.

```console
$ tenv use -v 1.6.0-rc1
Installation of OpenTofu 1.6.0-rc1
Write 1.6.0-rc1 in /home/dvaumoron/.tenv/.opentofu-version
```
</details>


<details><summary><b>TOFUENV_FORCE_REMOTE (alias TFENV_FORCE_REMOTE)</b></summary><br>
String (Default: false)

If set to true tenv detection of needed version will skip local check and verify compatibility on remote list.

`tenv` subcommands `detect` and `use` support a `--force-remote`, `-f` flag version.
</details>


<details><summary><b>TOFUENV_GITHUB_TOKEN</b></summary><br>
String (Default: "")

Allow to specify a GitHub token to increase [GitHub Rate limits for the REST API](https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api). Useful because OpenTofu binares are downloaded from the OpenTofu GitHub repository.

`tenv` subcommands `detect`, `install`, `list-remote` and `use` support a `--github-token`, `-t` flag version.
</details>


<details><summary><b>TOFUENV_OPENTOFU_PGP_KEY</b></summary><br>
String (Default: "")

Allow to specify a local file path to OpenTofu PGP public key, if not present download https://get.opentofu.org/opentofu.asc.

`tenv` subcommands `detect`, `ìnstall` and `use` support a `--key-file`, `-k` flag version.
</details>


<details><summary><b>TFENV_HASHICORP_PGP_KEY</b></summary><br>
String (Default: "")

Allow to specify a local file path to Hashicorp PGP public key, if not present download https://www.hashicorp.com/.well-known/pgp-key.txt.

`tenv tf` subcommands `detect`, `ìnstall` and `use` support a `--key-file`, `-k` flag version.
</details>


<details><summary><b>TOFUENV_REMOTE</b></summary><br>
String (Default: https://api.github.com/repos/opentofu/opentofu/releases)

To install OpenTofu from a remote other than the default (must comply with [Github REST API](https://docs.github.com/en/rest?apiVersion=2022-11-28))

`tenv` subcommands `detect`, `install`, `list-remote` and `use` support a `--remote-url`, `-u` flag version.
</details>


<details><summary><b>TFENV_REMOTE</b></summary><br>
String (Default: https://releases.hashicorp.com/terraform)

To install Terraform from a remote other than the default (must comply with [Hashicorp Release API](https://releases.hashicorp.com/docs/api/v1))

`tenv tf` subcommands `detect`, `install`, `list-remote` and `use` support a `--remote-url`, `-u` flag version.
</details>


<details><summary><b>TOFUENV_ROOT (alias TFENV_ROOT)</b></summary><br>
Path (Default: `$HOME/.tenv`)

The path to a directory where the local OpenTofu versions, Terraform versions and tenv configuration files exist.

`tenv` support a `--root-path`, `-r` flag version.
</details>


<details><summary><b>TOFUENV_TOFU_VERSION</b></summary><br>
String (Default: "")

If not empty string, this variable overrides OpenTofu version, specified in [`.opentofu-version`](#opentofu-version-file) files.
`tenv` subcommands `install` and `detect` also respects this variable.

e.g. with :

```console
$ tofu version
OpenTofu v1.6.0
on linux_amd64
```

then :

```console
$ TOFUENV_TOFU_VERSION=1.6.0-rc1 tofu version
OpenTofu v1.6.0-rc1
on linux_amd64
```
</details>


<details><summary><b>TFENV_TERRAFORM_VERSION</b></summary><br>
String (Default: "")

If not empty string, this variable overrides Terraform version, specified in `.terraform-version` files.

`tenv tf` subcommands `install` and `detect` also respects this variable.
</details>


<details><summary><b>TOFUENV_VERBOSE (alias TFENV_VERBOSE)</b></summary><br>
String (Default: false)

Active the verbose display of tenv.

`tenv` support a `--verbose`, `-v` flag version.
</details>


## .opentofu-version file

If you put a `.opentofu-version` file in the working directory, user home directory, or TOFUENV_ROOT directory, tenv detects it and uses the version written in it.
Note, that TOFUENV_TOFU_VERSION can be used to override version specified by `.opentofu-version` file.

Recognized values (same as `tenv use` command):

- an exact [Semver 2.0.0](https://semver.org/) version string to use
- a [version constraint](https://opentofu.org/docs/language/expressions/version-constraints) string (checked against versions available in TOFUENV_ROOT directory)
- `latest` or `latest-stable` (checked against versions available in TOFUENV_ROOT directory)
- `latest-allowed` or `min-required` to scan your OpenTofu files to detect which version is maximally allowed or minimally required.

See [required_version](https://opentofu.org/docs/language/settings#specifying-a-required-opentofu-version) docs.


## Technical details

### Project binaries

#### tofu
The `tofu` command in this project is a proxy to OpenTofu's `tofu` command  managed by `tenv`. The default resolution strategy is latest-allowed (without [TOFUENV_TOFU_VERSION](#tofuenv_tofu_version) environment variable or [`.opentofu-version`](#opentofu-version-file) file).

#### terraform
The `terraform` command in this project is a proxy to HashiCorp's `terraform` command managed by `tenv`. The default resolution strategy is latest-allowed (without [TFENV_TERRAFORM_VERSION](#tfenv_terraform_version) environment variable or `.terraform-version` file).


### Terraform support

tenv relies on `.terraform-version` files, [TFENV_HASHICORP_PGP_KEY](#tfenv_hashicorp_pgp_key), [TFENV_REMOTE](#tfenv_remote) and [TFENV_TERRAFORM_VERSION](#tfenv_terraform_version) specifically to manage Terraform versions.

`tenv tf` have the same managing subcommands for Terraform versions (`detect`, `install`, `list`, `list-remote`, `reset`, `uninstall` and `use`).

tenv checks the Terraform PGP signature (there is no cosign signature available).


## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Community
Have questions or suggestions? Reach out to us via:

* [GitHub Issues](LINK_TO_ISSUES)
* User/Developer Group: Join github community to get update of Harbor's news, features, releases, or to provide suggestion and feedback.
* Slack: Join tofuutils's community for discussion and ask questions: OpenTofu, channel: #tofuutils

## Authors
tenv is based on [tofuenv](https://github.com/tofuutils/tofuenv) and [](https://github.com/tofuutils/gotofuenv) projects and supported by tofuutils team with help from these awesome contributors:

<!-- markdownlint-disable no-inline-html -->
<a href="https://github.com/tofuutils/tenv/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=tofuutils/tenv" />
</a>


<a href="https://star-history.com/#tofuutils/tenv&Date">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=tofuutils/tenv&type=Date&theme=dark" />
    <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=tofuutils/tenv&type=Date" />
    <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=tofuutils/pre-commit-opentofu&type=Date" />
  </picture>
</a>

<!-- markdownlint-enable no-inline-html -->

## LICENSE
The tenv project is distributed under the Apache 2.0 license. See [LICENSE](LICENSE).