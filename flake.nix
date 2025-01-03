{
  description = "Use nix secrets in podman";

  inputs = { nixpkgs = { url = "github:nixos/nixpkgs/nixos-24.05"; }; };

  outputs = inputs@{ nixpkgs, self, ... }:
    let
      systems = [ "x86_64-linux" "aarch64-linux" ];

      forAllSystems = f: nixpkgs.lib.genAttrs systems (system: f system);

      nixpkgsFor = forAllSystems (system:
        import nixpkgs {
          inherit system;
          overlays = [ self.overlays.default ];
        });
    in {

      packages = forAllSystems (system:
        with nixpkgsFor.${system}; {
          inherit nix-podman-secrets;
          default = nix-podman-secrets;
        });

      overlays.default = final: prev:
        (import ./overlay.nix inputs self) final prev;

      nixosModules.default = (self:
        { lib, config, pkgs, ... }: {

          options.nix-podman-secrets = {
            podmanPackage = lib.mkOption {
              type = lib.types.package;
              default = pkgs.podman;
              description = "The podman package to use";
            };
          };

          config.environment.systemPackages =
            [ self.packages.${pkgs.system}.default ];

          config.system.activationScripts.syncNixPodmanSecrets =
            (lib.stringAfter ([
              "specialfs"
              "users"
              "groups"
              "setupSecrets"
            ])) ''
              [ -e /run/current-system ] || echo "populating podman secrets from nix secrets"
              PATH=$PATH:${
                lib.makeBinPath [
                  config.nix-podman-secrets.podmanPackage
                  self.packages.${pkgs.system}.nix-podman-secrets
                ]
              } ${
                self.packages.${pkgs.system}.nix-podman-secrets.outPath
              }/bin/nix-podman-secret-populate
            '';

        }) self;
    };
}
