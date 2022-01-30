with import <nixpkgs> {};

stdenv.mkDerivation {

  name = "auth-api";

  buildInputs = with pkgs; [
    jetbrains.goland
  ];

  shellHook = ''
    goland .
    exit
  '';
}

