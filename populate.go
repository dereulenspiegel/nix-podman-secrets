package main

import (
	"fmt"
)

func populatePodmanSecretsDB(nixsecretsPath string, debug bool) {
	debugLog(debug, "Listing nix secrets")
	nixSecretNames, err := listNixSecrets(nixsecretsPath)
	if err != nil {
		panic(fmt.Errorf("failed to list nix secret names: %w", err))
	}
	debugLog(debug, "Listing podman secrets")
	podmanSecrets, err := listPodmanSecrets()
	if err != nil {
		panic(fmt.Errorf("failed to list podman secrets: %w", err))
	}

	// Check if we need to remove secrets
	for _, secretName := range podmanSecrets {
		if !sliceContains(nixSecretNames, secretName) {
			debugLog(debug, "Deleting secrets %s from podman", secretName)
			deletePodmanSecret(secretName)
		}
	}

	// Create missing secrets
	for _, secretName := range nixSecretNames {
		if !sliceContains(podmanSecrets, secretName) {
			debugLog(debug, "Creating secret %s in podman", secretName)
			createPodmanSecret(secretName)
		}
	}
	debugLog(debug, "Finished syncing nix secrets to podman")
}

func debugLog(debug bool, message string, vals ...interface{}) {
	if debug {
		fmt.Printf(message, vals...)
	}
}

func sliceContains[T comparable](slice []T, elem T) bool {
	for _, e := range slice {
		if e == elem {
			return true
		}
	}
	return false
}
