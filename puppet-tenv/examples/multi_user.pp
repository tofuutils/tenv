# Configure tenv for multiple users with zsh
class { 'tenv':
  users => ['user1', 'user2', 'jenkins'],
  shell => 'zsh',
}
