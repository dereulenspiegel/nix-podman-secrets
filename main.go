package main

import (
	"errors"
	"fmt"
	"io"
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

	if len(os.Args) < 2 {
		panic(errors.New("you need to specify one of the subcommands list lookup or noop"))
	}
	cmdName := os.Args[1]

	switch cmdName {
	case CMD_LIST:
		listSecrets(os.Stdout, NIX_SECRET_DIR)
	case CMD_LOOKUP:
		secretId := os.Getenv(ENV_VAR_SECRET_ID)
		lookupSecret(os.Stdout, NIX_SECRET_DIR, secretId)
	case CMD_NOOP:
		fmt.Fprint(os.Stderr, "write access to nix managed secrets is not possible")
		os.Exit(1)
	default:
		panic(fmt.Errorf("unsupported command %s", cmdName))
	}
}

func lookupSecret(w io.Writer, secretDir, secretId string) {
	if secretId == "" {
		panic(errors.New("no SECRET_ID given for lookup"))
	}
	secretFilePath := filepath.Join(secretDir, secretId)
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

func listSecrets(w io.Writer, secretsDir string) {
	secretsDir, err := filepath.EvalSymlinks(secretsDir)
	if err != nil {
		panic(fmt.Errorf("failed to resolve secrets dir: %s", err))
	}
	secretFiles, err := os.ReadDir(secretsDir)
	if err != nil {
		panic(fmt.Errorf("can't list nix secrets: %s", err))
	}
	for _, secretFile := range secretFiles {
		fmt.Fprintln(w, secretFile.Name())
	}
}
