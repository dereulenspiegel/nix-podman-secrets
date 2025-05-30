{
  description = "Use nix secrets in podman";

  inputs = {
    nixpkgs = { url = "github:nixos/nixpkgs/nixos-24.11"; };
    sops-nix = {
      url = "github:Mic92/sops-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = inputs@{ nixpkgs, self, sops-nix, ... }:
    let
      systems = [ "x86_64-linux" "aarch64-linux" ];
      checkSystems = [ "x86_64-linux" "aarch64-linux" "aarch64-darwin" ];

      forAllSystems = f: nixpkgs.lib.genAttrs systems (system: f system);
      forAllCheckSystems = f:
        nixpkgs.lib.genAttrs checkSystems (system: f system);

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

          imports = [ sops-nix.nixosModules.sops ];

          config.environment.systemPackages =
            [ self.packages.${pkgs.system}.default ];

          config.systemd.services.sync-nix-podman-secrets = {
            enable = true;
            before = [ "podman.service" ];
            wantedBy = [ "multi-user.target" ];
            path = [
              config.nix-podman-secrets.podmanPackage
              self.packages.${pkgs.system}.default
            ];

            serviceConfig = {
              Type = "oneshot";
              ExecStart = "${
                  self.packages.${pkgs.system}.nix-podman-secrets.outPath
                }/bin/nix-podman-secret-populate";
            };
          };

        }) self;

      checks = forAllCheckSystems (system:
        let
          checkArgs = {
            pkgs = nixpkgs.legacyPackages.${system};
            inherit self;
          };
        in { main = import ./checks/main.nix checkArgs; });
    };
}
