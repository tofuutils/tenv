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
