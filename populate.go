package main

import (
	"fmt"
)

func populatePodmanSecretsDB(nixsecretsPath string) {

	nixSecretNames, err := listNixSecrets(nixsecretsPath)
	if err != nil {
		panic(fmt.Errorf("failed to list nix secret names: %w", err))
	}
	podmanSecrets, err := listPodmanSecrets()
	if err != nil {
		panic(fmt.Errorf("failed to list podman secrets: %w", err))
	}

	// Check if we need to remove secrets
	for _, secretName := range podmanSecrets {
		if !sliceContains(nixSecretNames, secretName) {
			deletePodmanSecret(secretName)
		}
	}

	// Create missing secrets
	for _, secretName := range nixSecretNames {
		if !sliceContains(podmanSecrets, secretName) {
			createPodmanSecret(secretName)
		}
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
