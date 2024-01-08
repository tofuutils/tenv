<p align="center">
  <img alt="golangci-lint logo" src="assets/logo.jpeg" height="150" />
  <h3 align="center">tenv</h3>
  <p align="center">Terraform/OpenTofu version manager</p>
</p>


[tenv](https://github.com/tofuutils/tenv) version manager that build on top of [tofuenv](https://github.com/tofuutils/tofuenv) and [tfenv](https://github.com/tfutils/tfenv) and manages [Terraform](https://www.terraform.io/) and [OpenTofu](https://opentofu.org/) binaries

## Support

Currently, tenv supports the following operating systems:

- macOS
  - 64bit
  - Arm (Apple Silicon)
- Linux
  - 64bit
  - Arm
- Windows (64bit) - only tested in git-bash - currently presumed failing due to symlink issues in git-bash

## Installation
WIP

## Environment variables
TENV_ROOT_PATH
TENV_TOFUENV_VERSION
TENV_TFENV_VERSION

## Contributors

This project exists thanks to all the people who contribute. [How to contribute](todo).

### Core Team

<details>
<summary>About core team</summary>

The tenv Core Team is a group of contributors that have demonstrated a lasting enthusiasm for the project and community.
The tenv Core Team has GitHub admin privileges on the repo.

## LICENSE
- [tenv inself](https://github.com/tofuutils/tenv/blob/main/LICENSE)
- [tofuenv](https://github.com/tofuutils/tofuenv/blob/main/LICENSE)
- [tfenv](https://github.com/tfutils/tfenv/blob/master/LICENSE)
  - tofuenv uses tfenv's source code
