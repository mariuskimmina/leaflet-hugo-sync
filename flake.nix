{
  description = "Sync Leaflet/Bluesky blog posts to Hugo-compatible markdown";

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
        packages = {
          default = pkgs.buildGoModule {
            pname = "leaflet-hugo-sync";
            version = "0.1.0";

            src = ./.;

            vendorHash = "sha256-ruru71ASUBsvOByM58X0eoJTfEhtdijpt9d2g6DXnvI=";

            subPackages = [ "cmd/leaflet-hugo-sync" ];

            ldflags = [
              "-s"
              "-w"
            ];

            meta = with pkgs.lib; {
              description = "Sync Leaflet (ATproto) blog posts to Hugo-compatible markdown";
              homepage = "https://github.com/mariuskimmina/leaflet-hugo-sync";
              license = licenses.mit;
              mainProgram = "leaflet-hugo-sync";
            };
          };

          leaflet-hugo-sync = self.packages.${system}.default;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
            go-tools
            delve
          ];

          shellHook = ''
            echo "leaflet-hugo-sync development environment"
            echo "Go version: $(go version)"
            echo "Available commands:"
            echo "  go build -o leaflet-hugo-sync ./cmd/leaflet-hugo-sync"
          '';
        };

        apps.default = {
          type = "app";
          program = "${self.packages.${system}.default}/bin/leaflet-hugo-sync";
        };
      }
    );
}
