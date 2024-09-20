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

      nixosModules.default = (self: { lib, config, pkgs, ... }:
        let
          cfg = config.nix-podman-secrets;
        in
        {
          options.nix-podman-secrets = {
            enable = lib.mkEnableOption "Enable nix-podman-secrets";
          };

          config = lib.mkIf cfg.enable {
            environment.systemPackages = [ self.packages.${pkgs.system}.default ];

            environment.etc."containers/containers.conf.d/999_nix-podman-secrets.conf" = {
              enable = cfg.enable;
              text = ''
                [secrets]
                driver = "shell"

                [secrets.opts]
                list = "/run/current-system/sw/bin/nix-podman-secrets list"
                lookup = "/run/current-system/sw/bin/nix-podman-secrets lookup"
                store = "/run/current-system/sw/bin/nix-podman-secrets noop"
                delete = "/run/current-system/sw/bin/nix-podman-secrets noop"
              '';
            };
          };
        }) self;
    };
}
