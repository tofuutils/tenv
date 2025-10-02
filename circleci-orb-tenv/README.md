A CircleCI Orb for installing and managing [tenv](https://github.com/tofuutils/tenv) - a version manager for OpenTofu, Terraform, Terragrunt, Terramate, and Atmos.

This still needs to be added to CircleCI Orbs Packages repository

## Features

- Install tenv with a single command
- Optional cosign installation for signature verification
- Support for specific versions or latest release
- Cache support for faster builds
- Multiple installation examples

## Usage

### Quick Start

```yaml
version: 2.1

orbs:
  tenv: tofuutils/tenv@1.0.0

workflows:
  main:
    jobs:
      - tenv/install:
          version: latest
```

### Install Command

Use the install command in your existing jobs:

```yaml
version: 2.1

orbs:
  tenv: tofuutils/tenv@1.0.0

jobs:
  deploy:
    docker:
      - image: cimg/base:stable
    steps:
      - checkout
      - tenv/install:
          version: latest
          install-cosign: true
      - run:
          name: Use Terraform
          command: |
            tenv tf install 1.6.0
            terraform version
```

### Complete Example

```yaml
version: 2.1

orbs:
  tenv: tofuutils/tenv@1.0.0

jobs:
  terraform-deploy:
    docker:
      - image: cimg/base:stable
    steps:
      - checkout
      - tenv/install:
          version: v2.6.1
          install-cosign: true
      - tenv/use-version:
          tool: terraform
          version: "1.6.0"
      - run:
          name: Terraform Init
          command: terraform init
      - run:
          name: Terraform Plan
          command: terraform plan

workflows:
  deploy:
    jobs:
      - terraform-deploy
```

## Commands

### install

Install tenv on the executor.

**Parameters:**

- `version` (string, default: "latest"): Version of tenv to install
- `install-cosign` (boolean, default: true): Install cosign for signature verification
- `github-token` (env_var_name, default: ""): GitHub token for API rate limits

**Example:**

```yaml
- tenv/install:
    version: v2.6.1
    install-cosign: true
```

### install-cosign

Install cosign for signature verification.

**Parameters:**

- `version` (string, default: "latest"): Version of cosign to install

**Example:**

```yaml
- tenv/install-cosign:
    version: latest
```

### use-version

Install and set a specific version of a tool.

**Parameters:**

- `tool` (enum: ["terraform", "tofu", "terragrunt", "terramate", "atmos"]): Tool to configure
- `version` (string): Version to install and use

**Example:**

```yaml
- tenv/use-version:
    tool: terraform
    version: "1.6.0"
```

## Jobs

### install

A simple job that installs tenv.

**Parameters:**

- `version` (string, default: "latest"): Version of tenv to install
- `install-cosign` (boolean, default: true): Install cosign
- `executor` (executor, default: docker): Executor to use

**Example:**

```yaml
workflows:
  main:
    jobs:
      - tenv/install:
          version: latest
```

### test-version

Install tenv and verify a specific tool version.

**Parameters:**

- `tenv-version` (string, default: "latest"): Version of tenv
- `tool` (enum): Tool to test
- `tool-version` (string): Tool version to test

**Example:**

```yaml
workflows:
  test:
    jobs:
      - tenv/test-version:
          tool: terraform
          tool-version: "1.6.0"
```

## Executors

The orb uses standard CircleCI Docker executors:

- `cimg/base:stable` - Default executor

## Contributing

Contributions are welcome! Please check the [GitHub repository](https://github.com/tofuutils/tenv) for contribution guidelines.

