package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

const (
	podmanBin = "/run/current-system/sw/bin/podman"

	nixPodmanSecretsBin = "/run/current-system/sw/bin/nix-podman-secrets"

	podmanSecretPrefix = "nix_"
)

func listPodmanSecrets() (secretNames []string, err error) {
	cmd := exec.Command(podmanBin,
		"secret", "list", "--format", "{{ .Name }}")
	outBuf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	cmd.Stdout = outBuf
	cmd.Stderr = errBuf
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets (%s): %w", errBuf.String(), err)
	}
	outParts := strings.Split(outBuf.String(), "\n")
	for _, part := range outParts {
		secretName := strings.TrimSpace(part)
		if len(secretName) > 0 && strings.HasPrefix(secretName, podmanSecretPrefix) {
			secretNames = append(secretNames, strings.TrimPrefix(secretName, podmanSecretPrefix))
		}
	}
	return
}

func deletePodmanSecret(secretName string) error {
	cmd := exec.Command(podmanBin, "secret", "delete", fmt.Sprintf("%s%s", podmanSecretPrefix, secretName))
	errBuf := &bytes.Buffer{}
	cmd.Stderr = errBuf
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to delete secret (%s): %w", errBuf.String(), err)
	}
	return nil
}

func createPodmanSecret(secretName string) error {
	cmd := exec.Command(podmanBin,
		"secret",
		"create",
		"--label", "source=nix",
		"--driver", "shell",
		"driver-opts", fmt.Sprintf("delete='%s noop',list='%s noop',lookup='%s lookup %s',store='%s noop'",
			nixPodmanSecretsBin,
			nixPodmanSecretsBin,
			nixPodmanSecretsBin,
			secretName,
			nixPodmanSecretsBin),
		fmt.Sprintf("%s%s", podmanSecretPrefix, secretName), "-")
	errBuff := &bytes.Buffer{}
	stdInBuff := bytes.NewBuffer([]byte("mock"))
	cmd.Stdin = stdInBuff
	cmd.Stderr = errBuff
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create secret (%s): %w", errBuff.String(), err)
	}
	return nil
}
