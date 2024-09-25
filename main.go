package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	CMD_LOOKUP   = "lookup"
	CMD_NOOP     = "noop"
	CMD_POPULATE = "populate"
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
	case CMD_LOOKUP:
		if len(os.Args) < 3 {
			panic(errors.New("secret name argument is missing"))
		}
		secretId := os.Args[2]
		lookupSecret(os.Stdout, NIX_SECRET_DIR, secretId)
	case CMD_NOOP:
		noop()
	case CMD_POPULATE:
		populatePodmanSecretsDB(NIX_SECRET_DIR)
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

func noop() {
	os.Exit(0)
}
