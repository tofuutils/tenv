# tenv

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

## LICENSE
- [tenv inself](https://github.com/tofuutils/tenv/blob/main/LICENSE)
- [tofuenv](https://github.com/tofuutils/tofuenv/blob/main/LICENSE)
- [tfenv](https://github.com/tfutils/tfenv/blob/master/LICENSE)
  - tofuenv uses tfenv's source code
