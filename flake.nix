{
  description = "cmdguard — validated Cobra CLI library with type-safe dependency injection";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };
    systems.url = "github:nix-systems/default";
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs@{ self, flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import inputs.systems;
      imports = [ inputs.treefmt-nix.flakeModule ];

      perSystem =
        { config, pkgs, ... }:
        let
          goPkg = pkgs.go_1_26;
        in
        {
          treefmt = {
            projectRootFile = "go.mod";
            programs = {
              nixfmt.enable = true;
              gofmt.enable = true;
              gofumpt.enable = true;
              goimports.enable = true;
            };
          };

          devShells = {
            default = pkgs.mkShellNoCC {
              packages = [
                goPkg
                pkgs.gopls
                pkgs.golangci-lint
              ];

              GOWORK = "off";
              GOEXPERIMENT = "jsonv2";

              shellHook = ''
                export GOTMPDIR="$HOME/.cache/go-tmp"
                mkdir -p "$GOTMPDIR"
                echo "cmdguard dev shell — Go $(go version | awk '{print $3}')"
              '';
            };

            ci = pkgs.mkShellNoCC {
              packages = [
                goPkg
                pkgs.golangci-lint
              ];

              GOWORK = "off";
              GOEXPERIMENT = "jsonv2";

              shellHook = ''
                export GOTMPDIR="$HOME/.cache/go-tmp"
                mkdir -p "$GOTMPDIR"
              '';
            };
          };

          checks.format = config.treefmt.build.check self;
        };
    };
}
