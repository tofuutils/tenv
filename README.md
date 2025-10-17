<!-- BADGES -->
<p align="center">
  <a href="https://github.com/tofuutils/tenv/releases"><img src="https://img.shields.io/github/v/release/tofuutils/tenv?style=for-the-badge" alt="Github release"></a>
  <a href="https://github.com/tofuutils/tenv/graphs/contributors"><img src="https://img.shields.io/github/contributors/tofuutils/tenv?style=for-the-badge" alt="Contributors"></a>
  <img src="https://img.shields.io/maintenance/yes/2025.svg?style=for-the-badge" alt="maintenance status">
  <a href="https://goreportcard.com/report/github.com/tofuutils/tenv/"><img src="https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=for-the-badge" alt="Go report"></a>
  <a href="https://codecov.io/gh/tofuutils/tenv"><img src="https://img.shields.io/codecov/c/github/tofuutils/tenv?token=BDU9X0BAZV&style=for-the-badge" alt="codecov"></a>
</p>


<!-- LOGO -->
<br />
<div align="center">
  <a>
    <img src="assets/logo.png" alt="Logo" width="200" height="200">
  </a>
  <h3 align="center">tenv</h3>
  <p align="center">
    OpenTofu, Terraform, Terragrunt, Terramate and Atmos version manager, written in Go.
    <br />
    ·
    <a href="https://github.com/tofuutils/tenv/issues/new?assignees=&labels=issue%3A+bug&projects=&template=bug_report.md&title=">Report Bug</a>
    ·
    <a href="https://github.com/tofuutils/tenv/issues/new?assignees=&labels=&projects=&template=feature_request.md&title=">Request Feature</a>
    ·
  </p>
</div>

<p align="center">
    <a href="https://devhunt.org/tool/tenv" title="DevHunt - Tool of the Week"><img src="./assets/devhunt-badge.png" width=160 alt="DevHunt - Tool of the Week" /></a>
</p>

<a id="about-the-project"></a>
## About The Project

Welcome to **tenv**, a versatile version manager for [OpenTofu](https://opentofu.org),
[Terraform](https://www.terraform.io/), [Terragrunt](https://terragrunt.gruntwork.io/), [Terramate](https://terramate.io/) and
[Atmos](https://atmos.tools), written in Go. Our tool simplifies the complexity of handling different versions of these powerful tools, ensuring developers and DevOps professionals can focus on what matters most - building and deploying efficiently.

**tenv** is a successor of [tofuenv](https://github.com/tofuutils/tofuenv) and [tfenv](https://github.com/tfutils/tfenv).

<a id="key-features"></a>
### Key Features

- Versatile version management: Easily switch between different versions of OpenTofu,
Terraform, Terragrunt, Terramate and Atmos.
- [Semver 2.0.0](https://semver.org/) Compatibility: Utilizes [go-version](https://github.com/hashicorp/go-version) for semantic versioning and use the [HCL](https://github.com/hashicorp/hcl) parser to extract required version constraint from OpenTofu/Terraform/Terragrunt files (see [required_version](#required_version) and [Terragrunt hcl](#terragrunt-hcl-file)).
- Signature verification: Supports [cosign](https://github.com/sigstore/cosign) (if present on your machine) and PGP (via [gopenpgp](https://github.com/ProtonMail/gopenpgp)), see [signature support](#signature-support).
- Intuitive installation: Simple installation process with Homebrew and manual options.
- Callable as [Go](https://go.dev) module, with a [Semver compatibility promise](https://semver.org/#summary) on [tenvlib](https://github.com/tofuutils/tenv/tree/main/versionmanager/tenvlib) wrapper package (get more information in [TENV_AS_LIB.md](https://github.com/tofuutils/tenv/blob/main/TENV_AS_LIB.md)).

<a id="difference-with-asdf"></a>
### Difference with asdf

[asdf-vm](https://asdf-vm.com/) share the same goals than **tenv** : simplify the usage of several version of tools.

asdf-vm is generic and extensible with a plugin system, key **tenv** differences :
- **tenv** is more specific and has features dedicated to OpenTofu, Terraform, Terragrunt, Terramate 
and Atmos, like [HCL](https://github.com/hashicorp/hcl) parsing based detection (see [Key Features](#key-features)).
- **tenv** is distributed as independent binaries and does not rely on any shell or other CLI executable.
- **tenv** does better in terms of performance and platform compatibility. It works uniformly across all modern operating systems,
including Linux, MacOS, Windows, BSD, and Solaris, whereas asdf-vm natively supports only Linux and MacOS.
- **tenv** checks the sha256 checksum and the signature of the checksum file with [cosign](https://github.com/sigstore/cosign). Check [Signature support](#signature-support) section for getting more information about it.
- **tenv** command compatibility: In nearly all places you can use the exact syntax that works in tfenv / tofuenv.
If you're coming from tfenv/tofuenv and comfortable with that way of working you can almost always use the same syntax with tenv.
- **tenv** performance: It sounds incredibly useful, though it might be tough to get a real apples to apples comparison since the
tools work differently in a lot of ways. The author of asdf did a great writeup of performance problems.

<a id="table-of-contents"></a>
## Table of Contents
<!-- TABLE OF CONTENTS -->
<details markdown="1">
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
    <li><a href="#verifying-signature">Verifying tenv Signatures</a></li>
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
If you need to enable cosign checks, install `cosign` (v.2.0+) tool via one of the following commands:

<details markdown="1"><summary><b>MacOS (Homebrew, MacPorts)</b></summary><br>

Installation via Homebrew:
```sh
brew install cosign
```

Installation via MacPorts:
```sh
sudo port install cosign
```

</details>


<details markdown="1"><summary><b>Windows (go install)</b></summary><br>

```sh
go install github.com/sigstore/cosign/v2/cmd/cosign@latest
```

</details>

<details markdown="1"><summary><b>Alpine Linux</b></summary><br>

```sh
apk add cosign
```

</details>

<details markdown="1"><summary><b>Arch Linux</b></summary><br>

```sh
sudo pacman -S cosign
```

</details>

<details markdown="1"><summary><b>Linux: RPM</b></summary><br>

```sh
LATEST_VERSION=$(curl https://api.github.com/repos/sigstore/cosign/releases/latest | jq -r .tag_name | tr -d "v")
curl -O -L "https://github.com/sigstore/cosign/releases/latest/download/cosign-${LATEST_VERSION}-1.x86_64.rpm"
sudo rpm -ivh cosign-${LATEST_VERSION}-1.x86_64.rpm
```

</details>
<details markdown="1"><summary><b>Linux: dkpg</b></summary><br>

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
<details markdown="1"><summary><b>Arch Linux (AUR, Nix)</b></summary><br>

This package is available on the Arch Linux User Repository.
It can be installed using the yay AUR helper:
```sh
yay tenv-bin
```

Installation via Nix package manager:
```sh
nix-env -i tenv
```

</details>

<details markdown="1"><summary><b>MacOS (Homebrew, MacPorts, Nix)</b></summary><br>

Installation via Homebrew:
```console
brew install tenv
```

Installation via MacPorts:
```console
sudo port install tenv
```

Installation via Nix package manager:
```console
nix-env -i tenv
```

</details>

<details markdown="1"><summary><b>Windows (Chocolatey, Scoop, Nix)</b></summary><br>

Installation via Chocolatey:
```console
choco install tenv
```

Installation via Scoop:
```console
scoop install tenv
```

Installation via Nix package manager:
```console
nix-env -i tenv
```

</details>

<details markdown="1"><summary><b>Linux: Snapcraft</b></summary><br>

```sh
snap install tenv
```

</details>

<details markdown="1"><summary><b>Alpine</b></summary><br>

```sh
apk add tenv --repository=http://dl-cdn.alpinelinux.org/alpine/edge/testing/
```

#### Installation via Cloudsmith Artifact Management platform.

Setup Cloudsmith repository [automatically](https://help.cloudsmith.io/docs/alpine-repository#public-repositories):
```sh
sudo apk add --no-cache bash
curl -1sLf 'https://dl.cloudsmith.io/public/tofuutils/tenv/cfg/setup/bash.alpine.sh' | sudo bash
```
Install via apk:
```sh
apk add tenv
```

</details>

<details markdown="1"><summary><b>Ubuntu</b></summary><br>

#### Install via dpkg

```sh
LATEST_VERSION=$(curl --silent https://api.github.com/repos/tofuutils/tenv/releases/latest | jq -r .tag_name)
curl -O -L "https://github.com/tofuutils/tenv/releases/latest/download/tenv_${LATEST_VERSION}_amd64.deb"
sudo dpkg -i "tenv_${LATEST_VERSION}_amd64.deb"
```

#### Install via Nix
Installation via Nix package manager:
```console
nix-env -i tenv
```

#### Installation via Cloudsmith Artifact Management platform.
Setup Cloudsmith repository [automatically](https://help.cloudsmith.io/docs/debian-repository#public-repositories):
```sh
curl -1sLf 'https://dl.cloudsmith.io/public/tofuutils/tenv/cfg/setup/bash.deb.sh' | sudo bash
```
Install via apt:
```sh
sudo apt install tenv
```

</details>

<details markdown="1"><summary><b>RedHat</b></summary><br>

#### Installation via Cloudsmith Artifact Management platform.
Setup Cloudsmith repository [automatically](https://help.cloudsmith.io/docs/redhat-repository#public-repositories):
```sh
curl -1sLf 'https://dl.cloudsmith.io/public/tofuutils/tenv/cfg/setup/bash.rpm.sh' | sudo bash
```
Install via yum:
```sh
sudo yum install tenv
```

</details>

<details markdown="1"><summary><b>NixOS</b></summary>

#### nix-env

```sh
nix-env -iA nixos.tenv
```

#### NixOS Configuration
Add the following Nix code to your NixOS Configuration, usually located in /etc/nixos/configuration.nix
```console
environment.systemPackages = [
    pkgs.tenv
  ];
```
#### nix-shell

```sh
nix-shell -p tenv
```
</details>

<a id="manual-installation"></a>
#### Manual Installation
Get the most recent packaged binaries (`.deb`, `.rpm`, `.apk`, `pkg.tar.zst `, `.zip` or `.tar.gz` format) by visiting the [release page](https://github.com/tofuutils/tenv/releases). After downloading, unzip the folder and seamlessly integrate it into your system's `PATH`.

<a id="docker-installation"></a>

#### Docker Installation

You can use dockerized version of tenv via the following command:

```sh
docker run -it --rm tofuutils/tenv:latest help
```
The docker container is not meant as a way to run tenv for CI pipelines, for local use, you should use one of the [packaged binaries](#manual-installation).
<a id="usage"></a>

<a id="shell-completion"></a>
### Install shell completion

> [!NOTE]
> If you install tenv via Brew, MacPorts, or Nix, completion will be installed automatically.

<details markdown="1"><summary><b>zsh</b></summary><br>

```console
tenv completion zsh > ~/.tenv.completion.zsh
echo "source \$HOME/.tenv.completion.zsh" >> ~/.zshrc
```
</details>

<details markdown="1"><summary><b>Oh My Zsh</b></summary><br>

```console
tenv completion zsh > ~/.oh-my-zsh/completions/_tenv
```

Make sure the completions folder `~/.oh-my-zsh/completions` is listed under `$fpath`:

```console
print -l $fpath
```
</details>

<details markdown="1"><summary><b>powershell</b></summary><br>

```console
tenv completion powershell | Out-String | Invoke-Expression
```
</details>

<details markdown="1"><summary><b>bash</b></summary><br>

```console
tenv completion bash > ~/.tenv.completion.bash
echo "source \$HOME/.tenv.completion.bash" >> ~/.bashrc
```
</details>

<details markdown="1"><summary><b>fish</b></summary><br>

```console
tenv completion fish > ~/.tenv.completion.fish
echo "source \$HOME/.tenv.completion.fish" >> ~/.config/fish/config.fish
```
</details>

## Usage

**tenv** supports [OpenTofu](https://opentofu.org),
[Terraform](https://www.terraform.io/), [Terragrunt](https://terragrunt.gruntwork.io/), [Terramate](https://terramate.io/) and
[Atmos](https://atmos.tools). To manage each binary you can use `tenv <tool> <command>`. Below is a list of tools and commands that use actual subcommands:

| tool (alias)        | env vars                   | description                                    |
| ------------------- | -------------------------- | ---------------------------------------------- |
| `tofu` (`opentofu`) | [TOFUENV_](#tofu-env-vars) | [OpenTofu](https://opentofu.org)               |
| `tf` (`terraform`)  | [TFENV_](#tf-env-vars)     | [Terraform](https://www.terraform.io/)         |
| `tg` (`terragrunt`) | [TG_](#tg-env-vars)        | [Terragrunt](https://terragrunt.gruntwork.io/) |
| `tm` (`terramate`)  | [TM_](#tm-env-vars)        | [Terramate](https://terramate.io/)             |
| `at` (`atmos`)      | [ATMOS_](#atmos-env-vars)  | [Atmos](https://atmos.tools)                   |


Without subcommand `tenv` display interactive menus to manage tools and their versions.
<p align="center">
  <a href="https://asciinema.org/a/670790">
    <img alt="tenv interactive" src="https://raw.githubusercontent.com/tofuutils/tenv/main/assets/tenv.gif" width="100%">
  </a>
</p>

<details markdown="1"><summary><b>tenv &lt;tool&gt; install [version]</b></summary><br>

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
tenv atmos install "~> 1.70"
tenv atmos install latest
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


<details markdown="1"><summary><b>tenv &lt;tool&gt; use  &lt;version&gt;</b></summary><br>

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
tenv atmos use latest
tenv tofu use latest-allowed
```

</details>


<details markdown="1"><summary><b>tenv &lt;tool&gt; detect</b></summary><br>

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
$ tenv atmos detect -q
Atmos 1.72.0 will be run from this directory.
```

</details>


<details markdown="1"><summary><b>tenv &lt;tool&gt; reset</b></summary><br>

Reset used version of tool (remove `TENV_ROOT/<TOOL>/version` file).

```console
$ tenv tofu reset
Removed /home/dvaumoron/.tenv/OpenTofu/version
```

</details>


<details markdown="1"><summary><b>tenv &lt;tool&gt; uninstall [version]</b></summary><br>

Uninstall versions of the tool (remove it from `TENV_ROOT` directory).

Without parameter, display an interactive list to select several versions.

If a parameter is passed, available parameter options:

- an exact [Semver 2.0.0](https://semver.org/) version string to remove (no confirmation required)
- a [version constraint](https://opentofu.org/docs/language/expressions/version-constraints) string
- `all`
- `but-last` (all versions except the highest installed)
- `not-used-for:<duration>`, `<duration>` in days or months, like "14d" or "2m"
- `not-used-since:<date>`, `<date>` format is YYYY-MM-DD, like "2024-06-30"

```console
$ tenv tofu uninstall v1.6.0-alpha4
Uninstallation of OpenTofu 1.6.0-alpha4 successful (directory /home/dvaumoron/.tenv/OpenTofu/1.6.0-alpha4 removed)
```

</details>


<details markdown="1"><summary><b>tenv &lt;tool&gt; list</b></summary><br>

List installed tool versions (located in `TENV_ROOT` directory), sorted in ascending version order.

`tenv <tool> list` has a `--descending`, `-d` flag to sort in descending order.

```console
$ tenv tofu list -v
* 1.6.0 (set by /home/dvaumoron/.tenv/OpenTofu/version)
  1.6.1
found 2 OpenTofu version(s) managed by tenv.
```

</details>


<details markdown="1"><summary><b>tenv &lt;tool&gt; list-remote</b></summary><br>

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


<details markdown="1"><summary><b>tenv &lt;tool&gt; constraint [expression]</b></summary><br>

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


<details markdown="1"><summary><b>tenv help [command]</b></summary><br>

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


<details markdown="1"><summary><b>tenv update-path</b></summary><br>

Display PATH updated with tenv directory location first. With GITHUB_ACTIONS set to true, write tenv directory location to GITHUB_PATH.

This command can be used when one of the managed tool is already installed on your system and hide the corresponding proxy (in that case `which tenv` and `which <tool>` will indicate different locations). The following shell call should resolve such issues :

```console
export PATH=$(tenv update-path)
```

</details>


<details markdown="1"><summary><b>tenv version</b></summary><br>

Display tenv current version.

```console
$ tenv version
tenv version v1.7.0
```


</details>


<a id="environment-variables"></a>
## Environment variables

**tenv** commands support global environment variables and variables by tool for : [OpenTofu](https://opentofu.org), [Terraform](https://www.terraform.io/), [TerraGrunt](https://terragrunt.gruntwork.io/) and [Atmos](https://atmos.tools).


<a id="tenv-vars"></a>
### Global tenv environment variables

<details markdown="1"><summary><b>TENV_ARCH</b></summary><br>

String (Default: current tenv binaries architecture)

Allow to override the default architecture for binaries downloading during installation.

`tenv <tool>` subcommands `detect`, `install` and `use` support a `--arch`, `-a`  flag version.

</details>


<details markdown="1"><summary><b>TENV_AUTO_INSTALL</b></summary><br>

String (Default: false)

If set to true **tenv** will automatically install missing tool versions needed.

`tenv <tool>` subcommands `detect` and `use` support a `--install`, `-i` enabling flag, and a `--no-install`, `-n` disabling flag.

</details>


<details markdown="1"><summary><b>TENV_FORCE_REMOTE</b></summary><br>

String (Default: false)

If set to true **tenv** detection of needed version will skip local check and verify compatibility on remote list.

`tenv <tool>` subcommands `detect` and `use` support a `--force-remote`, `-f` flag version.

</details>


<details markdown="1"><summary><b>TENV_GITHUB_TOKEN</b></summary><br>

String (Default: "")

Allow to specify a GitHub token to increase [GitHub Rate limits for the REST API](https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api). Useful because OpenTofu, Terragrunt and Atmos binaries are downloaded from GitHub repository.

`tenv tofu` and `tenv tg` subcommands `detect`, `install`, `list-remote` and `use` support a `--github-token`, `-t` flag version.

</details>


<details markdown="1"><summary><b>TENV_QUIET</b></summary><br>

String (Default: false)

If set to true **tenv** disable unnecessary output (including log level forced to off).

`tenv` subcommands support a `--quiet`, `-q` flag version.

</details>


<details markdown="1"><summary><b>TENV_LOG</b></summary><br>

String (Default: "warn")

Set **tenv** log level (possibilities sorted by decreasing verbosity : "trace", "debug", "info", "warn", "error", "off").

`tenv` support a `--verbose`, `-v` flag which set log level to "trace".

</details>


<details markdown="1"><summary><b>TENV_REMOTE_CONF</b></summary><br>

String (Default: `${TENV_ROOT}/remote.yaml`)

The path to a yaml file for [advanced remote configuration](#advanced-remote-configuration) (can be used to call artifact mirror).

`tenv <tool>` subcommands `detect`, `install`, `list-remote` and `use`  support a `--remote-conf`, `-c` flag version.

</details>


<details markdown="1"><summary><b>TENV_ROOT</b></summary><br>

String (Default: `${HOME}/.tenv`)

The path to a directory where the local OpenTofu versions, Terraform versions, Terragrunt versions and tenv configuration files exist.

`tenv` support a `--root-path`, `-r` flag version.

</details>


<details markdown="1"><summary><b>TENV_SKIP_LAST_USE</b></summary><br>

String (Default: false)

If set to true **tenv** disable tracking of last use date for installed versions. It allow to avoid warning message when **tenv** is installed as root user and run with a normal user by skipping the writing of `last-use.txt`. This will lead to misselection with `not-used-for` and `not-used-since` behavior of `tenv uninstall`.

</details>


<details markdown="1"><summary><b>TENV_VALIDATION</b></summary><br>

String (Default: signature)

Set **tenv** validation, known values are "signature" (check SHA256 and its signature, see [signature support](#signature-support)), "sha" (only check SHA256), "none" (no validation).

</details>


<details markdown="1"><summary><b>GITHUB_ACTIONS</b></summary><br>

String (Default: false)

If set to true **tenv** proxies exposes proxied output `stdout`, `stderr`, and `exitcode` by writing them into GITHUB_OUTPUT in [multiline format](https://docs.github.com/en/actions/using-workflows/workflow-commands-for-github-actions#multiline-strings). GitHub Actions set it (see [default environment variables](https://docs.github.com/en/actions/learn-github-actions/variables#default-environment-variables)).

</details>


<details markdown="1"><summary><b>GITHUB_OUTPUT</b></summary><br>

String (Default: "")

Needed when GITHUB_ACTIONS is set to true, path to a file to write proxied output.

</details>


<details markdown="1"><summary><b>GITHUB_PATH</b></summary><br>

String (Default: "")

Used by `tenv update-path` when GITHUB_ACTIONS is set to true, path to a file to write tenv directory location.

</details>


<a id="tofu-env-vars"></a>
### OpenTofu environment variables

<details markdown="1"><summary><b>TOFUENV_AGNOSTIC_PROXY</b></summary><br>

String (Default: false)

Switch `tofu` proxy to an agnostic proxy (behave like `tf`, see [resolution order](#project-binaries)).

</details>


<details markdown="1"><summary><b>TOFUENV_ARCH</b></summary><br>

Same as TENV_ARCH (compatibility with [tofuenv](https://github.com/tofuutils/tofuenv)).

</details>


<details markdown="1"><summary><b>TOFUENV_AUTO_INSTALL</b></summary><br>

Same as TENV_AUTO_INSTALL (compatibility with [tofuenv](https://github.com/tofuutils/tofuenv)).

#### Example 1
Use OpenTofu version 1.6.1 that is not installed, and auto installation stay disabled :

```console
$ tenv use 1.6.1
Written 1.6.1 in /home/dvaumoron/.tenv/OpenTofu/version
```

#### Example 2
Use OpenTofu version 1.6.0 that is not installed, and auto installation is enabled :

```console
$ TOFUENV_AUTO_INSTALL=true tenv tofu use 1.6.0
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


<details markdown="1"><summary><b>TOFUENV_FORCE_REMOTE</b></summary><br>

Same as TENV_FORCE_REMOTE.

</details>


<details markdown="1"><summary><b>TOFUENV_INSTALL_MODE</b></summary><br>

String (the default depend on TOFUENV_REMOTE, without change on it, it is "api" else it is "direct")

- "api" install mode retrieve download url of OpenTofu from [Github REST API](https://docs.github.com/en/rest?apiVersion=2022-11-28) (TOFUENV_REMOTE must comply with it).
- "direct" install mode generate download url of OpenTofu based on TOFUENV_REMOTE.
- "mirror" install mode generate download url with TOFUENV_URL_TEMPLATE (as specified in [TofuDL mirror specification](https://github.com/opentofu/tofudl/blob/mirror-spec/MIRROR-SPECIFICATION.md))

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details markdown="1"><summary><b>TOFUENV_LIST_MODE</b></summary><br>

String (the default depend on TOFUENV_LIST_URL, without change on it, it is "api" else it is "html")

- "api" list mode retrieve information of OpenTofu releases from [Github REST API](https://docs.github.com/en/rest?apiVersion=2022-11-28) (TOFUENV_LIST_URL must comply with it).
- "html" list mode extract information of OpenTofu releases from parsing an html page at TOFUENV_LIST_URL.
- "mirror" list mode retrieve information of OpenTofu releases at TOFUENV_LIST_URL as [TofuDL Mirroring format](https://get.opentofu.org/tofu/api.json)

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details markdown="1"><summary><b>TOFUENV_LIST_URL</b></summary><br>

String (Default: copy TOFUENV_REMOTE, default is overloaded by "https://get.opentofu.org/tofu/api.json" when TOFUENV_LIST_MODE is "mirror")

Allow to override the remote url only for the releases listing.

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details markdown="1"><summary><b>TOFUENV_OPENTOFU_PGP_KEY</b></summary><br>

String (Default: "")

Allow to specify a local file path or URL to OpenTofu PGP public key. If a URL is provided (starting with "http://" or "https://"), the key will be downloaded from that URL. If a local file path is provided, the key will be read from that location. If not set, the key will be downloaded from the default URL (https://get.opentofu.org/opentofu.asc).

`tenv tofu` subcommands `detect`, `ìnstall` and `use` support a `--key-file`, `-k` flag version.

</details>


<details markdown="1"><summary><b>TOFUENV_REMOTE</b></summary><br>

String (Default: https://api.github.com/repos/opentofu/opentofu/releases)

URL to install OpenTofu, when TOFUENV_REMOTE differ from its default value, TOFUENV_INSTALL_MODE is set to "direct" and TOFUENV_LIST_MODE is set to "html" (assume an artifact proxy usage).

`tenv tofu` subcommands `detect`, `install`, `list-remote` and `use` support a `--remote-url`, `-u` flag version.

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details markdown="1"><summary><b>TOFUENV_REMOTE_PASSWORD</b></summary><br>

String (Default: "")

Could be used with TOFUENV_REMOTE_USER to specify HTTP basic auth when same credential are used with TOFUENV_REMOTE and TOFUENV_LIST_URL (instead of `https://user:password@host.org` URL format).

</details>


<details markdown="1"><summary><b>TOFUENV_REMOTE_USER</b></summary><br>

String (Default: "")

Could be used with TOFUENV_REMOTE_PASSWORD to specify HTTP basic auth when same credential are used with TOFUENV_REMOTE and TOFUENV_LIST_URL (instead of `https://user:password@host.org` URL format).

</details>


<details markdown="1"><summary><b>TOFUENV_ROOT</b></summary><br>

Same as TENV_ROOT (compatibility with [tofuenv](https://github.com/tofuutils/tofuenv)).

</details>


<details markdown="1"><summary><b>TOFUENV_URL_TEMPLATE</b></summary><br>

String (Default: `https://github.com/opentofu/opentofu/releases/download/v{{ .Version }}/{{ .Artifact }}`)

Used when TOFUENV_INSTALL_MODE is "mirror" (see [TofuDL mirror specification](https://github.com/opentofu/tofudl/blob/mirror-spec/MIRROR-SPECIFICATION.md)).

</details>


<details markdown="1"><summary><b>TOFUENV_GITHUB_TOKEN</b></summary><br>

Same as TENV_GITHUB_TOKEN (compatibility with [tofuenv](https://github.com/tofuutils/tofuenv)).

</details>


<details markdown="1"><summary><b>TOFUENV_TOFU_DEFAULT_CONSTRAINT</b></summary><br>

String (Default: "")

If not empty string, this variable overrides OpenTofu default constraint, specified in ${TENV_ROOT}/OpenTofu/constraint file.

</details>


<details markdown="1"><summary><b>TOFUENV_TOFU_DEFAULT_VERSION</b></summary><br>

String (Default: "")

If not empty string, this variable overrides OpenTofu fallback version, specified in ${TENV_ROOT}/OpenTofu/version file.

</details>


<details markdown="1"><summary><b>TOFUENV_TOFU_VERSION</b></summary><br>

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

<details markdown="1"><summary><b>TFENV_AGNOSTIC_PROXY</b></summary><br>

String (Default: false)

Switch `terraform` proxy to an agnostic proxy (behave like `tf`, see [resolution order](#project-binaries)).

</details>


<details markdown="1"><summary><b>TFENV_ARCH</b></summary><br>

Same as TENV_ARCH (compatibility with [tfenv](https://github.com/tfutils/tfenv)).

</details>


<details markdown="1"><summary><b>TFENV_AUTO_INSTALL</b></summary><br>

Same as TENV_AUTO_INSTALL (compatibility with [tfenv](https://github.com/tfutils/tfenv)).

</details>


<details markdown="1"><summary><b>TFENV_FORCE_REMOTE</b></summary><br>

Same as TENV_FORCE_REMOTE.

</details>


<details markdown="1"><summary><b>TFENV_HASHICORP_PGP_KEY</b></summary><br>

Allow to specify a local file path or URL to Hashicorp PGP public key. If a URL is provided (starting with "http://" or "https://"), the key will be downloaded from that URL. If a local file path is provided, the key will be read from that location. If not set, the key will be downloaded from the default URL (https://www.hashicorp.com/.well-known/pgp-key.txt).

`tenv tf` subcommands `detect`, `ìnstall` and `use` support a `--key-file`, `-k` flag version.
</details>


<details markdown="1"><summary><b>TFENV_INSTALL_MODE</b></summary><br>

String (Default: "api")

- "api" install mode retrieve download url of Terraform from [Hashicorp Release API](https://releases.hashicorp.com/docs/api/v1) (TFENV_REMOTE must comply with it).
- "direct" install mode generate download url of Terraform based on TFENV_REMOTE.

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details markdown="1"><summary><b>TFENV_LIST_MODE</b></summary><br>

String (the default depend on TFENV_LIST_URL, without change on it, it is "api" else it is "html")

- "api" list mode retrieve information of Terraform releases from [Hashicorp Release API](https://releases.hashicorp.com/docs/api/v1) (TFENV_LIST_URL must comply with it).
- "html" list mode extract information of Terraform releases from parsing an html page in TFENV_LIST_URL.

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details markdown="1"><summary><b>TFENV_LIST_URL</b></summary><br>

String (Default: copy TFENV_REMOTE)

Allow to override the remote url only for the releases listing.

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details markdown="1"><summary><b>TFENV_REMOTE</b></summary><br>

String (Default: https://releases.hashicorp.com)

URL to install Terraform, changing it assume an artifact proxy use (TFENV_LIST_URL copy it, and if it differ from its default value, TFENV_LIST_MODE is set to "html", because an artifact proxy usage will not disturb the retrieving of index.json for a release, but will freeze the json list of releases).

`tenv tf` subcommands `detect`, `install`, `list-remote` and `use` support a `--remote-url`, `-u` flag version.

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details markdown="1"><summary><b>TFENV_REMOTE_PASSWORD</b></summary><br>

String (Default: "")

Could be used with TFENV_REMOTE_USER to specify HTTP basic auth when same credential are used with TFENV_REMOTE and TFENV_LIST_URL (instead of `https://user:password@host.org` URL format).

</details>


<details markdown="1"><summary><b>TFENV_REMOTE_USER</b></summary><br>

String (Default: "")

Could be used with TFENV_REMOTE_PASSWORD to specify HTTP basic auth when same credential are used with TFENV_REMOTE and TFENV_LIST_URL (instead of `https://user:password@host.org` URL format).

</details>


<details markdown="1"><summary><b>TFENV_ROOT</b></summary><br>

Same as TENV_ROOT (compatibility with [tfenv](https://github.com/tfutils/tfenv)).

</details>


<details markdown="1"><summary><b>TFENV_TERRAFORM_DEFAULT_CONSTRAINT</b></summary><br>

String (Default: "")

If not empty string, this variable overrides Terraform default constraint, specified in ${TENV_ROOT}/Terraform/constraint file.

</details>


<details markdown="1"><summary><b>TFENV_TERRAFORM_DEFAULT_VERSION</b></summary><br>

String (Default: "")

If not empty string, this variable overrides Terraform fallback version, specified in ${TENV_ROOT}/Terraform/version file.

</details>


<details markdown="1"><summary><b>TFENV_TERRAFORM_VERSION</b></summary><br>

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


<details markdown="1"><summary><b>TG_INSTALL_MODE</b></summary><br>

String (the default depend on TG_REMOTE, without change on it, it is "api" else it is "direct")

- "api" install mode retrieve download url of Terragrunt from [Github REST API](https://docs.github.com/en/rest?apiVersion=2022-11-28) (TG_REMOTE must comply with it).
- "direct" install mode generate download url of Terragrunt based on TG_REMOTE.

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details markdown="1"><summary><b>TG_LIST_MODE</b></summary><br>

String (the default depend on TG_LIST_URL, without change on it, it is "api" else it is "html")

- "api" list mode retrieve information of Terragrunt releases from [Github REST API](https://docs.github.com/en/rest?apiVersion=2022-11-28) (TG_LIST_URL must comply with it).
- "html" list mode extract information of Terragrunt releases from parsing an html page in TG_LIST_URL.

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details markdown="1"><summary><b>TG_LIST_URL</b></summary><br>

String (Default: copy TG_REMOTE)

Allow to override the remote url only for the releases listing.

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details markdown="1"><summary><b>TG_REMOTE</b></summary><br>

String (Default: https://api.github.com/repos/gruntwork-io/terragrunt/releases)

URL to install Terragrunt, when TG_REMOTE differ from its default value, TG_INSTALL_MODE is set to "direct" and TG_LIST_MODE is set to "html" (assume an artifact proxy usage).

`tenv tg` subcommands `detect`, `install`, `list-remote` and `use` support a `--remote-url`, `-u` flag version.

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details markdown="1"><summary><b>TG_REMOTE_PASSWORD</b></summary><br>

String (Default: "")

Could be used with TG_REMOTE_USER to specify HTTP basic auth when same credential are used with TG_REMOTE and TG_LIST_URL (instead of `https://user:password@host.org` URL format).

</details>


<details markdown="1"><summary><b>TG_REMOTE_USER</b></summary><br>

String (Default: "")

Could be used with TG_REMOTE_PASSWORD to specify HTTP basic auth when same credential are used with TG_REMOTE and TG_LIST_URL (instead of `https://user:password@host.org` URL format).

</details>


<details markdown="1"><summary><b>TG_DEFAULT_CONSTRAINT</b></summary><br>

String (Default: "")

If not empty string, this variable overrides Terragrunt default constraint, specified in ${TENV_ROOT}/Terragrunt/constraint file.

</details>


<details markdown="1"><summary><b>TG_DEFAULT_VERSION</b></summary><br>

String (Default: "")

If not empty string, this variable overrides Terragrunt fallback version, specified in ${TENV_ROOT}/Terragrunt/version file.

</details>


<details markdown="1"><summary><b>TG_VERSION</b></summary><br>

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


<a id="tm-env-vars"></a>
### Terramate environment variables


<details markdown="1"><summary><b>TM_INSTALL_MODE</b></summary><br>

String (the default depend on TM_REMOTE, without change on it, it is "api" else it is "direct")

- "api" install mode retrieve download url of Terramate from [Github REST API](https://docs.github.com/en/rest?apiVersion=2022-11-28) (TM_REMOTE must comply with it).
- "direct" install mode generate download url of Terramate based on TM_REMOTE.

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details markdown="1"><summary><b>TM_LIST_MODE</b></summary><br>

String (the default depend on TM_LIST_URL, without change on it, it is "api" else it is "html")

- "api" list mode retrieve information of Terramate releases from [Github REST API](https://docs.github.com/en/rest?apiVersion=2022-11-28) (TM_LIST_URL must comply with it).
- "html" list mode extract information of Terramate releases from parsing an html page in TM_LIST_URL.

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details markdown="1"><summary><b>TM_LIST_URL</b></summary><br>

String (Default: copy TM_REMOTE)

Allow to override the remote url only for the releases listing.

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details markdown="1"><summary><b>TM_REMOTE</b></summary><br>

String (Default: https://api.github.com/repos/terramate-io/terramate/releases)

URL to install Terramate, when TM_REMOTE differ from its default value, TM_INSTALL_MODE is set to "direct" and TM_LIST_MODE is set to "html" (assume an artifact proxy usage).

`tenv tm` subcommands `detect`, `install`, `list-remote` and `use` support a `--remote-url`, `-u` flag version.

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details markdown="1"><summary><b>TM_REMOTE_PASSWORD</b></summary><br>

String (Default: "")

Could be used with TM_REMOTE_USER to specify HTTP basic auth when same credential are used with TM_REMOTE and TM_LIST_URL (instead of `https://user:password@host.org` URL format).

</details>


<details markdown="1"><summary><b>TM_REMOTE_USER</b></summary><br>

String (Default: "")

Could be used with TM_REMOTE_PASSWORD to specify HTTP basic auth when same credential are used with TM_REMOTE and TM_LIST_URL (instead of `https://user:password@host.org` URL format).

</details>


<details markdown="1"><summary><b>TM_DEFAULT_CONSTRAINT</b></summary><br>

String (Default: "")

If not empty string, this variable overrides Terramate default constraint, specified in ${TENV_ROOT}/Terramate/constraint file.

</details>


<details markdown="1"><summary><b>TM_DEFAULT_VERSION</b></summary><br>

String (Default: "")

If not empty string, this variable overrides Terramate fallback version, specified in ${TENV_ROOT}/Terramate/version file.

</details>


<details markdown="1"><summary><b>TM_VERSION</b></summary><br>

String (Default: "")

If not empty string, this variable overrides Terramate version, specified in [`.terramate-version`](#terramate-version-files) files.

`tenv tm` subcommands `install` and `detect` also respects this variable.

e.g. with :

```console
$ terramate version
0.13.0
```

then :

```console
$ TM_VERSION=0.12.0 terramate version
0.12.0
```

</details>


<a id="atmos-env-vars"></a>
### Atmos environment variables


<details markdown="1"><summary><b>ATMOS_INSTALL_MODE</b></summary><br>

String (the default depend on ATMOS_REMOTE, without change on it, it is "api" else it is "direct")

- "api" install mode retrieve download url of Atmos from [Github REST API](https://docs.github.com/en/rest?apiVersion=2022-11-28) (ATMOS_REMOTE must comply with it).
- "direct" install mode generate download url of Atmos based on ATMOS_REMOTE.

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details markdown="1"><summary><b>ATMOS_LIST_MODE</b></summary><br>

String (the default depend on ATMOS_LIST_URL, without change on it, it is "api" else it is "html")

- "api" list mode retrieve information of Atmos releases from [Github REST API](https://docs.github.com/en/rest?apiVersion=2022-11-28) (ATMOS_LIST_URL must comply with it).
- "html" list mode extract information of Atmos releases from parsing an html page in ATMOS_LIST_URL.

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details markdown="1"><summary><b>ATMOS_LIST_URL</b></summary><br>

String (Default: copy ATMOS_REMOTE)

Allow to override the remote url only for the releases listing.

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details markdown="1"><summary><b>ATMOS_REMOTE</b></summary><br>

String (Default: https://api.github.com/repos/cloudposse/atmos/releases)

URL to install Atmos when ATMOS_REMOTE differ from its default value, ATMOS_INSTALL_MODE is set to "direct" and ATMOS_LIST_MODE is set to "html" (assume an artifact proxy usage).

`tenv atmos` subcommands `detect`, `install`, `list-remote` and `use` support a `--remote-url`, `-u` flag version.

See [advanced remote configuration](#advanced-remote-configuration) for more details.

</details>


<details markdown="1"><summary><b>ATMOS_REMOTE_PASSWORD</b></summary><br>

String (Default: "")

Could be used with ATMOS_REMOTE_USER to specify HTTP basic auth when same credential are used with ATMOS_REMOTE and ATMOS_LIST_URL (instead of `https://user:password@host.org` URL format).

</details>


<details markdown="1"><summary><b>ATMOS_REMOTE_USER</b></summary><br>

String (Default: "")

Could be used with ATMOS_REMOTE_PASSWORD to specify HTTP basic auth when same credential are used with ATMOS_REMOTE and ATMOS_LIST_URL (instead of `https://user:password@host.org` URL format).

</details>


<details markdown="1"><summary><b>ATMOS_DEFAULT_CONSTRAINT</b></summary><br>

String (Default: "")

If not empty string, this variable overrides Atmos default constraint, specified in ${TENV_ROOT}/Atmos/constraint file.

</details>


<details markdown="1"><summary><b>ATMOS_DEFAULT_VERSION</b></summary><br>

String (Default: "")

If not empty string, this variable overrides Atmos fallback version, specified in ${TENV_ROOT}/Atmos/version file.

</details>


<details markdown="1"><summary><b>ATMOS_VERSION</b></summary><br>

String (Default: "")

If not empty string, this variable overrides Atmos version, specified in [`.atmos-version`](#atmos-version-files) files.

`tenv atmos` subcommands `install` and `detect` also respects this variable.

e.g. with :

```console
$ atmos version
👽 Atmos v1.72.0 on linux/amd64

```

then :

```console
$ ATMOS_VERSION=1.70 atmos version
👽 Atmos v1.70.0 on linux/amd64
```

</details>

<a id="version-files"></a>
## version files

<a id="default-version-file"></a>
<details markdown="1"><summary><b>default version file</b></summary><br>

The `TENV_ROOT/<TOOL>/version` file is the tool default version used when no project specific or user specific are found. It can be written with `tenv <tool> use`.

</details>

<a id="opentofu-version-files"></a>
<details markdown="1"><summary><b>opentofu version files</b></summary><br>

If you put a `.opentofu-version` file in the working directory, one of its parent directory, or user home directory, **tenv** detects it and uses the version written in it.
Note, that TOFUENV_TOFU_VERSION can be used to override version specified by `.opentofu-version` file.

Recognize same values as `tenv tofu use` command.

See [required_version](https://opentofu.org/docs/language/settings#specifying-a-required-opentofu-version) docs.

</details>

<a id="terraform-version-files"></a>
<details markdown="1"><summary><b>terraform version files</b></summary><br>

If you put a `.terraform-version` or `.tfswitchrc` file in the working directory, one of its parent directory, or user home directory, **tenv** detects it and uses the version written in it.
Note, that TFENV_TERRAFORM_VERSION can be used to override version specified by those files.

Recognize same values as `tenv tf use` command.

See [required_version](https://developer.hashicorp.com/terraform/language/settings#specifying-a-required-terraform-version) docs.

</details>

<a id="terragrunt-version-files"></a>
<details markdown="1"><summary><b>terragrunt version files</b></summary><br>

If you put a `.terragrunt-version` or a `.tgswitchrc` file in the working directory, one of its parent directory, or user home directory, **tenv** detects it and uses the version written in it. **tenv** also detect a `version` field in a `.tgswitch.toml` in same places.
Note, that TG_VERSION can be used to override version specified by those files.

Recognize same values as `tenv tg use` command.

</details>


<a id="terragrunt-hcl-file"></a>
<details markdown="1"><summary><b>terragrunt.hcl or root.hcl file</b></summary><br>

[Terragrunt now recommends](https://terragrunt.gruntwork.io/docs/migrate/migrating-from-root-terragrunt-hcl/) using `root.hcl` instead of `terragrunt.hcl` as the root configuration file name.

If a `terragrunt.hcl`, `root.hcl`, or their `.json` equivalents exist in the working directory, a parent directory, or the user home directory, **tenv** will read constraints from the `terraform_version_constraint` or `terragrunt_version_constraint` field (depending on proxy or subcommand used).


If both `root.hcl` and `terragrunt.hcl` (or their `.json` versions) are present, `terragrunt.hcl` takes precedence.

</details>

<a id="terramate-version-files"></a>
<details markdown="1"><summary><b>terramate version files</b></summary><br>

If you put a `.terramate-version` file in the working directory, one of its parent directory, or user home directory, **tenv** detects it and uses the version written in it.
Note, that TM_VERSION can be used to override version specified by those files.

Recognize same values as `tenv tm use` command.

</details>

<a id="atmos-version-files"></a>
<details markdown="1"><summary><b>atmos version files</b></summary><br>

If you put a `.atmos-version` file in the working directory, one of its parent directory, or user home directory, **tenv** detects it and uses the version written in it.
Note, that ATMOS_VERSION can be used to override version specified by those files.

Recognize same values as `tenv atmos use` command.

</details>

<a id="required_version"></a>
<details markdown="1"><summary><b>required_version</b></summary><br>

the `latest-allowed` or `min-required` strategies scan through your IAC files (see list in [project binaries](#project-binaries)) and identify a version conforming to the constraint in the relevant files. They fallback to `latest` when no IAC files and no default constraint are found, and can optionally be used with a default constraint as detailed in [project binaries](#project-binaries).

Currently the format for [Terraform required_version](https://developer.hashicorp.com/terraform/language/settings#specifying-a-required-terraform-version) and [OpenTofu required_version](https://opentofu.org/docs/language/settings#specifying-a-required-opentofu-version) are very similar, however this may change over time, always refer to docs for the latest format specification.

example:

```HCL
version = ">= 1.2.0, < 2.0.0"
```

This would identify the latest version at or above 1.2.0 and below 2.0.0

</details>

<a id="technical-details"></a>
## Technical details

<a id="project-binaries"></a>
### Project binaries

All the proxy binaries return the exit code `42` on error happening before proxied command call.


<details markdown="1"><summary><b>tofu</b></summary><br>

The `tofu` command in this project is a proxy to OpenTofu's `tofu` command  managed by **tenv**.

The version resolution order is :

- TOFUENV_TOFU_VERSION environment variable
- `.opentofu-version` file
- `.tool-versions` [file](https://asdf-vm.com/manage/configuration.html#tool-versions)
- `terraform_version_constraint` from `terragrunt.hcl` file
- `terraform_version_constraint` from `terragrunt.hcl.json` file
- `terraform_version_constraint` from `root.hcl` file
- `terraform_version_constraint` from `root.hcl.json` file
- TOFUENV_TOFU_DEFAULT_VERSION environment variable
- `${TENV_ROOT}/OpenTofu/version` file (can be written with `tenv tofu use`)
- `latest-allowed`

The `latest-allowed` strategy rely on [required_version](#required_version) from .tofu, .tofu.json, .tf or .tf.json files with a fallback to `latest` when no constraint are found. Moreover it is possible to add a default constraint with TOFUENV_TOFU_DEFAULT_CONSTRAINT environment variable or `${TENV_ROOT}/OpenTofu/constraint` file (can be written with `tenv tofu constraint`). The default constraint is added while using `latest-allowed`, `min-required` or custom constraint. A default constraint with `latest-allowed` or `min-required` will avoid the fallback to `latest` when there is no .tf or .tf.json files.

</details>

<details markdown="1"><summary><b>terraform</b></summary><br>

The `terraform` command in this project is a proxy to HashiCorp's `terraform` command managed by **tenv**.

The version resolution order is :

- TFENV_TERRAFORM_VERSION environment variable
- `.terraform-version` file
- `.tfswitchrc` file
- `.tool-versions` [file](https://asdf-vm.com/manage/configuration.html#tool-versions)
- `terraform_version_constraint` from `terragrunt.hcl` file
- `terraform_version_constraint` from `terragrunt.hcl.json` file
- `terraform_version_constraint` from `root.hcl` file
- `terraform_version_constraint` from `root.hcl.json` file
- TFENV_TERRAFORM_DEFAULT_VERSION environment variable
- `${TENV_ROOT}/Terraform/version` file (can be written with `tenv tf use`)
- `latest-allowed`

The `latest-allowed` strategy rely on [required_version](#required_version) from .tf or .tf.json files with a fallback to `latest` when no constraint are found. Moreover it is possible to add a default constraint with TFENV_TERRAFORM_DEFAULT_CONSTRAINT environment variable or `${TENV_ROOT}/Terraform/constraint` file (can be written with `tenv tf constraint`). The default constraint is added while using `latest-allowed`, `min-required` or custom constraint. A default constraint with `latest-allowed` or `min-required` will avoid the fallback to `latest` when there is no .tf or .tf.json files.

</details>


<details markdown="1"><summary><b>terragrunt</b></summary><br>

The `terragrunt` command in this project is a proxy to Gruntwork's `terragrunt` command managed by **tenv**.

The version resolution order is :

- TG_VERSION environment variable
- `.terragrunt-version` file
- `.tgswitchrc` file
- `version` from `tgswitch.toml` file
- `.tool-versions` [file](https://asdf-vm.com/manage/configuration.html#tool-versions)
- `terragrunt_version_constraint` from `terragrunt.hcl` file
- `terragrunt_version_constraint` from `terragrunt.hcl.json` file
- `terragrunt_version_constraint` from `root.hcl` file
- `terragrunt_version_constraint` from `root.hcl.json` file
- TG_DEFAULT_VERSION environment variable
- `${TENV_ROOT}/Terragrunt/version` file (can be written with `tenv tg use`)
- `latest-allowed`

The `latest-allowed` strategy has no information for Terragrunt and will fallback to `latest` unless there is default constraint. Adding a default constraint could be done with TG_DEFAULT_CONSTRAINT environment variable or `${TENV_ROOT}/Terragrunt/constraint` file (can be written with `tenv tg constraint`). The default constraint is added while using `latest-allowed`, `min-required` or custom constraint. A default constraint with `latest-allowed` or `min-required` will avoid there fallback to `latest`.

</details>


<details markdown="1"><summary><b>terramate</b></summary><br>

The `terramate` command in this project is a proxy to Terramate's `terramate` command managed by **tenv**.

The version resolution order is :

- TM_VERSION environment variable
- `.terramate-version` file
- TM_DEFAULT_VERSION environment variable
- `${TENV_ROOT}/Terramate/version` file (can be written with `tenv tm use`)
- `latest-allowed`

The `latest-allowed` strategy has no information for Terramate and will fallback to `latest` unless there is default constraint. Adding a default constraint could be done with TM_DEFAULT_CONSTRAINT environment variable or `${TENV_ROOT}/Terramate/constraint` file (can be written with `tenv tm constraint`). The default constraint is added while using `latest-allowed`, `min-required` or custom constraint. A default constraint with `latest-allowed` or `min-required` will avoid there fallback to `latest`.

</details>


<details markdown="1"><summary><b>atmos</b></summary><br>

The `atmos` command in this project is a proxy to Cloudposse's `atmos` command managed by **tenv**.

The version resolution order is :

- ATMOS_VERSION environment variable
- `.atmos-version` file
- `.tool-versions` [file](https://asdf-vm.com/manage/configuration.html#tool-versions)
- ATMOS_DEFAULT_VERSION environment variable
- `${TENV_ROOT}/Atmos/version` file (can be written with `tenv atmos use`)
- `latest-allowed`

The `latest-allowed` strategy has no information for Atmos and will fallback to `latest`
unless there is default constraint. Adding a default constraint could be done with
ATMOS_DEFAULT_CONSTRAINT environment variable or `${TENV_ROOT}/Atmos/constraint` file (can
be written with `tenv atmos constraint`). The default constraint is added while using `latest-allowed`, `min-required` or custom constraint. A default constraint with `latest-allowed` or `min-required` will avoid there fallback to `latest`.

</details>

<details markdown="1"><summary><b>tf</b></summary><br>

The `tf` command is a proxy to `tofu` or `terraform` depending on the version files present in project.

The version resolution order is :

- `.opentofu-version` file (launch `tofu`)
- `tofu` version from `.tool-versions` [file](https://asdf-vm.com/manage/configuration.html#tool-versions)
- `terraform_version_constraint` from `terragrunt.hcl` file (launch `tofu`)
- `terraform_version_constraint` from `terragrunt.hcl.json` file (launch `tofu`)
- `terraform_version_constraint` from `root.hcl` file (launch `tofu`)
- `terraform_version_constraint` from `root.hcl.json` file (launch `tofu`)
- `.terraform-version` file (launch `terraform`)
- `.tfswitchrc` file  (launch `terraform`)
- `terraform` version from `.tool-versions` [file](https://asdf-vm.com/manage/configuration.html#tool-versions)
- fail with a message

</details>


<a id="advanced-remote-configuration"></a>
### Advanced remote configuration

This advanced configuration is meant to call artifact mirror (like [JFrog Artifactory](https://jfrog.com/artifactory)).

The yaml file from TENV_REMOTE_CONF path can have one part for each supported proxy : `tofu`, `terraform`, `terragrunt` and `atmos`.

<details markdown="1"><summary><b>yaml fields description</b></summary><br>

Each part can have the following string field : `install_mode`, `list_mode`, `list_url`, `url`, `new_base_url`, `old_base_url`, `selector` and `part`

With `install_mode` set to "direct", **tenv** skip the release information fetching and generate download url instead of reading them from API (overridden by `<TOOL>_INSTALL_MODE` env var).

With `list_mode` set to "html", **tenv** change the fetching of all releases information from API to parse the parent html page of artifact location, see `selector` and `part` (overridden by `<TOOL>_LIST_MODE` env var).

`url` allows to override the default remote url (overridden by flag or `<TOOL>_REMOTE` env var).

`list_url` allows to override the remote url only for the releases listing (overridden by `<TOOL>_LIST_URL` env var).

`old_base_url` and `new_base_url` are used as url rewrite rule (if an url start with the prefix, it will be changed to use the new base url).

If `old_base_url` and `new_base_url` are empty, **tenv** try to guess right behaviour based on previous fields.

`selector` is used to gather in a list all matching html node and `part` choose on which node part (attribute name or "#text" for inner text) a version will be extracted (selector default to "a" (html link) and part default to "href" (link target))

</details>


<details markdown="1"><summary><b>Examples</b></summary><br>

Those examples assume that a GitHub proxy at https://artifactory.example.com/artifactory/github have the same behavior than [JFrog Artifactory](https://jfrog.com/artifactory) :

- mirror https://github.com/opentofu/opentofu/releases/download/v1.6.0/tofu_1.6.0_linux_amd64.zip at https://artifactory.example.com/artifactory/github/opentofu/opentofu/releases/download/v1.6.0/tofu_1.6.0_linux_amd64.zip.
- have at https://artifactory.example.com/artifactory/github/opentofu/opentofu/releases/download an html page with links on existing sub folder like "v1.6.0/"

Example 1 : Retrieve Terraform binaries and list available releases from the mirror (TFENV_LIST_MODE is optional because TFENV_LIST_URL differ from its default(when TFENV_LIST_URL is not set, it copy TFENV_REMOTE)).

```console
TFENV_REMOTE=https://artifactory.example.com/artifactory/hashicorp
TFENV_LIST_MODE=html
```

Example 2 : Retrieve Terraform binaries from the mirror and list available releases from the Hashicorp releases API.

```console
TFENV_REMOTE=https://artifactory.example.com/artifactory/hashicorp
TFENV_LIST_URL=https://releases.hashicorp.com
```

Example 1 & 2, does not need install mode (by release index.json is figed in mirror without problem), however create a rewrite rule from "https://releases.hashicorp.com" to "https://artifactory.example.com/artifactory/hashicorp" to obtains correct download URLs.

Example 3 : Retrieve OpenTofu binaries and list available releases from the mirror (TOFUENV_INSTALL_MODE and TOFUENV_LIST_MODE are optional because overloading TOFUENV_REMOTE already change them).

```console
TOFUENV_REMOTE=https://artifactory.example.com/artifactory/github
TOFUENV_INSTALL_MODE=direct
TOFUENV_LIST_MODE=html
```

Example 4 : Retrieve OpenTofu binaries from the mirror and list available releases from the GitHub API (TOFUENV_INSTALL_MODE is optional because overloading TOFUENV_REMOTE already set it to "direct").

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


<a id="lockfile-support"></a>
### Lockfile support

<details markdown="1"><summary><b>Lockfile behavior</b></summary><br>

**tenv** uses lockfiles to ensure safe concurrent operations when multiple instances are run in parallel. This prevents race conditions during installation, uninstallation, and other operations that modify the local version cache.

**Lock file location and naming:**
- Default location: `${TENV_ROOT}/{tool}.lock` (e.g., `~/.tenv/OpenTofu.lock`)
- Can be customized using the `TENV_LOCK_PATH` environment variable
- Each tool (OpenTofu, Terraform, Terragrunt, Terramate, Atmos) uses its own lock file

**Parallel execution behavior:**
When multiple **tenv** instances attempt to run simultaneously:

1. The first instance creates its lock file successfully and proceeds with the operation
2. Subsequent instances fail to create the lock file (due to exclusive creation) and enter a retry loop
3. The retry mechanism waits 1 second between attempts and logs a warning message
4. Once the first instance completes and releases its lock, the next waiting instance can acquire the lock and proceed

**Example retry behavior:**
```console
$ tenv tofu install 1.6.0 & tenv tofu install 1.6.1
[1] Installing OpenTofu 1.6.0
[2] can not write .lock file, will retry: file already exists
[2] can not write .lock file, will retry: file already exists
[1] Installation of OpenTofu 1.6.0 successful
[2] Installing OpenTofu 1.6.1
[2] Installation of OpenTofu 1.6.1 successful
```

**Lock scope:**
- **Operations that modify versions** (`install`, `uninstall`) use locks to prevent conflicts
- **Read-only operations** (`list`, `list-remote`, `detect`) don't need locks
- **Batch operations** (`InstallMultiple` API) use a single lock for the entire batch
- Locks are automatically cleaned up when operations finish (success or failure)
- Interrupt signals (Ctrl+C) properly release locks to prevent deadlocks

</details>

<a id="signature-support"></a>
### Signature support

<details markdown="1"><summary><b>OpenTofu signature support</b></summary><br>

**tenv** checks the sha256 checksum and the signature of the checksum file with [cosign](https://github.com/sigstore/cosign) (if present on your machine) or PGP (via [gopenpgp](https://github.com/ProtonMail/gopenpgp)). However, unstable OpenTofu versions are signed only with cosign (in this case, if cosign is not found tenv will display a warning).

</details>

<details markdown="1"><summary><b>Terraform signature support</b></summary><br>

**tenv** checks the sha256 checksum and the PGP signature of the checksum file (via [gopenpgp](https://github.com/ProtonMail/gopenpgp), there is no cosign signature available).

</details>

<details markdown="1"><summary><b>Terragrunt signature support</b></summary><br>

**tenv** checks the sha256 checksum (there is no signature available).

</details>

<details markdown="1"><summary><b>Atmos signature support</b></summary><br>

**tenv** checks the sha256 checksum (there is no signature available).

</details>

<a id="verifying-signature"></a>
## Verifying tenv Signatures

You can use `cosign` to verify the signature of `tenv` releases. Below is an example installing the `.rpm` using `dnf` once we've verified the signatures/integrity.

> [!NOTE]
> The example below is a bash script that could be useful if you are wanting to automate installation of `tenv` in a developer environment. Adapt it to fit your specific use case.

```bash
# Get latest release
LATEST_VERSION=$(curl --silent https://api.github.com/repos/tofuutils/tenv/releases/latest | jq -r .tag_name) #v2.6.1

# Get checksum files
curl --silent -OL https://github.com/tofuutils/tenv/releases/download/${LATEST_VERSION}/tenv_${LATEST_VERSION}_checksums.txt
curl --silent -OL https://github.com/tofuutils/tenv/releases/download/${LATEST_VERSION}/tenv_${LATEST_VERSION}_checksums.txt.sig
curl --silent -OL https://github.com/tofuutils/tenv/releases/download/${LATEST_VERSION}/tenv_${LATEST_VERSION}_checksums.txt.pem

# Get RPM files
curl --silent -OL https://github.com/tofuutils/tenv/releases/download/${LATEST_VERSION}/tenv_${LATEST_VERSION}_amd64.rpm
curl --silent -OL https://github.com/tofuutils/tenv/releases/download/${LATEST_VERSION}/tenv_${LATEST_VERSION}_amd64.rpm.sig
curl --silent -OL https://github.com/tofuutils/tenv/releases/download/${LATEST_VERSION}/tenv_${LATEST_VERSION}_amd64.rpm.pem

# Verify signatures
cosign \
    verify-blob \
    --certificate-identity "https://github.com/tofuutils/tenv/.github/workflows/release.yml@refs/tags/${LATEST_VERSION}" \
    --signature "tenv_${LATEST_VERSION}_checksums.txt.sig" \
    --certificate "tenv_${LATEST_VERSION}_checksums.txt.pem" \
    --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
    "tenv_${LATEST_VERSION}_checksums.txt"

TENV_SIG_CHECK=$?

cosign \
    verify-blob \
    --certificate-identity "https://github.com/tofuutils/tenv/.github/workflows/release.yml@refs/tags/${LATEST_VERSION}" \
    --signature "tenv_${LATEST_VERSION}_amd64.rpm.sig" \
    --certificate "tenv_${LATEST_VERSION}_amd64.rpm.pem" \
    --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
    "tenv_${LATEST_VERSION}_amd64.rpm"

TENV_ASSET_CHECK=$?

# Check everything is good before installation
if [ "$TENV_SIG_CHECK" -eq "0" ] && [ "$TENV_ASSET_CHECK" -eq "0" ] && shasum -a 256 -c "tenv_${LATEST_VERSION}_checksums.txt" --ignore-missing
then
  dnf install "tenv_${LATEST_VERSION}_amd64.rpm" -y
  tenv --version
else
  echo "Signature verification and/or checksum checks failed!"
fi
```

<a id="contributing"></a>
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

Check out our [contributing guide](CONTRIBUTING.md) to get started.

Don't forget to give the project a star! Thanks again!

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
