package internal

import (
	"os"
	"path/filepath"
	"testing"
)

type podmanCmdMock struct {
	called bool
}

func (p *podmanCmdMock) cmd(in string) error {
	p.called = true
	return nil
}

func TestPopilateExistingSecrets(t *testing.T) {
	tmpNixSecretDir := filepath.Join(t.TempDir(), "nix-secrets")
	tmpMappingDir := filepath.Join(t.TempDir(), "mapping")
	err := os.MkdirAll(tmpNixSecretDir, 0700)
	if err != nil {
		t.Fatalf("failed to create test directories:%s", err)
	}
	err = os.MkdirAll(tmpMappingDir, 0700)
	if err != nil {
		t.Fatalf("failed to create test directories:%s", err)
	}

	if err := os.WriteFile(filepath.Join(tmpNixSecretDir, "secret_foo"), []byte("secret-value"), 0600); err != nil {
		t.Fatalf("failed to create nix secret: %s", err)
	}

	deletePodmanSecretMock := &podmanCmdMock{false}
	createPodmanSecretMock := &podmanCmdMock{false}

	PopulatePodmanSecretsDB(tmpNixSecretDir, tmpMappingDir, deletePodmanSecretMock.cmd, createPodmanSecretMock.cmd, false)

	if !createPodmanSecretMock.called {
		t.Fatalf("create podman secret was never called")
	}
}
