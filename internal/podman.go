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

func listPodmanSecrets(mappingDirPath string) (secretNames []string, removedSecretIDs []string, err error) {

	files, err := os.ReadDir(mappingDirPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list entries in mapping dir: %w", err)
	}

	for _, secretFile := range files {
		secretPath := filepath.Join(mappingDirPath, secretFile.Name())
		actualSecretFile, err := filepath.EvalSymlinks(secretPath)
		if err != nil {
			removedSecretIDs = append(removedSecretIDs, secretFile.Name())
			continue
		}
		secretName := filepath.Base(actualSecretFile)

		secretNames = append(secretNames, strings.TrimSpace(secretName))
	}
	return
}

func DeletePodmanSecretImpl(secretName string) error {
	cmd := exec.Command(podmanBin, "secret", "rm", secretName)
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
