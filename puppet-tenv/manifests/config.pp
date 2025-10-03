# @summary Configures tenv for users
#
# @api private
#
class tenv::config {
  $tenv::users.each |String $user| {
    # Get user home directory
    $user_home = $user ? {
      'root'  => '/root',
      default => "/home/${user}",
    }

    # Create tenv root directory
    file { "${user_home}/.tenv":
      ensure => directory,
      owner  => $user,
      group  => $user,
      mode   => '0755',
    }

    # Configure shell based on type
    case $tenv::shell {
      'bash': {
        $shell_config = "${user_home}/.bashrc"
        $completion_file = "${user_home}/.tenv.completion.bash"
      }
      'zsh': {
        $shell_config = "${user_home}/.zshrc"
        $completion_file = "${user_home}/.tenv.completion.zsh"
      }
      'fish': {
        $shell_config = "${user_home}/.config/fish/config.fish"
        $completion_file = "${user_home}/.tenv.completion.fish"
        
        # Ensure fish config directory exists
        file { "${user_home}/.config/fish":
          ensure => directory,
          owner  => $user,
          group  => $user,
          mode   => '0755',
        }
      }
      default: {
        fail("Unsupported shell: ${tenv::shell}")
      }
    }

    # Add tenv configuration to shell
    file_line { "tenv_root_${user}":
      path  => $shell_config,
      line  => "export TENV_ROOT=\"${tenv::root_path}\"",
      match => '^export TENV_ROOT=',
    }

    if $tenv::auto_install {
      file_line { "tenv_auto_install_${user}":
        path  => $shell_config,
        line  => 'export TENV_AUTO_INSTALL=true',
        match => '^export TENV_AUTO_INSTALL=',
      }
    }

    if $tenv::github_token {
      file_line { "tenv_github_token_${user}":
        path  => $shell_config,
        line  => "export TENV_GITHUB_TOKEN=\"${tenv::github_token}\"",
        match => '^export TENV_GITHUB_TOKEN=',
      }
    }

    if $tenv::update_path {
      file_line { "tenv_path_${user}":
        path  => $shell_config,
        line  => 'export PATH=$(tenv update-path)',
        match => '^export PATH=.*tenv update-path',
      }
    }

    # Setup shell completion
    if $tenv::setup_completion {
      exec { "tenv_completion_${user}_${tenv::shell}":
        command     => "tenv completion ${tenv::shell} > ${completion_file}",
        path        => ['/usr/bin', '/usr/local/bin', '/bin'],
        creates     => $completion_file,
        environment => ["HOME=${user_home}"],
        user        => $user,
        require     => File["${user_home}/.tenv"],
      }

      file_line { "source_tenv_completion_${user}":
        path    => $shell_config,
        line    => "source ${completion_file}",
        match   => "^source.*tenv.completion.${tenv::shell}",
        require => Exec["tenv_completion_${user}_${tenv::shell}"],
      }
    }
  }
}
