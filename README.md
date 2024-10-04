# nix-podman-secrets

This is a very simple program and flake to configure podman on a NixOS system
to use the secrets in `/run/secrets` populated i.e. by [sops-nix](https://github.com/Mic92/sops-nix).

To use this you can simply add this flake to your flakem i.e.:

```
inputs = {
    nix-podman-secrets = {
      url = "github:dereulenspiegel/nix-podman-secrets";

      inputs.nixpkgs.follows = "nixpkgs";
    };
}
```

and add the module to you nixosSystem module list, i.e.

```
    nixosConfigurations = {
        podman-host = nixpkgs.lib.nixosSystem {
          system = "x86_64-linux";
          modules = [
            inputs.nix-podman-secrets.nixosModules.default
          ]
    }
```

This will add the necessary packages and setup the activation scripts.
