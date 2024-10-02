package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	NIX_SECRET_DIR = "/run/secrets"
	MAPPING_SUBDIR = "podman-mapping"

	ENV_VAR_SECRET_ID = "SECRET_ID"
)

func EnsureMappingDirExists(nixSecretDir string) error {
	mappingDirPath := filepath.Join(nixSecretDir, MAPPING_SUBDIR)
	if stat, err := os.Stat(mappingDirPath); os.IsNotExist(err) {
		if err := os.Mkdir(mappingDirPath, 0700); err != nil {
			return fmt.Errorf("failed to create mapping dir: %w", err)
		}
	} else if !stat.IsDir() {
		return fmt.Errorf("mapping dir path %s exists, but is not a directory", mappingDirPath)
	}
	return nil
}

func ListNixSecrets(secretsDir string) (secretNames []string, err error) {
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
