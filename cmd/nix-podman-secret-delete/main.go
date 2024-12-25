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
		deleteSecret(internal.MAPPING_DIR, secretId)
	})
}

func deleteSecret(mappingDir, secretId string) {
	if err := os.Remove(filepath.Join(mappingDir, secretId)); err != nil {
		panic(fmt.Errorf("failed to remove mapping symlink for secret id %s: %w", secretId, err))
	}
}
