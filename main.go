package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	CMD_LIST   = "list"
	CMD_LOOKUP = "lookup"
	CMD_NOOP   = "noop"

	NIX_SECRET_DIR = "/run/secrets"

	ENV_VAR_SECRET_ID = "SECRET_ID"
)

func main() {
	cmdName := os.Args[1]

	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				fmt.Fprintf(os.Stderr, "an error occured: %s", err)
			} else {
				fmt.Fprintf(os.Stderr, "something unexpected happened (%T), %s", r, r)
			}
			os.Exit(1)
		}
	}()

	switch cmdName {
	case CMD_LIST:
		listSecrets()
	case CMD_LOOKUP:
		lookupSecret()
	case CMD_NOOP:
		fmt.Fprint(os.Stderr, "write access to nix managed secrets is not possible")
		os.Exit(1)
	default:
		panic(fmt.Errorf("unsupported command %s", cmdName))
	}
}

func lookupSecret() {
	secretId := os.Getenv(ENV_VAR_SECRET_ID)
	if secretId == "" {
		panic(errors.New("no SECRET_ID given for lookup"))
	}
	secretFilePath := filepath.Join(NIX_SECRET_DIR, secretId)
	secretBytes, err := os.ReadFile(secretFilePath)
	if err != nil {
		panic(fmt.Errorf("failed to read secret data from filesystem: %s", err))
	}
	fmt.Print(string(secretBytes))
}

func listSecrets() {
	secretFiles, err := os.ReadDir(NIX_SECRET_DIR)
	if err != nil {
		panic(fmt.Errorf("can't list nix secrets: %s", err))
	}
	for _, secretFile := range secretFiles {
		fmt.Println(secretFile.Name())
	}
}
