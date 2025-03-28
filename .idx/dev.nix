# To learn more about how to use Nix to configure your environment
# see: https://developers.google.com/idx/guides/customize-idx-env
{ pkgs, ... }: {
  # Which nixpkgs channel to use.
  channel = "stable-24.05"; # or "unstable"
  # Use https://search.nixos.org/packages to find packages
  packages = [
    pkgs.go
    pkgs.air
    pkgs.sqlc
    pkgs.gnumake42
    pkgs.go-swag
    pkgs.go-swagger
    pkgs.docker-compose
    pkgs.go-migrate
  ];
  # Sets environment variables in the workspace
  env = {};
  services = {
    docker = {
      enable = true;
    };
  };
  idx = {
    # Search for the extensions you want on https://open-vsx.org/ and use "publisher.id"
    extensions = [
      "golang.go"
      "rangav.vscode-thunder-client"
    ];
    workspace = {
      onCreate = {
        # Open editors for the following files by default, if they exist:
        default.openFiles = ["main.go"];
      };
      # Runs when a workspace is (re)started
      onStart= {
        run-server = "air";
      };
      # To run something each time the workspace is first created, use the `onStart` hook
    };
  };
}
