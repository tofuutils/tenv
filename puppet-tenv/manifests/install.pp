# @summary Installs tenv
#
# @api private
#
class tenv::install {
  include tenv::params

  # Determine version to install
  if $tenv::version == 'latest' {
    $version_url = 'https://api.github.com/repos/tofuutils/tenv/releases/latest'
    $version_cmd = "curl -s ${version_url} | grep '\"tag_name\":' | sed -E 's/.*\"([^\"]+)\".*/\\1/'"
    
    exec { 'get_tenv_version':
      command => "${version_cmd} > /tmp/tenv_version",
      path    => ['/usr/bin', '/usr/local/bin', '/bin'],
      creates => '/tmp/tenv_version',
    }
    
    $install_version = file('/tmp/tenv_version')
  } else {
    $install_version = $tenv::version
  }

  case $facts['os']['family'] {
    'Debian': {
      $package_url = "https://github.com/tofuutils/tenv/releases/download/${install_version}/tenv_${install_version}_amd64.deb"
      $package_path = "/tmp/tenv_${install_version}_amd64.deb"

      exec { 'download_tenv':
        command => "curl -sL ${package_url} -o ${package_path}",
        path    => ['/usr/bin', '/usr/local/bin', '/bin'],
        creates => $package_path,
      }

      package { 'tenv':
        ensure   => installed,
        provider => 'dpkg',
        source   => $package_path,
        require  => Exec['download_tenv'],
      }

      # Cleanup
      file { $package_path:
        ensure  => absent,
        require => Package['tenv'],
      }
    }

    'RedHat': {
      $package_url = "https://github.com/tofuutils/tenv/releases/download/${install_version}/tenv_${install_version}_amd64.rpm"
      $package_path = "/tmp/tenv_${install_version}_amd64.rpm"

      exec { 'download_tenv':
        command => "curl -sL ${package_url} -o ${package_path}",
        path    => ['/usr/bin', '/usr/local/bin', '/bin'],
        creates => $package_path,
      }

      package { 'tenv':
        ensure   => installed,
        provider => 'rpm',
        source   => $package_path,
        require  => Exec['download_tenv'],
      }

      # Cleanup
      file { $package_path:
        ensure  => absent,
        require => Package['tenv'],
      }
    }

    'Archlinux': {
      # Install from AUR using yay
      exec { 'install_tenv_aur':
        command => 'yay -S --noconfirm tenv-bin',
        path    => ['/usr/bin', '/usr/local/bin', '/bin'],
        unless  => 'pacman -Q tenv-bin',
      }
    }

    default: {
      fail("Unsupported OS family: ${facts['os']['family']}")
    }
  }

  # Verify installation
  exec { 'verify_tenv':
    command => 'tenv version',
    path    => ['/usr/bin', '/usr/local/bin', '/bin'],
    require => $facts['os']['family'] ? {
      'Archlinux' => Exec['install_tenv_aur'],
      default     => Package['tenv'],
    },
  }
}
