# Install tenv with cosign and enable auto-install
class { 'tenv':
  version        => 'v2.6.1',
  install_cosign => true,
  auto_install   => true,
}

