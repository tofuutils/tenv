<!-- BADGES -->
[![Github release](https://img.shields.io/github/v/release/tofuutils/tenv?style=for-the-badge)](https://github.com/tofuutils/tenv/releases) [![Contributors](https://img.shields.io/github/contributors/tofuutils/tenv?style=for-the-badge)](https://github.com/tofuutils/tenv/graphs/contributors) ![maintenance status](https://img.shields.io/maintenance/yes/2024.svg?style=for-the-badge) [![Go report](https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=for-the-badge)](https://goreportcard.com/report/github.com/tofuutils/tenv/) [![codecov](https://img.shields.io/codecov/c/github/tofuutils/tenv?token=BDU9X0BAZV&style=for-the-badge)](https://codecov.io/gh/tofuutils/tenv)


<!-- LOGO -->
<br />
<div align="center">
  <a>
    <img src="assets/logo.png" alt="Logo" width="200" height="200">
  </a>
<h3 align="center">tenv</h3>
  <p align="center">
    OpenTofu, Terraform and Terragrunt version manager, written in Go.
    <br />
    ·
    <a href="https://github.com/tofuutils/tenv/issues/new?assignees=&labels=issue%3A+bug&projects=&template=bug_report.md&title=">Report Bug</a>
    ·
    <a href="https://github.com/tofuutils/tenv/issues/new?assignees=&labels=&projects=&template=feature_request.md&title=">Request Feature</a>
  </p>
</div>

<a id="about-the-project"></a>
## About The Project

Welcome to **tenv**, a versatile version manager for [OpenTofu](https://opentofu.org), [Terraform](https://www.terraform.io/) and [Terragrunt](https://terragrunt.gruntwork.io/), written in Go. Our tool simplifies the complexity of handling different versions of these powerful tools, ensuring developers and DevOps professionals can focus on what matters most - building and deploying efficiently.

**tenv** is a successor of [tofuenv](https://github.com/tofuutils/tofuenv) and [tfenv](https://github.com/tfutils/tfenv).

<a id="key-features"></a>
### Key Features

- Versatile version management: Easily switch between different versions of OpenTofu, Terraform and Terragrunt.
- [Semver 2.0.0](https://semver.org/) Compatibility: Utilizes [go-version](https://github.com/hashicorp/go-version) for semantic versioning and use the [HCL](https://github.com/hashicorp/hcl) parser to extract required version constraint from OpenTofu/Terraform/Terragrunt files (see [required_version](#required_version) and [Terragrunt hcl](#terragrunt-hcl-file)).
- Signature verification: Supports [cosign](https://github.com/sigstore/cosign) (if present on your machine) and PGP (via [gopenpgp](https://github.com/ProtonMail/gopenpgp)), see [signature support](#signature-support).
- Intuitive installation: Simple installation process with Homebrew and manual options.

<a id="difference-with-asdf"></a>
### Difference with asdf

[asdf](https://asdf-vm.com/) share the same goals than **tenv** : simplify the usage of several version of tools.

asdf is generic and extensible with a plugin system, key **tenv** differences :

- **tenv** is more specific and has features dedicated to OpenTofu, Terraform and Terragrunt, like [HCL](https://github.com/hashicorp/hcl) parsing based detection (see [Key Features](#key-features)).
- **tenv** is distributed as independent binaries and does not rely on any shell or other CLI executable.

<a id="table-of-contents"></a>
## Table of Contents
<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#key-features">Key Features</a></li>
        <li><a href="#difference-with-asdf">Difference with asdf</a></li>
      </ul>
    </li>
    <li>
        <a href="#table-of-contents">Table of contents</a>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#environment-variables">Environment variables</a></li>
    <li><a href="#version-files">Version files</a></li>
    <li><a href="#technical-details">Technical details</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#community">Community</a></li>
    <li><a href="#authors">Authors</a></li>
    <li><a href="#licence">Licence</a></li>
  </ol>
</details>


<a id="getting-started"></a>
## Getting Started

<a id="prerequisites"></a>
### Prerequisites
If you need to enable cosign checks, install `cosign` tool via one of the following commands:

<details><summary><b>MacOS (Homebrew)</b></summary><br>

```sh
brew install cosign
```

</details>


<details><summary><b>Windows (Chocolatey)</b></summary><br>

```sh
choco install cosign
```

</details>

<details><summary><b>Alpine Linux</b></summary><br>

```sh
apk add cosign
```

</details>

<details><summary><b>Arch Linux</b></summary><br>

```sh
sudo pacman -S cosign
```

</details>

<details><summary><b>Linux: RPM</b></summary><br>

```sh
LATEST_VERSION=$(curl https://api.github.com/repos/sigstore/cosign/releases/latest | jq -r .tag_name | tr -d "v")
curl -O -L "https://github.com/sigstore/cosign/releases/latest/download/cosign-${LATEST_VERSION}-1.x86_64.rpm"
sudo rpm -ivh cosign-${LATEST_VERSION}.x86_64.rpm
```

</details>
<details><summary><b>Linux: dkpg</b></summary><br>

```sh
LATEST_VERSION=$(curl https://api.github.com/repos/sigstore/cosign/releases/latest | jq -r .tag_name | tr -d "v")
curl -O -L "https://github.com/sigstore/cosign/releases/latest/download/cosign_${LATEST_VERSION}_amd64.deb"
sudo dpkg -i cosign_${LATEST_VERSION}_amd64.deb
```

</details>


<a id="installation"></a>
### Installation

<a id="automatic-installation"></a>
#### Automatic Installation
<details><summary><b>Arch Linux (AUR)</b></summary><br>

This package is available on the Arch Linux User Repository.
It can be installed using the yay AUR helper:
```sh
yay tenv-bin
```

</details>

<details><summary><b>MacOS (Homebrew)</b></summary><br>

```console
brew tap tofuutils/tap
brew install tenv
```

</details>

<details><summary><b>Windows (Chocolatey)</b></summary><br>

```console
choco install tenv
```

</details>

<details><summary><b>Ubuntu</b></summary><br>

```sh
LATEST_VERSION=$(curl --silent https://api.github.com/repos/tofuutils/tenv/releases/latest | jq -r .tag_name)
curl -O -L "https://github.com/tofuutils/tenv/releases/latest/download/tenv_${LATEST_VERSION}_amd64.deb"
sudo dpkg -i "tenv_${LATEST_VERSION}_amd64.deb"
```

</details>


<a id="manual-installation"></a>
#### Manual Installation
Get the most recent packaged binaries (`.deb`, `.rpm`, `.apk`, `pkg.tar.zst `, `.zip` or `.tar.gz` format) by visiting the [release page](https://github.com/tofuutils/tenv/releases). After downloading, unzip the folder and seamlessly integrate it into your system's `PATH`.

<a id="docker-installation"></a>

#### Docker Installation

You can use dockerized version of tenv via the following command:

```sh
docker run --rm --it ghcr.io/tofuutils/tenv:latest help
```
The docker container is not meant as a way to run tenv for CI pipelines, for local use, you should use one of the [packaged binaries](#manual-installation).
<a id="usage"></a>

## Usage

**tenv** supports [OpenTofu](https://opentofu.org), [Terragrunt](https://terragrunt.gruntwork.io/) and [Terraform](https://www.terraform.io/). To manage each binary you can use `tenv <tool> <command>`. Below is a list of tools and commands that use actual subcommands:

| tool (alias)        | env vars                   | description                                    |
| ------------------- | -------------------------- | ---------------------------------------------- |
| `tofu` (`opentofu`) | [TOFUENV_](#tofu-env-vars) | [OpenTofu](https://opentofu.org)               |
| `tf` (`terraform`)  | [TFENV_](#tf-env-vars)     | [Terraform](https://www.terraform.io/)         |
| `tg` (`terragrunt`) | [TG_](#tg-env-vars)        | [Terragrunt](https://terragrunt.gruntwork.io/) |


<details><summary><b>tenv &lt;tool&gt; install [version]</b></summary><br>

Install a requested version of the tool (into `TENV_ROOT` directory from `<TOOL>_REMOTE` url).

Without a parameter, the version to use is resolved automatically (see resolution order in [tools description](#technical-details), with `latest` as default in place of `latest-allowed`).

If a parameter is passed, available options include:

- an exact [Semver 2.0.0](https://semver.org/) version string to install.
- a [version constraint](https://opentofu.org/docs/language/expressions/version-constraints) string (checked against versions available at `<TOOL>_REMOTE` url).
- `latest`, `latest-stable` (old name of `latest`) or `latest-pre` (include unstable version), which are checked against versions available at `<TOOL>_REMOTE` url.
- `latest:<re>` or `min:<re>` to get first version matching with `<re>` as a [regexp](https://github.com/google/re2/wiki/Syntax) after a descending or ascending version sort.
- `latest-allowed` or `min-required` to scan your IAC files to detect which version is maximally allowed or minimally required. See [required_version](#required_version) docs.

```console
tenv tofu install
tenv tofu install 1.6.0-beta5
tenv tf install "~> 1.6.0"
tenv tf install latest-pre
tenv tg install latest
tenv tg install latest-stable
tenv <tool> install latest-allowed
tenv <tool> install min-required
```

A complete display :

```console
$ tenv tofu install 1.6.0
Installing OpenTofu 1.6.0
Fetching release information from https://api.github.com/repos/opentofu/opentofu/releases/tags/v1.6.0
Downloading https://github.com/opentofu/opentofu/releases/download/v1.6.0/tofu_1.6.0_linux_amd64.zip
Downloading https://github.com/opentofu/opentofu/releases/download/v1.6.0/tofu_1.6.0_SHA256SUMS
Downloading https://github.com/opentofu/opentofu/releases/download/v1.6.0/tofu_1.6.0_SHA256SUMS.sig
Downloading https://github.com/opentofu/opentofu/releases/download/v1.6.0/tofu_1.6.0_SHA256SUMS.pem
Installation of OpenTofu 1.6.0 successful
```

</details>


<details><summary><b>tenv &lt;tool&gt; use [version]</b></summary><br>

Switch the default tool version to use (set in `TENV_ROOT/<TOOL>/version` file).

`tenv <tool> use` has a `--working-dir`, `-w` flag to write a [version file](#version-files) in working directory.

Available parameter options:

- an exact [Semver 2.0.0](https://semver.org/) version string to use.
- a [version constraint](https://opentofu.org/docs/language/expressions/version-constraints) string (checked against versions available in TENV_ROOT directory).
- `latest`, `latest-stable` (old name of `latest`) or `latest-pre` (include unstable version), which are checked against versions available in TENV_ROOT directory.
- `latest:<re>` or `min:<re>` to get first version matching with `<re>` as a [regexp](https://github.com/google/re2/wiki/Syntax) after a descending or ascending version sort.
- `latest-allowed` or `min-required` to scan your IAC files to detect which version is maximally allowed or minimally required. See [required_version](#required_version) docs.

```console
tenv tofu use v1.6.0-beta5
tenv tf use min-required
tenv tg use latest
tenv tofu use latest-allowed
```

</details>


<details><summary><b>tenv &lt;tool&gt; detect</b></summary><br>

Detect the used version of tool for the working directory.

```console
$ tenv tofu detect
No version files found for OpenTofu, fallback to latest-allowed strategy
Scan project to find .tf files
No OpenTofu version requirement found in project files, fallback to latest strategy
Found compatible version installed locally : 1.6.1
OpenTofu 1.6.1 will be run from this directory.
$ tenv tg detect -q
Terragrunt 0.55.1 will be run from this directory.
```

</details>


<details><summary><b>tenv &lt;tool&gt; reset</b></summary><br>

Reset used version of tool (remove `TENV_ROOT/<TOOL>/version` file).

```console
$ tenv tofu reset
Removed /home/dvaumoron/.tenv/OpenTofu/version
```

</details>


<details><summary><b>tenv &lt;tool&gt; uninstall [version]</b></summary><br>

Uninstall a specific version of the tool (remove it from `TENV_ROOT` directory without interpretation).

```console
$ tenv tofu uninstall v1.6.0-alpha4
Uninstallation of OpenTofu 1.6.0-alpha4 successful (directory /home/dvaumoron/.tenv/OpenTofu/1.6.0-alpha4 removed)
```

</details>


<details><summary><b>tenv &lt;tool&gt; list</b></summary><br>

List installed tool versions (located in `TENV_ROOT` directory), sorted in ascending version order.

`tenv <tool> list` has a `--descending`, `-d` flag to sort in descending order.

```console
$ tenv tofu list -v
* 1.6.0 (set by /home/dvaumoron/.tenv/OpenTofu/version)
  1.6.1
found 2 OpenTofu version(s) managed by tenv.
```

</details>


<details><summary><b>tenv &lt;tool&gt; list-remote</b></summary><br>

List installable tool versions (from `<TOOL>_REMOTE` url), sorted in ascending version order.

`tenv <tool> list-remote` has a `--descending`, `-d` flag to sort in descending order.

`tenv <tool> list-remote` has a `--stable`, `-s` flag to display only stable version.

```console
$ tenv tofu list-remote
Fetching all releases information from https://api.github.com/repos/opentofu/opentofu/releases
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
1.6.0-rc1
1.6.0 (installed)
1.6.1 (installed)
```

</details>


<details><summary><b>tenv &lt;tool&gt; constraint [expression]</b></summary><br>

Set or reset a default constraint expression for the tool.

```console
$ tenv tf constraint "<= 1.5.7"
Written <= 1.5.7 in /home/dvaumoron/.tenv/Terraform/constraint
```

Or without expression :

```console
$ tenv tg constraint
Removed /home/dvaumoron/.tenv/Terragrunt/constraint
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
  -a, --arch string          specify arch for binaries downloading (default "amd64")
  -f, --force-remote         force search on versions available at TFENV_REMOTE url
  -h, --help                 help for detect
  -k, --key-file string      local path to PGP public key file (replace check against remote one)
  -n, --no-install           disable installation of missing version
  -c, --remote-conf string   path to remote configuration file (advanced settings)
  -u, --remote-url string    remote url to install from

Global Flags:
  -q, --quiet              no unnecessary output (and no log)
  -r, --root-path string   local path to install versions of OpenTofu, Terraform and Terragrunt (default "/home/dvaumoron/.tenv")
  -v, --verbose            verbose output (and set log level to Trace)
```

```console
$ tenv tofu use -h
Switch the default OpenTofu version to use (set in TENV_ROOT/OpenTofu/version file)

Available parameter options:
- an exact Semver 2.0.0 version string to use
- a version constraint expression (checked against version available in TENV_ROOT directory)
- latest, latest-stable or latest-pre (checked against version available in TENV_ROOT directory)
- latest-allowed or min-required to scan your OpenTofu files to detect which version is maximally allowed or minimally required.

Usage:
  tenv tofu use version [flags]

Flags:
  -a, --arch string           specify arch for binaries downloading (default "amd64")
  -f, --force-remote          force search on versions available at TOFUENV_REMOTE url
  -t, --github-token string   GitHub token (increases GitHub REST API rate limits)
  -h, --help                  help for use
  -k, --key-file string       local path to PGP public key file (replace check against remote one)
  -n, --no-install            disable installation of missing version
  -c, --remote-conf string    path to remote configuration file (advanced settings)
  -u, --remote-url string     remote url to install from
  -w, --working-dir           create .opentofu-version file in working directory

Global Flags:
  -q, --quiet              no unnecessary output (and no log)
  -r, --root-path string   local path to install versions of OpenTofu, Terraform and Terragrunt (default "/home/dvaumoron/.tenv")
  -v, --verbose            verbose output (and set log level to Trace)
```

</details>


<details><summary><b>tenv update-path</b></summary><br>

Display PATH updated with tenv directory location first. With GITHUB_ACTIONS set to true, write tenv directory location to GITHUB_PATH.

This command can be used when one of the managed tool is already installed on your system and hide the corresponding proxy (in that case `which tenv` and `which <tool>` will indicate different locations). The following shell call should resolve such issues :

```console
export PATH=$(tenv update-path)
```

</details>


<details><summary><b>tenv version</b></summary><br>

Display tenv current version.

```console
$ tenv version
tenv version v1.7.0
```


</details>


<a id="environment-variables"></a>
## Environment variables

**tenv** commands support global environment variables and variables by tool for : [OpenTofu](https://opentofu.org), [Terraform](https://www.terraform.io/) and [TerraGrunt](https://terragrunt.gruntwork.io/).


<a id="tenv-vars"></a>
### Global tenv environment variables

<details><summary><b>TENV_ARCH</b></summary><br>

String (Default: current tenv binaries architecture)

Allow to override the default architecture for binaries downloading during installation.

`tenv <tool>` subcommands `detect`, `install` and `use` support a `--arch`, `-a`  flag version.

</details>


<details><summary><b>TENV_AUTO_INSTALL</b></summary><br>

String (Default: true)

If set to true **tenv** will automatically install missing tool versions needed.

`tenv <tool>` subcommands `detect` and `use` support a `--no-install`, `-n` disabling flag version.

</details>


<details><summary><b>TENV_FORCE_REMOTE</b></summary><br>

String (Default: false)

If set to true **tenv** detection of needed version will skip local check and verify compatibility on remote list.

`tenv <tool>` subcommands `detect` and `use` support a `--force-remote`, `-f` flag version.

</details>


<details><summary><b>TENV_GITHUB_TOKEN</b></summary><br>

String (Default: "")

Allow to specify a GitHub token to increase [GitHub Rate limits for the REST API](https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api). Useful because OpenTofu and Terragrunt binaries are downloaded from GitHub repository.

`tenv tofu` and `tenv tg` subcommands `detect`, `install`, `list-remote` and `use` support a `--github-token`, `-t` flag version.

</details>


<details><summary><b>TENV_QUIET</b></summary><br>

String (Default: false)

If set to true **tenv** disable unnecessary output (including log level forced to off).

`tenv` subcommands support a `--quiet`, `-q` flag version.

</details>


<details><summary><b>TENV_LOG</b></summary><br>

String (Default: "warn")

Set **tenv** log level (possibilities sorted by decreasing verbosity : "trace", "debug", "info", "warn", "error", "off").

`tenv` support a `--verbose`, `-v` flag which set log level to "trace".

</details>


<details><summary><b>TENV_REMOTE_CONF</b></summary><br>

String (Default: `${TENV_ROOT}/remote.yaml`)

The path to a yaml file for [advanced remote configuration](#advanced-remote-configuration) (can be used to call artifact mirror).

`tenv <tool>` subcommands `detect`, `install`, `list-remote` and `use`  support a `--remote-conf`, `-c` flag version.

</details>


<details><summary><b>TENV_ROOT</b></summary><br>

String (Default: `${HOME}/.tenv`)

The path to a directory where the local OpenTofu versions, Terraform versions, Terragrunt versions and tenv configuration files exist.

`tenv` support a `--root-path`, `-r` flag version.

</details>


<details><summary><b>GITHUB_ACTIONS</b></summary><br>

String (Default: false)

If set to true **tenv** proxies exposes proxied output `stdout`, `stderr`, and `exitcode` by writing them into GITHUB_OUTPUT in [multiline format](https://docs.github.com/en/actions/using-workflows/workflow-commands-for-github-actions#multiline-strings). GitHub Actions set it (see [default environment variables](https://docs.github.com/en/actions/learn-github-actions/variables#default-environment-variables)).

</details>


<details><summary><b>GITHUB_OUTPUT</b></summary><br>

String (Default: "")

Needed when GITHUB_ACTIONS is set to true, path to a file to write proxied output.

</details>


<details><summary><b>GITHUB_PATH</b></summary><br>

String (Default: "")

Used by `tenv update-path` when GITHUB_ACTIONS is set to true, path to a file to write tenv directory location.

</details>


<a id="tofu-env-vars"></a>
### OpenTofu environment variables

<details><summary><b>TOFUENV_ARCH</b></summary><br>

Same as TENV_ARCH (compatibility with [tofuenv](https://github.com/tofuutils/tofuenv)).

</details>


<details><summary><b>TOFUENV_AUTO_INSTALL</b></summary><br>

Same as TENV_AUTO_INSTALL (compatibility with [tofuenv](https://github.com/tofuutils/tofuenv)).

#### Example 1
Use OpenTofu version 1.6.1 that is not installed, and auto installation is disabled :

```console
$ TOFUENV_AUTO_INSTALL=false tenv use 1.6.1
Written 1.6.1 in /home/dvaumoron/.tenv/OpenTofu/version
```

#### Example 2
Use OpenTofu version 1.6.0 that is not installed, and auto installation stay enabled :

```console
$ tenv tofu use 1.6.0
Installing OpenTofu 1.6.0
Fetching release information from https://api.github.com/repos/opentofu/opentofu/releases/tags/v1.6.0
Downloading https://github.com/opentofu/opentofu/releases/download/v1.6.0/tofu_1.6.0_linux_amd64.zip
Downloading https://github.com/opentofu/opentofu/releases/download/v1.6.0/tofu_1.6.0_SHA256SUMS
Downloading https://github.com/opentofu/opentofu/releases/download/v1.6.0/tofu_1.6.0_SHA256SUMS.sig
Downloading https://github.com/opentofu/opentofu/releases/download/v1.6.0/tofu_1.6.0_SHA256SUMS.pem
Installation of OpenTofu 1.6.0 successful
Written 1.6.0 in /home/dvaumoron/.tenv/OpenTofu/version
```

</details>


<details><summary><b>TOFUENV_FORCE_REMOTE</b></summary><br>

Same as TENV_FORCE_REMOTE.

</details>


<details><summary><b>TOFUENV_INSTALL_MODE</b></summary><br>

String (Default: "api")

- "api" install mode retrieve download url of OpenTofu from [Github REST API](https://docs.github.com/en/rest?apiVersion=2022-11-28) (TOFUENV_REMOTE must comply with it).
- "direct" install mode generate download url of OpenTofu based on TOFUENV_REMOTE.

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details><summary><b>TOFUENV_LIST_MODE</b></summary><br>

String (Default: "api")

- "api" list mode retrieve information of OpenTofu releases from [Github REST API](https://docs.github.com/en/rest?apiVersion=2022-11-28) (TOFUENV_LIST_URL must comply with it).
- "html" list mode extract information of OpenTofu releases from parsing an html page in TOFUENV_LIST_URL.

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details><summary><b>TOFUENV_LIST_URL</b></summary><br>

String (Default: "")

Allow to override the remote url only for the releases listing, default value depend on TOFUENV_LIST_MODE :

- with "api" mode use default API URL (https://api.github.com/repos/opentofu/opentofu/releases)
- with "html" mode same as TOFUENV_REMOTE

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details><summary><b>TOFUENV_OPENTOFU_PGP_KEY</b></summary><br>

String (Default: "")

Allow to specify a local file path to OpenTofu PGP public key, if not present download https://get.opentofu.org/opentofu.asc.

`tenv tofu` subcommands `detect`, `ìnstall` and `use` support a `--key-file`, `-k` flag version.

</details>


<details><summary><b>TOFUENV_REMOTE</b></summary><br>

String (Default: https://api.github.com/repos/opentofu/opentofu/releases)

URL to install OpenTofu, when TOFUENV_REMOTE differ from its default value, TOFUENV_INSTALL_MODE is set to direct (assume an artifact proxy usage, however releases listing continue to use API).

`tenv tofu` subcommands `detect`, `install`, `list-remote` and `use` support a `--remote-url`, `-u` flag version.

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details><summary><b>TOFUENV_ROOT</b></summary><br>

Same as TENV_ROOT (compatibility with [tofuenv](https://github.com/tofuutils/tofuenv)).

</details>


<details><summary><b>TOFUENV_GITHUB_TOKEN</b></summary><br>

Same as TENV_GITHUB_TOKEN (compatibility with [tofuenv](https://github.com/tofuutils/tofuenv)).

</details>


<details><summary><b>TOFUENV_TOFU_DEFAULT_CONSTRAINT</b></summary><br>

String (Default: "")

If not empty string, this variable overrides OpenTofu default constraint, specified in ${TENV_ROOT}/OpenTofu/constraint file.

</details>


<details><summary><b>TOFUENV_TOFU_DEFAULT_VERSION</b></summary><br>

String (Default: "")

If not empty string, this variable overrides OpenTofu fallback version, specified in ${TENV_ROOT}/OpenTofu/version file.

</details>


<details><summary><b>TOFUENV_TOFU_VERSION</b></summary><br>

String (Default: "")

If not empty string, this variable overrides OpenTofu version, specified in [`.opentofu-version`](#opentofu-version-files) files.

`tenv tofu` subcommands `install` and `detect` also respects this variable.

e.g. with :

```console
$ tofu version
OpenTofu v1.6.1
on linux_amd64
```

then :

```console
$ TOFUENV_TOFU_VERSION=1.6.0 tofu version
OpenTofu v1.6.0
on linux_amd64
```

</details>


<a id="tf-env-vars"></a>
### Terraform environment variables

<details><summary><b>TFENV_ARCH</b></summary><br>

Same as TENV_ARCH (compatibility with [tfenv](https://github.com/tfutils/tfenv)).

</details>


<details><summary><b>TFENV_AUTO_INSTALL</b></summary><br>

Same as TENV_AUTO_INSTALL (compatibility with [tfenv](https://github.com/tfutils/tfenv)).

</details>


<details><summary><b>TFENV_FORCE_REMOTE</b></summary><br>

Same as TENV_FORCE_REMOTE.

</details>


<details><summary><b>TFENV_HASHICORP_PGP_KEY</b></summary><br>

String (Default: "")

Allow to specify a local file path to Hashicorp PGP public key, if not present download https://www.hashicorp.com/.well-known/pgp-key.txt.

`tenv tf` subcommands `detect`, `ìnstall` and `use` support a `--key-file`, `-k` flag version.

</details>


<details><summary><b>TFENV_INSTALL_MODE</b></summary><br>

String (Default: "api")

- "api" install mode retrieve download url of Terraform from [Hashicorp Release API](https://releases.hashicorp.com/docs/api/v1) (TFENV_REMOTE must comply with it).
- "direct" install mode generate download url of Terraform based on TFENV_REMOTE.

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details><summary><b>TFENV_LIST_MODE</b></summary><br>

String (Default: "api")

- "api" list mode retrieve information of Terraform releases from [Hashicorp Release API](https://releases.hashicorp.com/docs/api/v1) (TFENV_LIST_URL must comply with it).
- "html" list mode extract information of Terraform releases from parsing an html page in TFENV_LIST_URL.

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details><summary><b>TFENV_LIST_URL</b></summary><br>

String (Default: "")

Allow to override the remote url only for the releases listing, default value depend on TFENV_LIST_MODE :

- with "api" mode use default API URL (https://releases.hashicorp.com)
- with "html" mode same as TFENV_REMOTE

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details><summary><b>TFENV_REMOTE</b></summary><br>

String (Default: https://releases.hashicorp.com)

URL to install Terraform (can differ from its default, an artifact proxy will not disturb the retrieving of release index.json, however releases listing continue to use API).

`tenv tf` subcommands `detect`, `install`, `list-remote` and `use` support a `--remote-url`, `-u` flag version.

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details><summary><b>TFENV_ROOT</b></summary><br>

Same as TENV_ROOT (compatibility with [tfenv](https://github.com/tfutils/tfenv)).

</details>


<details><summary><b>TFENV_TERRAFORM_DEFAULT_CONSTRAINT</b></summary><br>

String (Default: "")

If not empty string, this variable overrides Terraform default constraint, specified in ${TENV_ROOT}/Terraform/constraint file.

</details>


<details><summary><b>TFENV_TERRAFORM_DEFAULT_VERSION</b></summary><br>

String (Default: "")

If not empty string, this variable overrides Terraform fallback version, specified in ${TENV_ROOT}/Terraform/version file.

</details>


<details><summary><b>TFENV_TERRAFORM_VERSION</b></summary><br>

String (Default: "")

If not empty string, this variable overrides Terraform version, specified in [`.terraform-version`](#terraform-version-files) files.

`tenv tf` subcommands `install` and `detect` also respects this variable.

e.g. with :

```console
$ terraform version
Terraform v1.7.2
on linux_amd64
```

then :

```console
$ TFENV_TERRAFORM_VERSION=1.7.0 terraform version
Terraform v1.7.0
on linux_amd64

Your version of Terraform is out of date! The latest version
is 1.7.2. You can update by downloading from https://www.terraform.io/downloads.html
```

</details>


<a id="tg-env-vars"></a>
### Terragrunt environment variables


<details><summary><b>TG_INSTALL_MODE</b></summary><br>

String (Default: "api")

- "api" install mode retrieve download url of Terragrunt from [Github REST API](https://docs.github.com/en/rest?apiVersion=2022-11-28) (TG_REMOTE must comply with it).
- "direct" install mode generate download url of Terragrunt based on TG_REMOTE.

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details><summary><b>TG_LIST_MODE</b></summary><br>

String (Default: "api")

- "api" list mode retrieve information of Terragrunt releases from [Github REST API](https://docs.github.com/en/rest?apiVersion=2022-11-28) (TG_LIST_URL must comply with it).
- "html" list mode extract information of Terragrunt releases from parsing an html page in TG_LIST_URL.

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details><summary><b>TG_LIST_URL</b></summary><br>

String (Default: "")

Allow to override the remote url only for the releases listing, default value depend on TG_LIST_MODE :

- with "api" mode use default API URL (https://api.github.com/repos/gruntwork-io/terragrunt/releases)
- with "html" mode same as TG_REMOTE

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details><summary><b>TG_REMOTE</b></summary><br>

String (Default: https://api.github.com/repos/gruntwork-io/terragrunt/releases)

URL to install Terragrunt, when TG_REMOTE differ from its default value, TG_INSTALL_MODE is set to "direct" (assume an artifact proxy usage, however releases listing continue to use API).

`tenv tg` subcommands `detect`, `install`, `list-remote` and `use` support a `--remote-url`, `-u` flag version.

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>

<details><summary><b>TG_DEFAULT_CONSTRAINT</b></summary><br>

String (Default: "")

If not empty string, this variable overrides Terragrunt default constraint, specified in ${TENV_ROOT}/Terragrunt/constraint file.

</details>


<details><summary><b>TG_DEFAULT_VERSION</b></summary><br>

String (Default: "")

If not empty string, this variable overrides Terragrunt fallback version, specified in ${TENV_ROOT}/Terragrunt/version file.

</details>


<details><summary><b>TG_VERSION</b></summary><br>

String (Default: "")

If not empty string, this variable overrides Terragrunt version, specified in [`.terragrunt-version`](#terragrunt-version-files) files.

`tenv tg` subcommands `install` and `detect` also respects this variable.

e.g. with :

```console
$ terragrunt -v
terragrunt version v0.55.1
```

then :

```console
$ TG_VERSION=0.54.1 terragrunt -v
terragrunt version v0.54.1
```

</details>


<a id="version-files"></a>
## version files

<a id="default-version-file"></a>
<details><summary><b>default version file</b></summary><br>

The `TENV_ROOT/<TOOL>/version` file is the tool default version used when no project specific or user specific are found. It can be written with `tenv <tool> use`.

</details>

<a id="opentofu-version-files"></a>
<details><summary><b>opentofu version files</b></summary><br>

If you put a `.opentofu-version` file in the working directory, one of its parent directory, or user home directory, **tenv** detects it and uses the version written in it.
Note, that TOFUENV_TOFU_VERSION can be used to override version specified by `.opentofu-version` file.

Recognize same values as `tenv tofu use` command.

See [required_version](https://opentofu.org/docs/language/settings#specifying-a-required-opentofu-version) docs.

</details>

<a id="terraform-version-files"></a>
<details><summary><b>terraform version files</b></summary><br>

If you put a `.terraform-version` or `.tfswitchrc` file in the working directory, one of its parent directory, or user home directory, **tenv** detects it and uses the version written in it.
Note, that TFENV_TERRAFORM_VERSION can be used to override version specified by those files.

Recognize same values as `tenv tf use` command.

See [required_version](https://developer.hashicorp.com/terraform/language/settings#specifying-a-required-terraform-version) docs.

</details>

<a id="terragrunt-version-files"></a>
<details><summary><b>terragrunt version files</b></summary><br>

If you put a `.terragrunt-version` or a `.tgswitchrc` file in the working directory, one of its parent directory, or user home directory, **tenv** detects it and uses the version written in it. **tenv** also detect a `version` field in a `.tgswitch.toml` in same places.
Note, that TG_VERSION can be used to override version specified by those files.

Recognize same values as `tenv tg use` command.

</details>


<a id="terragrunt-hcl-file"></a>
<details><summary><b>terragrunt.hcl file</b></summary><br>

If you have a terragrunt.hcl or terragrunt.hcl.json in the working directory, one of its parent directory, or user home directory, **tenv** will read constraint from `terraform_version_constraint` or `terragrunt_version_constraint` field in it (depending on proxy or subcommand used).

</details>

<a id="required_version"></a>
<details><summary><b>required_version</b></summary><br>

the `latest-allowed` or `min-required` strategies scan through your IAC files (.tf or .tf.json) and identify a version conforming to the constraint in the relevant files. They fallback to `latest` when no IAC files and no default constraint are found.

Currently the format for [Terraform required_version](https://developer.hashicorp.com/terraform/language/settings#specifying-a-required-terraform-version) and [OpenTofu required_version](https://opentofu.org/docs/language/settings#specifying-a-required-opentofu-version) are very similar, however this may change over time, always refer to docs for the latest format specification.

example:

```HCL
version = ">= 1.2.0, < 2.0.0"
```

This would identify the latest version at or above 1.2.0 and below 2.0.0

</details>

<a id="technical-details"></a>
## Technical details

### Project binaries

<details><summary><b>tofu</b></summary><br>

The `tofu` command in this project is a proxy to OpenTofu's `tofu` command  managed by **tenv**.

The version resolution order is :

- TOFUENV_TOFU_VERSION environment variable
- `.opentofu-version` file
- `terraform_version_constraint` from `terragrunt.hcl` file
- `terraform_version_constraint` from `terragrunt.hcl.json` file
- TOFUENV_TOFU_DEFAULT_VERSION environment variable
- `${TENV_ROOT}/OpenTofu/version` file (can be written with `tenv tofu use`)
- `latest-allowed`

The `latest-allowed` strategy rely on [required_version](#required_version) from .tf or .tf.json files with a fallback to `latest` when no constraint are found. Moreover it is possible to add a default constraint with TOFUENV_TOFU_DEFAULT_CONSTRAINT environment variable or `${TENV_ROOT}/OpenTofu/constraint` file (can be written with `tenv tofu constraint`). The default constraint is added while using `latest-allowed`, `min-required` or custom constraint. A default constraint with `latest-allowed` or `min-required` will avoid the fallback to `latest` when there is no .tf or .tf.json files.

</details>

<details><summary><b>terraform</b></summary><br>

The `terraform` command in this project is a proxy to HashiCorp's `terraform` command managed by **tenv**.

The version resolution order is :

- TFENV_TERRAFORM_VERSION environment variable
- `.terraform-version` file
- `.tfswitchrc` file
- `terraform_version_constraint` from `terragrunt.hcl` file
- `terraform_version_constraint` from `terragrunt.hcl.json` file
- TFENV_TERRAFORM_DEFAULT_VERSION environment variable
- `${TENV_ROOT}/Terraform/version` file (can be written with `tenv tf use`)
- `latest-allowed`

The `latest-allowed` strategy rely on [required_version](#required_version) from .tf or .tf.json files with a fallback to `latest` when no constraint are found. Moreover it is possible to add a default constraint with TFENV_TERRAFORM_DEFAULT_CONSTRAINT environment variable or `${TENV_ROOT}/Terraform/constraint` file (can be written with `tenv tf constraint`). The default constraint is added while using `latest-allowed`, `min-required` or custom constraint. A default constraint with `latest-allowed` or `min-required` will avoid the fallback to `latest` when there is no .tf or .tf.json files.

</details>


<details><summary><b>terragrunt</b></summary><br>

The `terragrunt` command in this project is a proxy to Gruntwork's `terragrunt` command managed by **tenv**.

The version resolution order is :

- TG_VERSION environment variable
- `.terragrunt-version` file
- `.tgswitchrc` file
- `version` from `tgswitch.toml` file
- `terragrunt_version_constraint` from `terragrunt.hcl` file
- `terragrunt_version_constraint` from `terragrunt.hcl.json` file
- TG_DEFAULT_VERSION environment variable
- `${TENV_ROOT}/Terragrunt/version` file (can be written with `tenv tg use`)
- `latest-allowed`

The `latest-allowed` strategy has no information for Terragrunt and will fallback to `latest` unless there is default constraint. Adding a default constraint could be done with TG_DEFAULT_CONSTRAINT environment variable or `${TENV_ROOT}/Terragrunt/constraint` file (can be written with `tenv tg constraint`). The default constraint is added while using `latest-allowed`, `min-required` or custom constraint. A default constraint with `latest-allowed` or `min-required` will avoid there fallback to `latest`.

</details>


<details><summary><b>tf</b></summary><br>

The `tf` command is a proxy to `tofu` or `terraform` depending on the version files present in project.

The version resolution order is :

- `.opentofu-version` file (launch `tofu`)
- `terraform_version_constraint` from `terragrunt.hcl` file (launch `tofu`)
- `terraform_version_constraint` from `terragrunt.hcl.json` file (launch `tofu`)
- `.terraform-version` file (launch `terraform`)
- `.tfswitchrc` file  (launch `terraform`)
- fail with a message


</details>


<a id="advanced-remote-configuration"></a>
### Advanced remote configuration

This advanced configuration is meant to call artifact mirror (like [JFrog Artifactory](https://jfrog.com/artifactory)).

The yaml file from TENV_REMOTE_CONF path can have one part for each supported proxy : `tofu`, `terraform` and `terragrunt`.

<details><summary><b>yaml fields description</b></summary><br>

Each part can have the following string field : `install_mode`, `list_mode`, `list_url`, `url`, `new_base_url`, `old_base_url`, `selector` and `part`

With `install_mode` set to "direct", **tenv** skip the release information fetching and generate download url instead of reading them from API (overridden by `<TOOL>_INSTALL_MODE` env var).

With `list_mode` set to "html", **tenv** change the fetching of all releases information from API to parse the parent html page of artifact location, see `selector` and `part` (overridden by `<TOOL>_LIST_MODE` env var).

`url` allows to override the default remote url (overridden by flag or `<TOOL>_REMOTE` env var).

`list_url` allows to override the remote url only for the releases listing (overridden by `<TOOL>_LIST_URL` env var).

`old_base_url` and `new_base_url` are used as url rewrite rule (if an url start with the prefix, it will be changed to use the new base url).

If `old_base_url` and `new_base_url` are empty, **tenv** try to guess right behaviour based on previous fields.

`selector` is used to gather in a list all matching html node and `part` choose on which node part (attribute name or "#text" for inner text) a version will be extracted (selector default to "a" (html link) and part default to "href" (link target))

</details>


<details><summary><b>Examples</b></summary><br>

Those examples assume that a GitHub proxy at https://artifactory.example.com/artifactory/github have the same behavior than [JFrog Artifactory](https://jfrog.com/artifactory) :

- mirror https://github.com/opentofu/opentofu/releases/download/v1.6.0/tofu_1.6.0_linux_amd64.zip at https://artifactory.example.com/artifactory/github/opentofu/opentofu/releases/download/v1.6.0/tofu_1.6.0_linux_amd64.zip.
- have at https://artifactory.example.com/artifactory/github/opentofu/opentofu/releases/download an html page with links on existing sub folder like "v1.6.0/"

Example 1 : Retrieve Terraform binaries and list available releases from the mirror.

```console
TFENV_REMOTE=https://artifactory.example.com/artifactory/hashicorp
TFENV_LIST_MODE=html
```

Example 2 : Retrieve Terraform binaries from the mirror and list available releases from the Hashicorp releases API (TFENV_LIST_URL is optional because it default to https://releases.hashicorp.com with the default list mode "api").

```console
TFENV_REMOTE=https://artifactory.example.com/artifactory/hashicorp
TFENV_LIST_URL=https://releases.hashicorp.com
```

Example 1 & 2, does not need install mode (by release index.json is figed in mirror without problem), however create a rewrite rule from "https://releases.hashicorp.com" to "https://artifactory.example.com/artifactory/hashicorp" to obtains correct download URLs.

Example 3 : Retrieve OpenTofu binaries and list available releases from the mirror (TOFUENV_INSTALL_MODE is optional because overloading TOFUENV_REMOTE already set it to "direct").

```console
TOFUENV_REMOTE=https://artifactory.example.com/artifactory/github
TOFUENV_INSTALL_MODE=direct
TOFUENV_LIST_MODE=html
```

Example 4 : Retrieve OpenTofu binaries from the mirror and list available releases from the GitHub API (TOFUENV_INSTALL_MODE is optional because overloading TOFUENV_REMOTE already set it to "direct", and TOFUENV_LIST_URL is optional because it default to https://api.github.com/repos/opentofu/opentofu/releases with the default list mode "api").

```console
TOFUENV_REMOTE=https://artifactory.example.com/artifactory/github
TOFUENV_INSTALL_MODE=direct
TOFUENV_LIST_URL=https://api.github.com/repos/opentofu/opentofu/releases
```

Example 3 & 4, does not create a rewrite rule (the direct install mode build correct download URLs).

Example 1 & 4 can be merged in a remote.yaml :

```yaml
tofu:
  url: "https://artifactory.example.com/artifactory/github"
  install_mode: "direct"
  list_url: "https://api.github.com/repos/opentofu/opentofu/releases"
terraform:
  url: "https://artifactory.example.com/artifactory/hashicorp"
  list_mode: "html"
```

</details>


<a id="signature-support"></a>
### Signature support

<details><summary><b>OpenTofu signature support</b></summary><br>

**tenv** checks the sha256 checksum and the signature of the checksum file with [cosign](https://github.com/sigstore/cosign) (if present on your machine) or PGP (via [gopenpgp](https://github.com/ProtonMail/gopenpgp)). However, unstable OpenTofu versions are signed only with cosign (in this case, if cosign is not found tenv will display a warning).

</details>

<details><summary><b>Terraform signature support</b></summary><br>

**tenv** checks the sha256 checksum and the PGP signature of the checksum file (via [gopenpgp](https://github.com/ProtonMail/gopenpgp), there is no cosign signature available).

</details>

<details><summary><b>Terragrunt signature support</b></summary><br>

**tenv** checks the sha256 checksum (there is no signature available).

</details>

<a id="contributing"></a>
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


<a id="community"></a>
## Community
Have questions or suggestions? Reach out to us via:

* [GitHub Issues](LINK_TO_ISSUES)
* User/Developer Group: Join github community to get update of Harbor's news, features, releases, or to provide suggestion and feedback.
* Slack: Join tofuutils's community for discussion and ask questions: OpenTofu, channel: #tofuutils


<a id="authors"></a>
## Authors
tenv is based on [tofuenv](https://github.com/tofuutils/tofuenv) and [gotofuenv](https://github.com/tofuutils/gotofuenv) projects and supported by tofuutils team with help from these awesome contributors:

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

<a id="licence"></a>
## LICENSE
The tenv project is distributed under the Apache 2.0 license. See [LICENSE](LICENSE).
