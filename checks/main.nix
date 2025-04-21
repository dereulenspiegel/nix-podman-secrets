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

        system.stateVersion = "24.11";
      };
  };

  testScript = ''
    start_all()

  '';
}
