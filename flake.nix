{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.11";

    flake-parts.url = "github:hercules-ci/flake-parts";

    alejandra = {
      url = "github:kamadorueda/alejandra/3.0.0";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    rust-overlay = {
      url = "github:oxalica/rust-overlay";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = inputs @ {self, ...}:
    inputs.flake-parts.lib.mkFlake {inherit inputs;} (toplevel @ {withSystem, ...}: {
      systems = ["aarch64-darwin" "aarch64-linux" "x86_64-linux"];

      perSystem = {
        config,
        self',
        inputs',
        pkgs,
        system,
        ...
      }: {
        _module.args.pkgs = import inputs.nixpkgs {
          localSystem = system;
          config = {
            allowUnfree = true;
            allowAliases = true;
          };
          overlays = [inputs.rust-overlay.overlays.default];
        };

        devShells = import ./nix/devshells.nix pkgs;

        packages = import ./nix/package-sets.nix pkgs;

        formatter = inputs'.alejandra.packages.default;
      };
    });
}
