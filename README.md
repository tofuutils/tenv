# GoTofuEnv

[OpenTofu](https://opentofu.org) version manager (inspired by [tofuenv](https://github.com/tofuutils/tofuenv), written in Go)

Handle [Semver 2.0.0](https://semver.org/) with [go-version](https://github.com/hashicorp/go-version) and use the [HCL](https://github.com/hashicorp/hcl) parser to extract required version constraint from OpenTofu files.

## Installation

### Automatic

Install via [Homebrew](https://brew.sh/)

```console
$ brew tap dvaumoron/tap
$ brew install tofuenv
```

### Manual

Get the last packaged binaries (use .deb, .rpm, .apk or .zip) found [here](https://github.com/dvaumoron/gotofuenv/releases).

For the .zip case, the unzipped folder must be added to your PATH.

## Usage

### tofu

This project version of `tofu` command is a proxy to OpenTofu `tofu` command  managed by `gotofuenv`.

### gotofuenv install [version]

Install a requested version of OpenTofu (into TOFUENV_ROOT directory from TOFUENV_REMOTE url).

Without parameter the version to use is resolved automatically via TOFUENV_TOFU_VERSION or [`.opentofu-version`](#opentofu-version-file) files
(searched in working directory, user home directory and TOFUENV_ROOT directory).
Use "latest-stable" when none are found.

If a parameter is passed, available options:

- an exact [Semver 2.0.0](https://semver.org/) version string to install
- a [version constraint](https://opentofu.org/docs/language/expressions/version-constraints) string (checked against available at TOFUENV_REMOTE url)
- latest or latest-stable (checked against available at TOFUENV_REMOTE url)
- latest-allowed or min-required to scan your OpenTofu files to detect which version is maximally allowed or minimally required.

See [required_version](https://opentofu.org/docs/language/settings#specifying-a-required-opentofu-version) docs.

```console
$ gotofuenv install 1.6.0-beta5
$ gotofuenv install "~> 1.6.0"
$ gotofuenv install latest
$ gotofuenv install latest-stable
$ gotofuenv install latest-allowed
$ gotofuenv install min-required
```

### Environment Variables

Both command support the following environment variables.

#### TOFUENV_AUTO_INSTALL

String (Default: true)

If set to true gotofuenv will automatically install missing OpenTofu version needed (fallback to latest-allowed strategy when no [`.opentofu-version`](#opentofu-version-file) files are found).

`gotofuenv use` support a `--no-install`, `-n` disabling flag version.

Example: use 1.6.0-rc1 version that is not installed, and auto installation is disabled. (-v flag is equivalent to `TOFUENV_VERBOSE=true`)

```console
$ TOFUENV_AUTO_INSTALL=false gotofuenv use -v 1.6.0-rc1
Write 1.6.0-rc1 in /home/dvaumoron/.gotofuenv/.opentofu-version
```

Example: use 1.6.0-rc1 version that is not installed, and auto installation stay enabled.

```console
$ gotofuenv use -v 1.6.0-rc1
Installation of OpenTofu 1.6.0-rc1
Search asset tofu_1.6.0-rc1_linux_amd64.zip for release v1.6.0-rc1
Write 1.6.0-rc1 in /home/dvaumoron/.gotofuenv/.opentofu-version
```

#### TOFUENV_GITHUB_TOKEN

String (Default: "")

Allow to specify a GitHub token to increase [GitHub Rate limits for the REST API](https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api). Useful because OpenTofu binares are downloaded from the OpenTofu GitHub repository.

`gotofuenv` support a `--github-token`, `-t` flag version.

#### TOFUENV_REMOTE

String (Default: https://api.github.com/repos/opentofu/opentofu/releases)

To install from a remote other than the default (must comply with [github REST API](https://docs.github.com/en/rest?apiVersion=2022-11-28))

`gotofuenv` support a `--remote-url`, `-u` flag version.

#### TOFUENV_ROOT

Path (Default: `$HOME/.gotofuenv`)

The path to a directory where the local OpenTofu versions and GoTofuEnv configuration files exist.

`gotofuenv` support a `--root-path`, `-r` flag version.

#### TOFUENV_TOFU_VERSION

String (Default: "")

If not empty string, this variable overrides OpenTofu version, specified in [`.opentofu-version`](#opentofu-version-file) files.
`gotofuenv install` command also respects this variable.

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

#### TOFUENV_VERBOSE

String (Default: false)

Active the verbose display of gotofuenv.

`gotofuenv` support a `--verbose`, `-v` flag version.

### gotofuenv use version

Switch the default OpenTofu version to use (set in [`.opentofu-version`](#opentofu-version-file) file in TOFUENV_ROOT).

`gotofuenv use` has a `--force-remote`, `-f` flag to force remote search.

`gotofuenv use` has a `--working-dir`, `-w` flag to write [`.opentofu-version`](#opentofu-version-file) file in working directory.

Available parameter options:

- an exact [Semver 2.0.0](https://semver.org/) version string to use
- a [version constraint](https://opentofu.org/docs/language/expressions/version-constraints) string (checked against available in TOFUENV_ROOT directory)
- latest or latest-stable (checked against available in TOFUENV_ROOT directory)
- latest-allowed or min-required to scan your OpenTofu files to detect which version is maximally allowed or minimally required.

See [required_version](https://opentofu.org/docs/language/settings#specifying-a-required-opentofu-version) docs.

```console
$ gotofuenv use min-required
$ gotofuenv use v1.6.0-beta5
$ gotofuenv use latest
$ gotofuenv use latest-allowed
```

### gotofuenv reset

Reset used version of OpenTofu (remove .opentofu-version file from TOFUENV_ROOT).

```console
$ gotofuenv reset
```

### gotofuenv uninstall version

Uninstall a specific version of OpenTofu (remove it from TOFUENV_ROOT directory without interpretation).

```console
$ gotofuenv uninstall v1.6.0-alpha4
```

### gotofuenv list

List installed OpenTofu versions (located in TOFUENV_ROOT directory), sorted in ascending version order.

`gotofuenv list` has a `--descending`, `-d` flag to sort in descending order.

```console
$ gotofuenv list
  1.6.0-rc1 
* 1.6.0 (set by /home/dvaumoron/.gotofuenv/.opentofu-version )
```

### gotofuenv list-remote

List installable OpenTofu versions (from TOFUENV_REMOTE url), sorted in ascending version order.

`gotofuenv list-remote` has a `--descending`, `-d` flag to sort in descending order.

`gotofuenv list-remote` has a `--stable`, `-s` flag to display only stable version.

```console
$ gotofuenv list-remote
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

### gotofuenv help [command]

Help about any command.

```console
$ gotofuenv help reset
Reset used version of OpenTofu (remove .opentofu-version file from TOFUENV_ROOT).

Usage:
  gotofuenv reset [flags]

Flags:
  -h, --help   help for reset

Global Flags:
  -t, --github-token string   GitHub token (increases GitHub REST API rate limits)
  -u, --remote-url string     remote url to install from (default "https://api.github.com/repos/opentofu/opentofu/releases")
  -r, --root-path string      local path to install OpenTofu versions (default "/home/dvaumoron/.gotofuenv")
  -v, --verbose               verbose output
```

## .opentofu-version file

If you put a `.opentofu-version` file  in working directory, user home directory or TOFUENV_ROOT directory, gotofuenv detects it and uses the version written in it.
Note, that TOFUENV_TOFU_VERSION can be used to override version specified by `.opentofu-version` file.

Recognized value (same as `gotofuenv use` command) :

- an exact [Semver 2.0.0](https://semver.org/) version string to use
- a [version constraint](https://opentofu.org/docs/language/expressions/version-constraints) string (checked against available in TOFUENV_ROOT directory)
- latest or latest-stable (checked against available in TOFUENV_ROOT directory)
- latest-allowed or min-required to scan your OpenTofu files to detect which version is maximally allowed or minimally required.

See [required_version](https://opentofu.org/docs/language/settings#specifying-a-required-opentofu-version) docs.

## LICENSE

The GoTofuEnv project is released under the Apache 2.0 license. See [LICENSE](LICENSE).
