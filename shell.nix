{ pkgs ? import <nixpkgs> { } }:

with pkgs;

mkShell {
  buildInputs = [
    go
    syft
    grype
    docker
    docker-credential-helpers    
    trivy
  ];
}
