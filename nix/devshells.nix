pkgs: let
  inherit (import ./utils.nix) getSubDirFiles buildCallPackage;
in
  builtins.listToAttrs (builtins.map (buildCallPackage pkgs "devshell.nix") (getSubDirFiles ../src "devshell.nix"))
