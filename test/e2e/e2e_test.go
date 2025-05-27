//go:build e2e

package e2e

import (
	"bytes"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInstallTerraform(t *testing.T) {
	//
	// Check the basic installation of a specific version.
	//
	tenvBin := os.Getenv("TENV_BIN")

	cmd := exec.Command(tenvBin, "tf", "install", "1.10.5", "-v")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()

	require.NoError(t, err, "Expected no error during the installation process")
	require.Contains(t, out.String(), "Installing Terraform 1.10.5")
	require.Contains(t, out.String(), "Installation of Terraform 1.10.5 successful")
}

func TestTFenvVersionEnvVariable(t *testing.T) {
	//
	// Check that tenv detects the version from the env,
	// but does not install it by default.
	//
	tenvBin := os.Getenv("TENV_BIN")

	cmd := exec.Command(tenvBin, "tf", "detect")

	env := os.Environ()
	env = append(env, "TFENV_TERRAFORM_VERSION=1.10.0")
	cmd.Env = env

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()

	require.NoError(t, err, "Expected no error during the installation process")
	require.NotContains(t, out.String(), "Installation of Terraform 1.10.0 successful")
	require.Contains(t, out.String(), "Resolved version from TFENV_TERRAFORM_VERSION : 1.10.0")
}

func TestTFenvVersionEnvVariableInstall(t *testing.T) {
	//
	// Check that tenv detects the version from the env,
	// and install it if '-i' flag provided.
	//
	tenvBin := os.Getenv("TENV_BIN")

	cmd := exec.Command(tenvBin, "tf", "detect", "-i")

	env := os.Environ()
	env = append(env, "TFENV_TERRAFORM_VERSION=1.10.0")
	cmd.Env = env

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()

	require.NoError(t, err, "Expected no error during the installation process")
	require.Contains(t, out.String(), "Installing Terraform 1.10.0")
	require.Contains(t, out.String(), "Installation of Terraform 1.10.0 successful")
}

func TestTFenvVersionLastUse(t *testing.T) {
	//
	// Check that the version file can be read by anyone.
	//
	tenvBin := os.Getenv("TENV_BIN")

	env := os.Environ()
	env = append(env, "TENV_ROOT=/usr/local/share/tenv", "TENV_AUTO_INSTALL=true")
	var out bytes.Buffer

	cmd_install := exec.Command("sudo", "--preserve-env=TENV_ROOT,TENV_AUTO_INSTALL", tenvBin, "tofu", "use", "latest")
	cmd_install.Env = env
	cmd_install.Stdout = &out
	cmd_install.Stderr = &out

	cmd_version := exec.Command("tofu", "--version")
	cmd_version.Env = env
	cmd_version.Stdout = &out
	cmd_version.Stderr = &out

	_ = cmd_install.Run()
	err := cmd_version.Run()

	require.NoErrorf(t, err, "Expected no error during the version check. Output:\n%s", out.String())
}

func TestTFenvTerragruntVersionDetect(t *testing.T) {
	//
	// Check that tenv detects the terragrunt version from the root.hcl file,
	// but does not install it by default.
	//
	tenvBin := os.Getenv("TENV_BIN")

	fileContent := `terragrunt_version_constraint = "0.69.1"`
	_ = os.WriteFile("root.hcl", []byte(fileContent), 0644)

	cmd := exec.Command(tenvBin, "tg", "detect")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()

	require.NoError(t, err, "Expected no error during the detect")
	require.Contains(t, out.String(), "0.69.1")
	require.NotContains(t, out.String(), "Installation of Terragrunt 0.69.1 successful")
}

func TestTofuVersionFromFile(t *testing.T) {
	t.Parallel()

	tenvBin := os.Getenv("TENV_BIN")

	// Test .opentofu-version file
	err := os.WriteFile(".opentofu-version", []byte("1.6.0"), 0644)
	if err != nil {
		t.Fatal("Failed to create .opentofu-version file:", err)
	}
	defer os.Remove(".opentofu-version")

	cmd := exec.Command(tenvBin, "tofu", "detect")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
	require.NoError(t, err, "Expected no error during version detection")
	require.Contains(t, out.String(), "1.6.0")
}

func TestTofuVersionFromAsdf(t *testing.T) {
	t.Parallel()

	tenvBin := os.Getenv("TENV_BIN")

	// Test .tool-versions file
	toolVersions := `opentofu 1.6.1
terraform 1.7.0
`
	err := os.WriteFile(".tool-versions", []byte(toolVersions), 0644)
	if err != nil {
		t.Fatal("Failed to create .tool-versions file:", err)
	}
	defer os.Remove(".tool-versions")

	cmd := exec.Command(tenvBin, "tofu", "detect")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
	require.NoError(t, err, "Expected no error during version detection")
	require.Contains(t, out.String(), "1.6.1")
}

func TestTofuVersionFromTerragrunt(t *testing.T) {
	t.Parallel()

	tenvBin := os.Getenv("TENV_BIN")

	// Test terragrunt.hcl file
	hclContent := `
terraform {
  source = "..."

  # Configure OpenTofu version constraint
  opentofu_version_constraint = "~> 1.6.0"
}
`
	err := os.WriteFile("terragrunt.hcl", []byte(hclContent), 0644)
	if err != nil {
		t.Fatal("Failed to create terragrunt.hcl file:", err)
	}
	defer os.Remove("terragrunt.hcl")

	cmd := exec.Command(tenvBin, "tofu", "detect")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
	require.NoError(t, err, "Expected no error during version detection")
	require.Contains(t, out.String(), "1.6.0")
}

func TestTofuVersionFromEnv(t *testing.T) {
	t.Parallel()

	tenvBin := os.Getenv("TENV_BIN")

	cmd := exec.Command(tenvBin, "tofu", "detect")
	env := os.Environ()
	env = append(env, "TOFUENV_TOFU_VERSION=1.6.2")
	cmd.Env = env

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	require.NoError(t, err, "Expected no error during version detection")
	require.Contains(t, out.String(), "1.6.2")
}

func TestTofuVersionFromHCL(t *testing.T) {
	t.Parallel()

	tenvBin := os.Getenv("TENV_BIN")

	// Test .tf file
	hclContent := `
terraform {
  required_version = ">= 1.6.0"
}
`
	err := os.WriteFile("main.tf", []byte(hclContent), 0644)
	if err != nil {
		t.Fatal("Failed to create main.tf file:", err)
	}
	defer os.Remove("main.tf")

	cmd := exec.Command(tenvBin, "tofu", "detect")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
	require.NoError(t, err, "Expected no error during version detection")
	require.Contains(t, out.String(), "1.6.0")
}

func TestTofuVersionFromJSON(t *testing.T) {
	t.Parallel()

	tenvBin := os.Getenv("TENV_BIN")

	// Test .tf.json file
	jsonContent := `{
  "terraform": {
    "required_version": ">= 1.6.0"
  }
}`
	err := os.WriteFile("main.tf.json", []byte(jsonContent), 0644)
	if err != nil {
		t.Fatal("Failed to create main.tf.json file:", err)
	}
	defer os.Remove("main.tf.json")

	cmd := exec.Command(tenvBin, "tofu", "detect")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
	require.NoError(t, err, "Expected no error during version detection")
	require.Contains(t, out.String(), "1.6.0")
}

func TestAtmosVersionFromFile(t *testing.T) {
	t.Parallel()

	tenvBin := os.Getenv("TENV_BIN")

	// Test .atmos-version file
	err := os.WriteFile(".atmos-version", []byte("1.130.0"), 0644)
	if err != nil {
		t.Fatal("Failed to create .atmos-version file:", err)
	}
	defer os.Remove(".atmos-version")

	cmd := exec.Command(tenvBin, "atmos", "detect")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
	require.NoError(t, err, "Expected no error during version detection")
	require.Contains(t, out.String(), "1.130.0")
}

func TestAtmosVersionFromAsdf(t *testing.T) {
	t.Parallel()

	tenvBin := os.Getenv("TENV_BIN")

	// Test .tool-versions file
	toolVersions := `atmos 1.130.1
terraform 1.7.0
`
	err := os.WriteFile(".tool-versions", []byte(toolVersions), 0644)
	if err != nil {
		t.Fatal("Failed to create .tool-versions file:", err)
	}
	defer os.Remove(".tool-versions")

	cmd := exec.Command(tenvBin, "atmos", "detect")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
	require.NoError(t, err, "Expected no error during version detection")
	require.Contains(t, out.String(), "1.130.1")
}

func TestAtmosVersionFromEnv(t *testing.T) {
	t.Parallel()

	tenvBin := os.Getenv("TENV_BIN")

	cmd := exec.Command(tenvBin, "atmos", "detect")
	env := os.Environ()
	env = append(env, "ATMOS_VERSION=1.130.2")
	cmd.Env = env

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	require.NoError(t, err, "Expected no error during version detection")
	require.Contains(t, out.String(), "1.130.2")
}

func TestAtmosInstallAndUse(t *testing.T) {
	t.Parallel()

	tenvBin := os.Getenv("TENV_BIN")

	// Test installation of specific version
	cmd := exec.Command(tenvBin, "atmos", "install", "1.130.0")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	require.NoError(t, err, "Expected no error during installation")
	require.Contains(t, out.String(), "Installation of Atmos 1.130.0 successful")

	// Test using the installed version
	cmd = exec.Command(tenvBin, "atmos", "use", "1.130.0")
	out.Reset()
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
	require.NoError(t, err, "Expected no error when switching version")
	require.Contains(t, out.String(), "Now using Atmos 1.130.0")
}

func TestTerragruntVersionFromAsdf(t *testing.T) {
	t.Parallel()

	tenvBin := os.Getenv("TENV_BIN")

	// Test .tool-versions file
	toolVersions := `terragrunt 0.71.0
terraform 1.7.0
`
	err := os.WriteFile(".tool-versions", []byte(toolVersions), 0644)
	if err != nil {
		t.Fatal("Failed to create .tool-versions file:", err)
	}
	defer os.Remove(".tool-versions")

	cmd := exec.Command(tenvBin, "tg", "detect")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
	require.NoError(t, err, "Expected no error during version detection")
	require.Contains(t, out.String(), "0.71.0")
}

func TestTerragruntVersionFromEnv(t *testing.T) {
	t.Parallel()

	tenvBin := os.Getenv("TENV_BIN")

	cmd := exec.Command(tenvBin, "tg", "detect")
	env := os.Environ()
	env = append(env, "TG_VERSION=0.71.1")
	cmd.Env = env

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	require.NoError(t, err, "Expected no error during version detection")
	require.Contains(t, out.String(), "0.71.1")
}

func TestTerragruntVersionFromHCL(t *testing.T) {
	t.Parallel()

	tenvBin := os.Getenv("TENV_BIN")

	// Test terragrunt.hcl file with version constraint
	hclContent := `
terragrunt_version_constraint = "~> 0.71.0"

terraform {
  source = "..."
}
`
	err := os.WriteFile("terragrunt.hcl", []byte(hclContent), 0644)
	if err != nil {
		t.Fatal("Failed to create terragrunt.hcl file:", err)
	}
	defer os.Remove("terragrunt.hcl")

	cmd := exec.Command(tenvBin, "tg", "detect")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
	require.NoError(t, err, "Expected no error during version detection")
	require.Contains(t, out.String(), "0.71.0")
}

func TestTerragruntVersionFromJSON(t *testing.T) {
	t.Parallel()

	tenvBin := os.Getenv("TENV_BIN")

	// Test terragrunt.json file with version constraint
	jsonContent := `{
  "terragrunt_version_constraint": "~> 0.71.0",
  "terraform": {
    "source": "..."
  }
}`
	err := os.WriteFile("terragrunt.json", []byte(jsonContent), 0644)
	if err != nil {
		t.Fatal("Failed to create terragrunt.json file:", err)
	}
	defer os.Remove("terragrunt.json")

	cmd := exec.Command(tenvBin, "tg", "detect")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
	require.NoError(t, err, "Expected no error during version detection")
	require.Contains(t, out.String(), "0.71.0")
}

func TestTerragruntInstallAndUse(t *testing.T) {
	t.Parallel()

	tenvBin := os.Getenv("TENV_BIN")

	// Test installation of specific version
	cmd := exec.Command(tenvBin, "tg", "install", "0.71.0")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	require.NoError(t, err, "Expected no error during installation")
	require.Contains(t, out.String(), "Installation of Terragrunt 0.71.0 successful")

	// Test using the installed version
	cmd = exec.Command(tenvBin, "tg", "use", "0.71.0")
	out.Reset()
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
	require.NoError(t, err, "Expected no error when switching version")
	require.Contains(t, out.String(), "Now using Terragrunt 0.71.0")
}

func TestTerraformVersionFromFile(t *testing.T) {
	t.Parallel()

	tenvBin := os.Getenv("TENV_BIN")

	// Test .terraform-version file
	err := os.WriteFile(".terraform-version", []byte("1.7.0"), 0644)
	if err != nil {
		t.Fatal("Failed to create .terraform-version file:", err)
	}
	defer os.Remove(".terraform-version")

	cmd := exec.Command(tenvBin, "tf", "detect")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
	require.NoError(t, err, "Expected no error during version detection")
	require.Contains(t, out.String(), "1.7.0")
}

func TestTerraformVersionFromTfswitchrc(t *testing.T) {
	t.Parallel()

	tenvBin := os.Getenv("TENV_BIN")

	// Test .tfswitchrc file
	err := os.WriteFile(".tfswitchrc", []byte("1.7.1"), 0644)
	if err != nil {
		t.Fatal("Failed to create .tfswitchrc file:", err)
	}
	defer os.Remove(".tfswitchrc")

	cmd := exec.Command(tenvBin, "tf", "detect")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
	require.NoError(t, err, "Expected no error during version detection")
	require.Contains(t, out.String(), "1.7.1")
}

func TestTerraformVersionFromAsdf(t *testing.T) {
	t.Parallel()

	tenvBin := os.Getenv("TENV_BIN")

	// Test .tool-versions file
	toolVersions := `terraform 1.7.2
terragrunt 0.71.0
`
	err := os.WriteFile(".tool-versions", []byte(toolVersions), 0644)
	if err != nil {
		t.Fatal("Failed to create .tool-versions file:", err)
	}
	defer os.Remove(".tool-versions")

	cmd := exec.Command(tenvBin, "tf", "detect")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
	require.NoError(t, err, "Expected no error during version detection")
	require.Contains(t, out.String(), "1.7.2")
}

func TestTerraformVersionFromTerragruntHCL(t *testing.T) {
	t.Parallel()

	tenvBin := os.Getenv("TENV_BIN")

	// Test terragrunt.hcl file with terraform version constraint
	hclContent := `
terraform {
  source = "..."

  # Configure Terraform version constraint
  terraform_version_constraint = "~> 1.7.0"
}
`
	err := os.WriteFile("terragrunt.hcl", []byte(hclContent), 0644)
	if err != nil {
		t.Fatal("Failed to create terragrunt.hcl file:", err)
	}
	defer os.Remove("terragrunt.hcl")

	cmd := exec.Command(tenvBin, "tf", "detect")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
	require.NoError(t, err, "Expected no error during version detection")
	require.Contains(t, out.String(), "1.7.0")
}

func TestTerraformVersionFromHCL(t *testing.T) {
	t.Parallel()

	tenvBin := os.Getenv("TENV_BIN")

	// Test .tf file with required version
	hclContent := `
terraform {
  required_version = ">= 1.7.0"
}
`
	err := os.WriteFile("main.tf", []byte(hclContent), 0644)
	if err != nil {
		t.Fatal("Failed to create main.tf file:", err)
	}
	defer os.Remove("main.tf")

	cmd := exec.Command(tenvBin, "tf", "detect")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
	require.NoError(t, err, "Expected no error during version detection")
	require.Contains(t, out.String(), "1.7.0")
}

func TestTerraformVersionFromJSON(t *testing.T) {
	t.Parallel()

	tenvBin := os.Getenv("TENV_BIN")

	// Test .tf.json file with required version
	jsonContent := `{
  "terraform": {
    "required_version": ">= 1.7.0"
  }
}`
	err := os.WriteFile("main.tf.json", []byte(jsonContent), 0644)
	if err != nil {
		t.Fatal("Failed to create main.tf.json file:", err)
	}
	defer os.Remove("main.tf.json")

	cmd := exec.Command(tenvBin, "tf", "detect")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
	require.NoError(t, err, "Expected no error during version detection")
	require.Contains(t, out.String(), "1.7.0")
}
