package main

import (
	"bytes"
	"testing"
)

var secretsDir = "./test-secrets-dir"

func TestListSecrets(t *testing.T) {

	out := &bytes.Buffer{}
	listSecrets(out, secretsDir)

	expectedList := "secret1\nsecret2\n"

	if expectedList != out.String() {
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

func TestMockStore(t *testing.T) {
	in := &bytes.Buffer{}
	in.Write([]byte("irrelevant, secret is managed by nix"))

	mockStore(in, secretsDir, "secret1")
}
