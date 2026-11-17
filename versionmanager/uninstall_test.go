package versionmanager

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/tofuutils/tenv/v4/config"
	terraformretriever "github.com/tofuutils/tenv/v4/versionmanager/retriever/terraform"
)

func errsFromJoin(err error) []error {
	type joinError interface{ Unwrap() []error }

	var jerr joinError
	if !errors.As(err, &jerr) {
		return nil
	}

	return jerr.Unwrap()
}

func buildManager(t *testing.T, toolFolder string) VersionManager {
	t.Helper()

	conf, confErr := config.DefaultConfig()
	if confErr != nil {
		t.Fatal(confErr)
	}

	conf.RootPath = t.TempDir()
	conf.LockPath = conf.RootPath
	conf.InitDisplayer(false)

	return Make(&conf, "", toolFolder, nil, terraformretriever.Make(&conf), nil)
}

func TestUninstall(t *testing.T) {
	t.Parallel()

	t.Run("Uninstall_not_installed_exact_version", func(t *testing.T) {
		t.Parallel()

		mgr := buildManager(t, "Terraform")
		uninstallErr := mgr.Uninstall("1.0.0")
		if uninstallErr == nil {
			t.Fatal("expected error when uninstalling a version that is not installed, got nil")
		}

		if !errors.Is(uninstallErr, errVersionNotInstalled) {
			t.Errorf("expected errVersionNotInstalled, got: %v", uninstallErr)
		}
	})

	t.Run("UninstallMultiple_returns_error_for_not_installed", func(t *testing.T) {
		t.Parallel()

		mgr := buildManager(t, "Terraform")
		uninstallErr := mgr.UninstallMultiple([]string{"1.0.0"})
		if uninstallErr == nil {
			t.Fatal("expected error when uninstalling a version that is not installed, got nil")
		}

		if !errors.Is(uninstallErr, errVersionNotInstalled) {
			t.Errorf("expected errVersionNotInstalled, got: %v", uninstallErr)
		}
	})

	t.Run("UninstallMultiple_returns_joined_errors_for_multiple_not_installed", func(t *testing.T) {
		t.Parallel()

		mgr := buildManager(t, "OpenTofu")
		uninstallErr := mgr.UninstallMultiple([]string{"1.0.0", "1.1.0"})
		if uninstallErr == nil {
			t.Fatal("expected error when uninstalling versions that are not installed, got nil")
		}

		unwrapped := errsFromJoin(uninstallErr)
		if len(unwrapped) == 0 {
			t.Fatal("expected joined error with at least one wrapped error")
		}

		found := false
		for _, u := range unwrapped {
			if errors.Is(u, errVersionNotInstalled) {
				found = true

				break
			}
		}

		if !found {
			t.Errorf("expected errVersionNotInstalled wrapped in joined error, got: %v", uninstallErr)
		}
	})

	t.Run("Uninstall_installed_version_succeeds", func(t *testing.T) {
		t.Parallel()

		mgr := buildManager(t, "Terraform")
		installPath, pathErr := mgr.InstallPath()
		if pathErr != nil {
			t.Fatal(pathErr)
		}

		versionDir := filepath.Join(installPath, "1.2.3")
		if mkdirErr := os.MkdirAll(versionDir, 0o755); mkdirErr != nil {
			t.Fatal(mkdirErr)
		}

		uninstallErr := mgr.Uninstall("1.2.3")
		if uninstallErr != nil {
			t.Fatalf("unexpected error uninstalling installed version: %v", uninstallErr)
		}

		if _, statErr := os.Stat(versionDir); !os.IsNotExist(statErr) {
			t.Error("expected version directory to be removed after uninstall")
		}
	})

	t.Run("Uninstall_same_version_twice_returns_error", func(t *testing.T) {
		t.Parallel()

		mgr := buildManager(t, "Terragrunt")
		installPath, pathErr := mgr.InstallPath()
		if pathErr != nil {
			t.Fatal(pathErr)
		}

		versionDir := filepath.Join(installPath, "0.49.0")
		if mkdirErr := os.MkdirAll(versionDir, 0o755); mkdirErr != nil {
			t.Fatal(mkdirErr)
		}

		firstErr := mgr.Uninstall("0.49.0")
		if firstErr != nil {
			t.Fatalf("first uninstall should succeed, got: %v", firstErr)
		}

		secondErr := mgr.Uninstall("0.49.0")
		if secondErr == nil {
			t.Fatal("expected error on second uninstall of already-removed version, got nil")
		}

		if !errors.Is(secondErr, errVersionNotInstalled) {
			t.Errorf("expected errVersionNotInstalled on second uninstall, got: %v", secondErr)
		}
	})
}
