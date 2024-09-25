package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	NIX_SECRET_DIR = "/run/secrets"

	ENV_VAR_SECRET_ID = "SECRET_ID"
)

func listNixSecrets(secretsDir string) (secretNames []string, err error) {
	secretsDir, err = filepath.EvalSymlinks(secretsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve secrets dir: %w", err)
	}
	secretFiles, err := os.ReadDir(secretsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read secrets dir: %w", err)
	}
	for _, secretFile := range secretFiles {
		secretNames = append(secretNames, secretFile.Name())
	}
	return
}
