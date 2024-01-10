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

Support same [environment variables](#Environment Variables) as `gotofuenv`.

### gotofuenv install [version]

Install a requested version of OpenTofu (into GOTOFUENV_ROOT directory from GOTOFUENV_REMOTE url).

Without parameter the version to use is resolved automatically via GOTOFUENV_TOFU_VERSION or `.opentofu-version` files
(searched in working directory, user home directory and GOTOFUENV_ROOT directory).
Use "latest" when none are found.

If a parameter is passed, available options:

- an exact [Semver 2.0.0](https://semver.org/) version string to install
- a Semver constraint string (checked against available at GOTOFUENV_REMOTE url)
- latest (checked against available at GOTOFUENV_REMOTE url)
- latest-allowed is a syntax to scan your OpenTofu files to detect which version is maximally allowed.
- min-required is a syntax to scan your OpenTofu files to detect which version is minimally required.

See [required_version](https://opentofu.org/docs/language/settings/) docs.

```console
$ gotofuenv install 1.6.0-beta5
$ gotofuenv install ">= 1.6.0-rc1" 
$ gotofuenv install latest
$ gotofuenv install latest-allowed
$ gotofuenv install min-required
```

### Environment Variables

#### GOTOFUENV_GITHUB_TOKEN

String (Default: "")

Allow to specify a GitHub token to increase [GitHub Rate limits for the REST API](https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api). Useful because OpenTofu binares are downloaded from the OpenTofu GitHub repository.

`gotofuenv` support a `-t` flag version.

#### GOTOFUENV_AUTO_INSTALL

String (Default: true)

If set to true gotofuenv will automatically install missing OpenTofu version needed (fallback to latest-allowed strategy when no `.opentofu-version` files are found).

`gotofuenv use` support a `-n` disabling flag version.

Example: use 1.6.0-rc1 version that is not installed, and auto installation is disabled. (-v flag is equivalent to `GOTOFUENV_VERBOSE=true`)

```console
$ GOTOFUENV_AUTO_INSTALL=false gotofuenv use -v 1.6.0-rc1
Write 1.6.0-rc1 in /home/dvaumoron/.gotofuenv/.opentofu-version
```

Example: use 1.6.0-rc1 version that is not installed, and auto installation stay enabled.

```console
$ gotofuenv use -v 1.6.0-rc1
Installation of OpenTofu 1.6.0-rc1
Search asset tofu_1.6.0-rc1_linux_amd64.zip for release v1.6.0-rc1
Write 1.6.0-rc1 in /home/dvaumoron/.gotofuenv/.opentofu-version
```

#### GOTOFUENV_VERBOSE

String (Default: false)

Active the verbose display of gotofuenv.

`gotofuenv` support a `-v` flag version.

#### GOTOFUENV_REMOTE

String (Default: https://api.github.com/repos/opentofu/opentofu/releases)

To install from a remote other than the default (must comply with github REST API)

#### GOTOFUENV_ROOT

Path (Default: `$HOME/.gotofuenv`)

The path to a directory where the local tofu versions and configuration files exist.

#### GOTOFUENV_TOFU_VERSION

String (Default: "")

If not empty string, this variable overrides OpenTofu version, specified in `.opentofu-version` files.
`gotofuenv install` command also respects this variable.

e.g. with :

```console
$ tofu version
OpenTofu v1.6.0
on linux_amd64
```

then :

```console
$ GOTOFUENV_TOFU_VERSION=1.6.0-rc1 tofu version
OpenTofu v1.6.0-rc1
on linux_amd64
```

### gotofuenv use version

Switch the default OpenTofu version to use (set in `.opentofu-version` file in GOTOFUENV_ROOT).

`gotofuenv use` has a `-w` flag to write `.opentofu-version` file in working directory.

Available parameter options:

- an exact [Semver 2.0.0](https://semver.org/) version string to use
- a Semver constraint string (checked against available in GOTOFUENV_ROOT directory)
- latest (checked against available in GOTOFUENV_ROOT directory)
- latest-allowed is a syntax to scan your OpenTofu files to detect which version is maximally allowed.
- min-required is a syntax to scan your OpenTofu files to detect which version is minimally required.

```console
$ tofuenv use min-required
$ tofuenv use v1.6.0-beta5
$ tofuenv use latest
$ tofuenv use latest-allowed
```

### tofuenv uninstall version

Uninstall a specific version of OpenTofu (remove it from GOTOFUENV_ROOT directory without interpretation).

```console
$ tofuenv uninstall 0.7.0
$ tofuenv uninstall latest
$ tofuenv uninstall latest:^0.8
```

### tofuenv list

List installed versions

```console
$ tofuenv list
  1.6.0-alpha5
* 1.6.0-rc1 (set by /opt/.tofuenv/version)
```

### gotofuenv list-remote

List installable OpenTofu versions (from GOTOFUENV_REMOTE url), sorted in ascending version order.

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

## .opentofu-version file

If you put a `.opentofu-version` file  in working directory, user home directory or GOTOFUENV_ROOT directory, gotofuenv detects it and uses the version written in it.
Note, that GOTOFUENV_TOFU_VERSION can be used to override version specified by `.opentofu-version` file.

Available value recognized (same as `gotofuenv use` command) :

- an exact [Semver 2.0.0](https://semver.org/) version string to use
- a Semver constraint string (checked against available in GOTOFUENV_ROOT directory)
- latest (checked against available in GOTOFUENV_ROOT directory)
- latest-allowed is a syntax to scan your OpenTofu files to detect which version is maximally allowed.
- min-required is a syntax to scan your OpenTofu files to detect which version is minimally required.

## LICENSE

The GoTofuEnv project is released under the Apache 2.0 license. See [LICENSE](LICENSE).
