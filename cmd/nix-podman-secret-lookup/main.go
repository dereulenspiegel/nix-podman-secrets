package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dereulenspiegel/nix-podman-secrets/internal"
)

func main() {
	internal.WrapMain(func() {
		secretId := os.Getenv("SECRET_ID")
		lookupSecret(os.Stdout, internal.MAPPING_DIR, secretId)
	})
}

func lookupSecret(w io.Writer, mappingDirPath, secretId string) {
	if secretId == "" {
		panic(errors.New("no SECRET_ID given for lookup"))
	}
	secretFilePath := filepath.Join(mappingDirPath, secretId)
	secretFilePath, err := filepath.EvalSymlinks(secretFilePath)
	if err != nil {
		panic(fmt.Errorf("failed to resolve secrets dir: %s", err))
	}
	secretBytes, err := os.ReadFile(secretFilePath)
	if err != nil {
		panic(fmt.Errorf("failed to read secret data from filesystem: %s", err))
	}
	fmt.Fprint(w, string(secretBytes))
}
