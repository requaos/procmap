{
  description = "A terminal-based process visualization tool";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "procmap";
          version = "0.1.0";

          src = ./.;

          vendorHash = "sha256-8LxD44XZycnB2rLhwuR7EO0VbLGFI4zzjyyttrOsQzA=";

          meta = with pkgs.lib; {
            description = "Terminal-based process visualization tool";
            homepage = "https://github.com/req/procmap";
            license = licenses.mit;
            maintainers = [ ];
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            delve
            go-tools
          ];

          shellHook = ''
            echo "🚀 procmap development environment"
            echo ""
            echo "Tools available:"
            echo "  go $(go version | cut -d' ' -f3) - Go toolchain"
            echo "  gopls - Language server"
            echo "  delve - Debugger"
            echo "  staticcheck - Linter"
            echo ""
            echo "Quick commands:"
            echo "  go run main.go - Run the app"
            echo "  go test ./... - Run tests"
            echo "  nix build - Build with Nix"
            echo ""

            # Check if running nushell
            if [ -n "$NU_VERSION" ]; then
              echo "📝 Nushell detected! Run: source .nix-shell.nu"
              echo ""
            fi
          '';
        };
      }
    );
}
