# nix-podman-secrets

This is a very simple program and flake to configure podman on a NixOS system
to use the secrets in `/run/secrets` populated i.e. by [sops-nix](https://github.com/Mic92/sops-nix).
Since these secrets are populated and managed by Nix it is not possible for podman to
create new secrets or delete existing secrets.
