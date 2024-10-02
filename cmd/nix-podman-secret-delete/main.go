package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dereulenspiegel/nix-podman-secrets/internal"
)

func main() {
	internal.WrapMain(func() {
		secretId := os.Getenv("SECRET_ID")
		deleteSecret(internal.NIX_SECRET_DIR, secretId)
	})
}

func deleteSecret(nixSecretDir, secretId string) {
	pathToDelete, err := filepath.EvalSymlinks(filepath.Join(nixSecretDir, internal.MAPPING_SUBDIR, secretId))
	if err != nil {
		panic(fmt.Errorf("failed to evaluate mapping file path: %w", err))
	}
	if err := os.Remove(pathToDelete); err != nil {
		panic(fmt.Errorf("failed to delete secret mapping: %w", err))
	}
}
