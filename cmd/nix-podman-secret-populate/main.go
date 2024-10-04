package main

import (
	"os"

	"github.com/dereulenspiegel/nix-podman-secrets/internal"
)

const (
	ENV_DEBUG = "NIX_PODMAN_SECRETS_DEBUG"
)

func main() {
	debug := false
	if os.Getenv(ENV_DEBUG) == "true" {
		debug = true
	}
	internal.WrapMain(func() {
		internal.PopulatePodmanSecretsDB(
			internal.NIX_SECRET_DIR,
			internal.MAPPING_DIR,
			internal.DeletePodmanSecretImpl,
			internal.CreatePodmanSecretImpl,
			debug)
	})
}
