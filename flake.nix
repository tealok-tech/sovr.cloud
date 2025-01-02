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
	  services.sovr-server.enable = lib.mkEnableOption "sovr-server";
	};

	config = lib.mkIf config.services.sovr-server.enable {
	  systemd.services.sovr-server = {
	    serviceConfig.ExecStart = "${self.packages.${pkgs.system}.default}/bin/main";
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
