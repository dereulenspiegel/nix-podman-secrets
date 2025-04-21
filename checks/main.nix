(import ./lib.nix) {
  name = "main-server";
  nodes = {
    machine = { self, pkgs, ... }:
      let
        additionalSSHKey = pkgs.writeTextFile {
          name = "ssh.key";
          text = ''
            -----BEGIN OPENSSH PRIVATE KEY-----
            b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
            QyNTUxOQAAACDxbeKYjUHvMHgLbDCFVDTpgYaZ/gSDRBPNYoreGnDdrwAAAJg4/aWbOP2l
            mwAAAAtzc2gtZWQyNTUxOQAAACDxbeKYjUHvMHgLbDCFVDTpgYaZ/gSDRBPNYoreGnDdrw
            AAAEBdrJQzHoB31buNmiEI+WBLsTQ5zwS/ZF1BzPkkMAFWC/Ft4piNQe8weAtsMIVUNOmB
            hpn+BINEE81iit4acN2vAAAAFHRpbGxAdXR0ZXJseS1hYnlzbWFsAQ==
            -----END OPENSSH PRIVATE KEY-----

          '';
        };
      in {
        imports = [ self.nixosModules.default ];
        sops = {
          age = { sshKeyPaths = [ "${additionalSSHKey}" ]; };
          secrets = { test = { sopsFile = ./test.secret.yaml; }; };
        };

        virtualisation.containers.enable = true;
        virtualisation = {
          podman = {
            enable = true;

            # Create a `docker` alias for podman, to use it as a drop-in replacement
            dockerCompat = true;

            # Required for containers under podman-compose to be able to talk to each other.
            defaultNetwork.settings.dns_enabled = true;
          };
        };

        system.stateVersion = "24.11";
      };
  };

  testScript = ''
    start_all()
    machine.wait_for_unit("multi-user.target")
    machine.succeed("podman secret list")
    machine.succeed("podman secret exists test")
  '';
}
