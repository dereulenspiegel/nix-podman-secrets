package main

import (
	"bytes"
	"strings"
	"testing"
)

var secretsDir = "./test-secrets-dir"

func TestListSecrets(t *testing.T) {

	secrets, err := listNixSecrets(secretsDir)
	if err != nil {
		t.Fatalf("failed to read nix secrets: %s", err)
	}

	expectedList := "secret1\nsecret2"

	if expectedList != strings.Join(secrets, "\n") {
		t.Fatal("secret list does not match expected format")
	}
}

func TestLookupSecret(t *testing.T) {
	out := &bytes.Buffer{}
	secretId := "secret1"
	lookupSecret(out, secretsDir, secretId)

	expectedSecretData := "foo"

	if expectedSecretData != out.String() {
		t.Fatalf("Looked up secrets data does not match, expected '%s', got '%s'", expectedSecretData, out.String())
	}
}
