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
        pkgs = import nixpkgs { inherit system; };
      });
    in
    {
      packages = forAllSystems ({ pkgs }: {
        default = pkgs.buildGo123Module rec {
          pname = "sovr-server";
          version = "1.0.0";
          #subPackages = [ "sovr" ];
          src = ./.;
          vendorHash = "sha256-YOVPsS3iWg7P2awcvfRPT+cPOpJBfw6T2IzIJTseq+k=";
        };
      });
    };
}
