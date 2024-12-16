package internal

import (
	"fmt"
)

func PopulatePodmanSecretsDB(nixsecretsPath, mappingDirPath string, deletePodmanSecret DeletePodmanSecretFunc, createPodmanSecret CreatePodmanSecretFunc, debug bool) {
	debugLog(debug, "Listing nix secrets")
	if err := EnsureMappingDirExists(mappingDirPath); err != nil {
		panic(err)
	}
	nixSecretNames, err := ListNixSecrets(nixsecretsPath)
	if err != nil {
		panic(fmt.Errorf("failed to list nix secret names: %w", err))
	}
	debugLog(debug, "Listing podman secrets")
	podmanSecrets, removedSecretIDs, err := listPodmanSecrets(mappingDirPath)
	if err != nil {
		panic(fmt.Errorf("failed to list podman secrets: %w", err))
	}

	for _, secretID := range removedSecretIDs {
		debugLog(debug, "Deleting secret with id %s from podman", secretID)
		if err := deletePodmanSecret(secretID); err != nil {
			panic(fmt.Errorf("failed to remove secret by id (%s): %w", secretID, err))
		}
	}

	// Check if we need to remove secrets
	for _, secretName := range podmanSecrets {
		if !sliceContains(nixSecretNames, secretName) {
			debugLog(debug, "Deleting secrets %s from podman", secretName)
			if err := deletePodmanSecret(secretName); err != nil {
				panic(fmt.Errorf("failed to delete secret %s: %w", secretName, err))
			}
		}
	}

	// Create missing secrets
	for _, secretName := range nixSecretNames {
		if !sliceContains(podmanSecrets, secretName) {
			debugLog(debug, "Creating secret %s in podman", secretName)
			if err := createPodmanSecret(secretName); err != nil {
				panic(fmt.Errorf("failed to create secret %s: %w", secretName, err))
			}
		}
	}
	debugLog(debug, "Finished syncing nix secrets to podman")
}

func debugLog(debug bool, message string, vals ...interface{}) {
	if debug {
		fmt.Printf(message+"\n", vals...)
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
