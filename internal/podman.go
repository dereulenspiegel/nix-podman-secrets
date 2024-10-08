package internal

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	podmanBin = "podman"

	nixPodmanSecretsBin = "nix-podman-secret"
)

type DeletePodmanSecretFunc func(string) error
type CreatePodmanSecretFunc func(string) error

func listPodmanSecrets(mappingDirPath string) (secretNames []string, err error) {

	files, err := os.ReadDir(mappingDirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to list entries in mapping dir: %w", err)
	}

	for _, secretFile := range files {
		secretPath := filepath.Join(mappingDirPath, secretFile.Name())
		secretNameBytes, err := os.ReadFile(secretPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read secret name from file %s: %w", secretPath, err)
		}
		secretNames = append(secretNames, strings.TrimSpace(string(secretNameBytes)))
	}
	return
}

func DeletePodmanSecretImpl(secretName string) error {
	cmd := exec.Command(podmanBin, "secret", "delete", secretName)
	errBuf := &bytes.Buffer{}
	cmd.Stderr = errBuf
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to delete secret (%s): %w", errBuf.String(), err)
	}
	return nil
}

func CreatePodmanSecretImpl(secretName string) error {
	cmd := exec.Command(podmanBin,
		"secret",
		"create",
		"--label", "source=nix",
		"--driver", "shell",
		"--driver-opts", fmt.Sprintf("delete=%s-delete,list=%s-list,lookup=%s-lookup,store=%s-store",
			nixPodmanSecretsBin,
			nixPodmanSecretsBin,
			nixPodmanSecretsBin,
			nixPodmanSecretsBin),
		secretName, "-")
	errBuff := &bytes.Buffer{}
	stdInBuff := bytes.NewBuffer([]byte(secretName))
	cmd.Stdin = stdInBuff
	cmd.Stderr = errBuff
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create secret (%s): %w", errBuff.String(), err)
	}
	return nil
}
