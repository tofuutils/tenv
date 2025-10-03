# @summary Manages tenv installation and configuration
#
# This class installs and configures tenv, a version manager for
# OpenTofu, Terraform, Terragrunt, Terramate, and Atmos.
#
# @param version
#   Version of tenv to install. Use 'latest' for the most recent version.
#
# @param install_cosign
#   Whether to install cosign for signature verification.
#
# @param configure_shell
#   Whether to configure shell environments for users.
#
# @param setup_completion
#   Whether to setup shell completion.
#
# @param shell
#   Shell type to configure (bash, zsh, or fish).
#
# @param auto_install
#   Enable automatic installation of missing tool versions.
#
# @param root_path
#   Root directory for tenv installations.
#
# @param users
#   Array of users to configure tenv for.
#
# @param arch
#   Architecture for binary downloads.
#
# @param github_token
#   GitHub token for API rate limits (optional).
#
# @param manage_prerequisites
#   Whether to manage prerequisite packages.
#
# @param update_path
#   Whether to update PATH in shell configuration.
#
# @example Basic usage
#   include tenv
#
# @example Custom configuration
#   class { 'tenv':
#     version        => 'v2.6.1',
#     install_cosign => true,
#     users          => ['user1', 'user2'],
#     shell          => 'zsh',
#   }
#
class tenv (
  String $version                      = 'latest',
  Boolean $install_cosign              = true,
  Boolean $configure_shell             = true,
  Boolean $setup_completion            = true,
  Enum['bash', 'zsh', 'fish'] $shell   = 'bash',
  Boolean $auto_install                = false,
  Stdlib::Absolutepath $root_path      = "/root/.tenv",
  Array[String] $users                 = ['root'],
  String $arch                         = 'amd64',
  Optional[String] $github_token       = undef,
  Boolean $manage_prerequisites        = true,
  Boolean $update_path                 = true,
) {
  # Validate OS support
  unless $facts['os']['family'] in ['Debian', 'RedHat', 'Archlinux'] {
    fail("Unsupported operating system: ${facts['os']['family']}")
  }

  # Include parameter class
  include tenv::params

  # Manage prerequisite packages
  if $manage_prerequisites {
    case $facts['os']['family'] {
      'Debian': {
        ensure_packages(['curl', 'jq', 'unzip', 'ca-certificates'])
      }
      'RedHat': {
        ensure_packages(['curl', 'jq', 'unzip', 'ca-certificates'])
      }
      'Archlinux': {
        ensure_packages(['curl', 'jq', 'unzip', 'ca-certificates'])
      }
      default: {
        fail("Unsupported OS family: ${facts['os']['family']}")
      }
    }
  }

  # Install cosign if requested
  if $install_cosign {
    contain tenv::cosign
  }

  # Install tenv
  contain tenv::install

  # Configure tenv
  if $configure_shell {
    contain tenv::config
  }

  # Ordering
  if $manage_prerequisites {
    Package['curl', 'jq', 'unzip', 'ca-certificates']
    -> Class['tenv::install']
  }

  if $install_cosign {
    Class['tenv::cosign']
    -> Class['tenv::install']
  }

  if $configure_shell {
    Class['tenv::install']
    -> Class['tenv::config']
  }
}
