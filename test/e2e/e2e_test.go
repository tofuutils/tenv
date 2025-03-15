//go:build e2e

package e2e

import (
	"bytes"
	"os/exec"
	"os"
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
