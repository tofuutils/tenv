Ansible role to install and configure [tenv](https://github.com/tofuutils/tenv) - a version manager for OpenTofu, Terraform, Terragrunt, Terramate, and Atmos.

## Requirements

- Ansible 2.10 or higher
- Supported OS:
  - Ubuntu 20.04, 22.04, 24.04
  - Debian 11, 12
  - RHEL/CentOS/Rocky 8, 9
  - Arch Linux

## Role Variables

Available variables are listed below, along with default values (see `defaults/main.yml`):

```yaml
# Version to install ("latest" or specific version like "v2.6.1")
tenv_version: latest

# Install cosign for signature verification
tenv_install_cosign: true

# Configure shell environment
tenv_configure_shell: true

# Update PATH automatically
tenv_update_path: true

# Setup shell completion
tenv_setup_completion: true

# Shell type (bash, zsh, fish)
tenv_shell: bash

# Enable auto-install of tool versions
tenv_auto_install: false

# Root path for tenv installations
tenv_root_path: "{{ ansible_env.HOME }}/.tenv"

# Users to configure tenv for
tenv_users:
  - "{{ ansible_user_id }}"

# Architecture
tenv_arch: amd64

# GitHub token for API rate limits (optional)
tenv_github_token: ""
```

## Dependencies

None.

## Example Playbook

### Basic Installation

```yaml
---
- hosts: all
  become: yes
  roles:
    - role: ansible-role-tenv
```

### Custom Configuration

```yaml
---
- hosts: all
  become: yes
  roles:
    - role: ansible-role-tenv
      vars:
        tenv_version: v2.6.1
        tenv_auto_install: true
        tenv_shell: zsh
        tenv_users:
          - user1
          - user2
```

### Minimal Installation (without cosign)

```yaml
---
- hosts: all
  become: yes
  roles:
    - role: ansible-role-tenv
      vars:
        tenv_install_cosign: false
        tenv_setup_completion: false
```

## Post-Installation

After installation, you can start using tenv:

```bash
# Install a specific Terraform version
tenv tf install 1.6.0

# Install latest OpenTofu
tenv tofu install latest

# Use a specific version
tenv tf use 1.6.0

# List installed versions
tenv tf list

# Interactive mode
tenv
```

## Testing

To test this role:

```bash
cd tests/
ansible-playbook -i inventory test.yml
```

