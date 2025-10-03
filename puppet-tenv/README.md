# tenv

#### Table of Contents

1. [Description](#description)
2. [Setup](#setup)
    * [What tenv affects](#what-tenv-affects)
    * [Setup requirements](#setup-requirements)
    * [Beginning with tenv](#beginning-with-tenv)
3. [Usage](#usage)
    * [Basic Installation](#basic-installation)
    * [Custom Configuration](#custom-configuration)
    * [Multi-User Setup](#multi-user-setup)
    * [Advanced Examples](#advanced-examples)
4. [Reference](#reference)
    * [Classes](#classes)
    * [Parameters](#parameters)
5. [Development](#development)
    * [Testing](#testing)
6. [Support](#support)

## Description

This Puppet module installs and configures [tenv](https://github.com/tofuutils/tenv), a versatile version manager for Infrastructure as Code tools written in Go.

tenv manages multiple versions of:
- **OpenTofu** - Open-source Terraform alternative
- **Terraform** - HashiCorp's infrastructure provisioning tool
- **Terragrunt** - Terraform wrapper for DRY configurations
- **Terramate** - Infrastructure as Code orchestration tool
- **Atmos** - Universal tool for DevOps and cloud automation

### Key Features

- **Version Management**: Easily switch between different versions of IaC tools
- **Signature Verification**: Optional cosign installation for binary verification
- **Multi-User Support**: Configure tenv for multiple users simultaneously
- **Shell Integration**: Automatic PATH updates and completion for bash, zsh, and fish
- **Semantic Versioning**: Support for version constraints and latest releases
- **HCL Parsing**: Automatic version detection from `.tf` and `.hcl` files

## Setup

### What tenv affects

This module will:

* Install tenv binary via your system's package manager (apt, yum, pacman)
* Optionally install cosign for signature verification
* Create `~/.tenv` directories for specified users
* Modify shell configuration files (`.bashrc`, `.zshrc`, or `.config/fish/config.fish`)
* Set up shell completion scripts
* Configure environment variables (`TENV_ROOT`, `TENV_AUTO_INSTALL`, etc.)

### Setup Requirements

**Required Puppet Modules:**
- `puppetlabs/stdlib` (>= 6.0.0 < 10.0.0)
- `puppetlabs/apt` (>= 7.0.0 < 10.0.0) - for Debian-based systems

**System Requirements:**
- Puppet >= 6.0.0
- Internet connectivity for downloading packages
- Root or sudo access for package installation

**Supported Operating Systems:**
- Ubuntu 20.04, 22.04, 24.04
- Debian 11, 12
- RHEL/CentOS/Rocky 8, 9
- Arch Linux

### Beginning with tenv

The simplest way to get started is:

```puppet
include tenv
```

This will install tenv with default settings for the root user.

## Usage

### Basic Installation

Install tenv with default parameters:

```puppet
class { 'tenv':
  version => 'latest',
}
```

### Custom Configuration

Install a specific version with cosign enabled:

```puppet
class { 'tenv':
  version        => 'v2.6.1',
  install_cosign => true,
  auto_install   => true,
}
```

### Multi-User Setup

Configure tenv for multiple users with zsh:

```puppet
class { 'tenv':
  users => ['developer', 'devops', 'jenkins'],
  shell => 'zsh',
}
```

### Advanced Examples

#### Complete Configuration with All Options

```puppet
class { 'tenv':
  version              => 'v2.6.1',
  install_cosign       => true,
  configure_shell      => true,
  setup_completion     => true,
  shell                => 'bash',
  auto_install         => true,
  root_path            => '/opt/tenv',
  users                => ['user1', 'user2'],
  arch                 => 'amd64',
  github_token         => lookup('github_api_token', String, 'first', undef),
  manage_prerequisites => true,
  update_path          => true,
}
```

#### CI/CD User Configuration

Configure for Jenkins or other CI/CD users:

```puppet
class { 'tenv':
  users           => ['jenkins', 'gitlab-runner'],
  auto_install    => true,
  configure_shell => true,
  shell           => 'bash',
}
```

#### Minimal Installation (No Shell Configuration)

Install tenv without modifying shell configurations:

```puppet
class { 'tenv':
  install_cosign      => false,
  configure_shell     => false,
  setup_completion    => false,
  manage_prerequisites => true,
}
```

#### Using Hiera for Configuration

Create `data/common.yaml`:

```yaml
---
tenv::version: 'v2.6.1'
tenv::install_cosign: true
tenv::auto_install: true
tenv::users:
  - 'developer'
  - 'devops'
tenv::shell: 'zsh'
tenv::github_token: "%{lookup('github_token')}"
```

Then in your manifest:

```puppet
include tenv
```

#### Role-Based Configuration

Create a role for infrastructure engineers:

```puppet
class role::infrastructure_engineer {
  class { 'tenv':
    users        => ['engineer1', 'engineer2'],
    auto_install => true,
    shell        => 'zsh',
  }
}
```

## Reference

### Classes

#### `tenv`

Main class for managing tenv installation and configuration.

#### `tenv::install`

Private class that handles tenv binary installation. Called by main class.

#### `tenv::config`

Private class that configures shell environments for users. Called by main class.

#### `tenv::cosign`

Private class that installs cosign for signature verification. Called by main class when `install_cosign => true`.

#### `tenv::params`

Private class containing OS-specific parameters.

### Parameters

#### `version`

Data type: `String`

Version of tenv to install. Use `'latest'` for the most recent release or specify a version like `'v2.6.1'`.

Default: `'latest'`

```puppet
class { 'tenv':
  version => 'v2.6.1',
}
```

#### `install_cosign`

Data type: `Boolean`

Whether to install cosign for signature verification of downloaded binaries.

Default: `true`

```puppet
class { 'tenv':
  install_cosign => false,
}
```

#### `configure_shell`

Data type: `Boolean`

Whether to configure shell environments for specified users.

Default: `true`

```puppet
class { 'tenv':
  configure_shell => false,
}
```

#### `setup_completion`

Data type: `Boolean`

Whether to setup shell completion scripts.

Default: `true`

```puppet
class { 'tenv':
  setup_completion => false,
}
```

#### `shell`

Data type: `Enum['bash', 'zsh', 'fish']`

Shell type to configure. Determines which shell configuration files are modified.

Default: `'bash'`

```puppet
class { 'tenv':
  shell => 'zsh',
}
```

#### `auto_install`

Data type: `Boolean`

Enable automatic installation of missing tool versions when detected.

Default: `false`

```puppet
class { 'tenv':
  auto_install => true,
}
```

#### `root_path`

Data type: `Stdlib::Absolutepath`

Root directory where tenv stores downloaded tool versions.

Default: `'/root/.tenv'` (or `/home/${user}/.tenv` for non-root users)

```puppet
class { 'tenv':
  root_path => '/opt/tenv',
}
```

#### `users`

Data type: `Array[String]`

List of system users to configure tenv for.

Default: `['root']`

```puppet
class { 'tenv':
  users => ['user1', 'user2', 'jenkins'],
}
```

#### `arch`

Data type: `String`

Architecture for binary downloads.

Default: `'amd64'`

```puppet
class { 'tenv':
  arch => 'arm64',
}
```

#### `github_token`

Data type: `Optional[String]`

GitHub personal access token for increased API rate limits. Useful for environments with many requests.

Default: `undef`

```puppet
class { 'tenv':
  github_token => 'ghp_xxxxxxxxxxxx',
}
```

#### `manage_prerequisites`

Data type: `Boolean`

Whether to manage prerequisite packages (curl, jq, unzip, ca-certificates).

Default: `true`

```puppet
class { 'tenv':
  manage_prerequisites => false,
}
```

#### `update_path`

Data type: `Boolean`

Whether to update PATH environment variable in shell configuration.

Default: `true`

```puppet
class { 'tenv':
  update_path => false,
}
```

### Network Requirements

The module requires internet connectivity to:
- Download tenv packages from GitHub releases
- Download cosign packages (if enabled)
- Query GitHub API for latest versions

### Arch Linux Requirements

On Arch Linux, the module expects the `yay` AUR helper to be installed for tenv installation from AUR. If `yay` is not present, the module will install `base-devel` but won't automatically install `yay` itself.

### Shell Configuration

The module modifies user shell configuration files. If you have custom shell configurations, review the changes to ensure compatibility.

### User Permissions

Users configured with tenv need appropriate permissions to:
- Read/write to their home directories
- Execute tenv commands
- Install tool versions in `TENV_ROOT`

## Development

### Testing

#### Prerequisites

Install development dependencies:

```bash
bundle install
```

#### Running Tests

**Unit Tests:**

```bash
# Run all unit tests
bundle exec rake spec

# Run specific test file
bundle exec rspec spec/classes/init_spec.rb

# Run with coverage
bundle exec rake spec SPEC_OPTS='--format documentation'
```

**Syntax Validation:**

```bash
# Validate manifests
bundle exec rake validate

# Run puppet-lint
bundle exec rake lint

# Check syntax of all files
bundle exec rake syntax
```

**Acceptance Tests:**

```bash
# Run full acceptance test suite
bundle exec rake beaker

# Run on specific platform
BEAKER_set=ubuntu-2204-x64 bundle exec rake beaker
```

**All Checks:**

```bash
# Run all tests and validations
bundle exec rake test
```

#### Test Structure

```
spec/
├── spec_helper.rb              # RSpec configuration
├── spec_helper_acceptance.rb   # Beaker configuration
├── default_facts.yml           # Default facts for tests
├── classes/                    # Unit tests
│   ├── init_spec.rb
│   ├── install_spec.rb
│   ├── config_spec.rb
│   └── cosign_spec.rb
└── acceptance/                 # Integration tests
    └── class_spec.rb
```

### Release Process

1. Update version in `metadata.json`
2. Update `CHANGELOG.md`
3. Run full test suite
4. Build module: `pdk build`
5. Tag release: `git tag v1.0.0`
6. Push tags: `git push --tags`
7. Publish to Puppet Forge: `puppet module publish pkg/tofuutils-tenv-1.0.0.tar.gz`

### Development Tools

**Puppet Development Kit (PDK):**

```bash
# Validate module
pdk validate

# Run unit tests
pdk test unit

# Create new class
pdk new class myclass
```

**Debugging:**

```bash
# Test manifest application
puppet apply --noop examples/basic.pp

# Verbose output
puppet apply --debug examples/basic.pp

# Check catalog
puppet apply --catalog examples/basic.pp


### Environment Variables

After installation, these environment variables are configured:

- `TENV_ROOT`: Root directory for tenv installations
- `TENV_AUTO_INSTALL`: Automatic version installation (if enabled)
- `TENV_GITHUB_TOKEN`: GitHub token for API rate limits (if provided)
- `PATH`: Updated to include tenv binaries

### Files Modified

The module modifies these files:

- `~/.bashrc` (if shell is bash)
- `~/.zshrc` (if shell is zsh)
- `~/.config/fish/config.fish` (if shell is fish)
- `~/.tenv.completion.{bash,zsh,fish}` (completion scripts)
- `~/.tenv/` (tenv root directory)

---

**Version**: 1.0.0  
**Last Updated**: 2025-10-04
