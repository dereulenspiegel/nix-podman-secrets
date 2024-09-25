{
  description = "Use nix secrets in podman";

  inputs = {
    nixpkgs = {
      url = "github:nixos/nixpkgs/nixos-24.05";
    };
  };

  outputs = inputs@{ nixpkgs, self, ... }:
    let
      systems = [ "x86_64-linux" "aarch64-linux" ];

      forAllSystems = f: nixpkgs.lib.genAttrs systems (system: f system);

      nixpkgsFor = forAllSystems (system: import nixpkgs {
        inherit system;
        overlays = [ self.overlays.default ];
      });
    in
    {

      packages = forAllSystems (system: with nixpkgsFor.${system};{
        inherit nix-podman-secrets;
        default = nix-podman-secrets;
      });

      overlays.default = final: prev: (import ./overlay.nix inputs self) final prev;

      nixosModules.default = (self: { lib, config, pkgs, ... }: {
        environment.systemPackages = [ self.packages.${pkgs.system}.default ];

        system.activationScripts.syncNixPodmanSecrets = ''
          /run/current-system/sw/bin/nix-podman-secrets populate
        '';

      }) self;
    };
}
