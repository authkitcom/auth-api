with import <nixpkgs> {};

stdenv.mkDerivation {

  name = "demo";

  buildInputs = with pkgs; [
    go_1_17
    protobuf
    go-protobuf
  ];

  shellHook = ''
    export GOPATH=$HOME/go
    export PATH=$PATH:$HOME/go/bin
    export GOPRIVATE=gitlab.authkit.com/*
    '';

  hardeningDisable = [ "fortify" ];

}

