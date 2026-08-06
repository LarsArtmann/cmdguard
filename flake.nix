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

          apps = {
            check-all = {
              type = "app";
              program = toString (pkgs.writeShellScript "check-all" ''
                set -euo pipefail
                export GOEXPERIMENT=jsonv2

                echo "=== Build ==="
                ${goPkg}/bin/go build ./...
                for mod in glamour prompts spinner telemetry flightrecorder; do
                  (cd "$mod" && ${goPkg}/bin/go build ./...)
                done

                echo "=== Test (race) ==="
                ${goPkg}/bin/go test ./... -count=1 -timeout 120s -race
                for mod in glamour prompts spinner telemetry flightrecorder; do
                  (cd "$mod" && ${goPkg}/bin/go test ./... -count=1 -timeout 120s -race)
                done

                echo "=== Lint ==="
                ${pkgs.golangci-lint}/bin/golangci-lint run ./...
                for mod in glamour prompts spinner telemetry flightrecorder; do
                  (cd "$mod" && ${pkgs.golangci-lint}/bin/golangci-lint run ./...)
                done

                echo "=== Format check ==="
                nix --extra-experimental-features 'nix-command flakes' flake check

                echo "=== go mod tidy check ==="
                export GOWORK=off
                for dir in . glamour prompts spinner telemetry flightrecorder; do
                  if [ -f "$dir/go.mod" ]; then
                    (cd "$dir" && ${goPkg}/bin/go mod tidy)
                    if ! (cd "$dir" && git diff --exit-code go.mod go.sum >/dev/null 2>&1); then
                      (cd "$dir" && git restore go.mod go.sum 2>/dev/null || true)
                      echo "FAIL: $dir/go.mod or go.sum not tidy. Run: (cd $dir && go mod tidy)"
                      exit 1
                    fi
                    echo "  OK: $dir/go.mod tidy"
                  fi
                done

                echo ""
                echo "All checks passed!"
              '');
            };
          };
        };
    };
}
