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

      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      packages = forAllSystems (system:
        let pkgs = nixpkgsFor.${system};
        in {
          nix-podman-secrets = pkgs.buildGoModule {
            pname = "nix-podman-secrets";
            version = "0.1.0";

            src = ./.;
            doCheck = false;

            vendorHash = null;

            #outputs = [ "bin" ];

            meta = {
              description = "Simple tool for podman secrets shell driver to access nix secrets";
              homepage = "https://github.com/dereulenspiegel/nix-podman-secrets";
            };
          };
        });

      defaultPackage = forAllSystems (system: self.packages.${system}.nix-podman-secrets);

      defaultApp = forAllSystems (system: {
        type = "app";
        program = "${self.packages.${system}.nix-podman-secrets}/bin/nix-podman-secrets";
      });

      nixosModules.default = forAllSystems (system: { lib, config, ... }:
        with lib;
        let
          cfg = config.nix-podman-secrets;
        in
        {
          options.nix-podman-secrets = {
            enable = mkEnableOption "Enable nix-podman-secrets";
          };
          config = mkIf cfg.enable {
            systemPackages = [ self.packages.${system}.nix-podman-secrets ];

            environment.etc."containers/containers.conf.d/999_nix-podman-secrets.conf" = {
              enable = cfg.enable;
              text = ''
                [secrets]
                driver = "shell"

                [secrets.opts]
                list = /run/current-system/sw/bin/nix-podman-secrets list
                lookup = /run/current-system/sw/bin/nix-podman-secrets lookup
                store = /run/current-system/sw/bin/nix-podman-secrets noop
                delete = /run/current-system/sw/bin/nix-podman-secrets noop
              '';
            };
          };
        });
    };
}
