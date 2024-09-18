inputs: flake-self:
let
  inherit inputs;
  inherit flake-self;
in
self: super: {
  nix-podman-secrets = super.callPackage ./package.nix;
}
