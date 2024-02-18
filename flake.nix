{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
  };

  outputs = {
    self,
    nixpkgs,
    flake-parts,
  } @ inputs:
    flake-parts.lib.mkFlake {inherit inputs;} {
      systems = [
        "x86_64-linux"
      ];

      perSystem = {
        config,
        system,
        ...
      }: let
        pkgs = nixpkgs.legacyPackages.${system};
      in {
        packages.terraform-lockfile-insights = pkgs.hello;
        packages.default = self.packages.${system}.terraform-lockfile-insights;
      };
    };
}
