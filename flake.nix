{
  description = "Sovr.cloud server";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.11";
  };

  outputs = { self, nixpkgs }:
    let
      allSystems = [
        "x86_64-linux" # 64-bit Intel/AMD Linux
        "aarch64-linux" # 64-bit ARM Linux
        "x86_64-darwin" # 64-bit Intel macOS
        "aarch64-darwin" # 64-bit ARM macOS
      ];
      forAllSystems = f: nixpkgs.lib.genAttrs allSystems (system: f {
        inherit system;
        pkgs = import nixpkgs { inherit system; };
      });
    in
    {
      nixosModules.default = { config, pkgs, lib, ... }: {
        options = {
          services.sovr-server = {
            enable = lib.mkEnableOption "sovr-server service";
            port = lib.mkOption {
              type = lib.types.port;
              default = 8080;
              description = "the port to serve requests on";
            };
            sessionSecret = lib.mkOption {
              type = lib.types.str;
              default = "please-change-me";
              description = "the secret used for generating session cookies";
            };
          };
        };

        config = lib.mkIf config.services.sovr-server.enable {
          systemd.services.sovr-server = {
            after = [ "network.target" ];
            enable = true;
            environment = {
              GIN_MODE = "release";
              PORT = "${toString config.services.sovr-server.port}";
              RELYING_PARTY_DISPLAY_NAME = "Sovr Cloud";
              RELYING_PARTY_ID = "www.sovr.cloud";
              RELYING_PARTY_ORIGINS = "https://www.sovr.cloud";
              SESSION_SECRET = "${toString config.services.sovr-server.sessionSecret}";
              TEMPLATES_DIR = "/opt/sovr.cloud/templates";
              TRUSTED_PROXIES = "127.0.0.1,::1";
            };
            serviceConfig = {
              ExecStart = "${self.packages.${pkgs.system}.default}/bin/sovr.cloud";
              Group = "deploy";
              Restart = "on-failure";
              Type = "simple";
              User = "deploy";
              WorkingDirectory = "/opt/sovr.cloud";
            };
            wantedBy = [ "multi-user.target" ];
          };
        };
      };
      packages = forAllSystems ({ pkgs, system }: {
        default = self.packages.${system}.sovr-server;
        sovr-server = pkgs.buildGoModule rec {
          overrideModAttrs = (oldAttrs: {
            preBuild = /* bash */ ''
              export GOPROXY=https://goproxy.io
            '';
          });
          pname = "sovr-server";
          src = ./.;
          vendorHash = "sha256-YOVPsS3iWg7P2awcvfRPT+cPOpJBfw6T2IzIJTseq+k=";
          version = "1.0.0";
        };
      });
    };
}
