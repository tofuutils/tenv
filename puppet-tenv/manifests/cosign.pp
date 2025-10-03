# @summary Installs cosign for signature verification
#
# @api private
#
class tenv::cosign {
  include tenv::params

  # Get latest cosign version
  $version_url = 'https://api.github.com/repos/sigstore/cosign/releases/latest'
  $version_cmd = "curl -s ${version_url} | grep '\"tag_name\":' | sed -E 's/.*\"v([^\"]+)\".*/\\1/'"

  exec { 'get_cosign_version':
    command => "${version_cmd} > /tmp/cosign_version",
    path    => ['/usr/bin', '/usr/local/bin', '/bin'],
    creates => '/tmp/cosign_version',
  }

  case $facts['os']['family'] {
    'Debian': {
      exec { 'download_cosign':
        command => "curl -sL https://github.com/sigstore/cosign/releases/latest/download/cosign_\$(cat /tmp/cosign_version)_amd64.deb -o /tmp/cosign.deb",
        path    => ['/usr/bin', '/usr/local/bin', '/bin'],
        creates => '/tmp/cosign.deb',
        require => Exec['get_cosign_version'],
      }

      package { 'cosign':
        ensure   => installed,
        provider => 'dpkg',
        source   => '/tmp/cosign.deb',
        require  => Exec['download_cosign'],
      }

      file { '/tmp/cosign.deb':
        ensure  => absent,
        require => Package['cosign'],
      }
    }

    'RedHat': {
      exec { 'download_cosign':
        command => "curl -sL https://github.com/sigstore/cosign/releases/latest/download/cosign-\$(cat /tmp/cosign_version)-1.x86_64.rpm -o /tmp/cosign.rpm",
        path    => ['/usr/bin', '/usr/local/bin', '/bin'],
        creates => '/tmp/cosign.rpm',
        require => Exec['get_cosign_version'],
      }

      package { 'cosign':
        ensure   => installed,
        provider => 'rpm',
        source   => '/tmp/cosign.rpm',
        require  => Exec['download_cosign'],
      }

      file { '/tmp/cosign.rpm':
        ensure  => absent,
        require => Package['cosign'],
      }
    }

    'Archlinux': {
      package { 'cosign':
        ensure => installed,
      }
    }

    default: {
      fail("Unsupported OS family: ${facts['os']['family']}")
    }
  }
}
