# @summary Default parameters for tenv module
#
# @api private
#
class tenv::params {
  case $facts['os']['family'] {
    'Debian': {
      $package_provider = 'dpkg'
      $package_extension = 'deb'
      $cosign_package_format = 'deb'
    }
    'RedHat': {
      $package_provider = 'rpm'
      $package_extension = 'rpm'
      $cosign_package_format = 'rpm'
    }
    'Archlinux': {
      $package_provider = 'pacman'
      $package_extension = 'pkg.tar.zst'
      $cosign_package_format = 'deb' # Will use AUR instead
    }
    default: {
      fail("Unsupported OS family: ${facts['os']['family']}")
    }
  }
}
