{ lib, buildGoModule }:

buildGoModule {
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
}
