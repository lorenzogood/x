pkgs: let
  inherit (import ./utils.nix) getSubDirFiles reverseAttrToList;

  build = path: let
    name = builtins.baseNameOf (toString path);
    package-sets = (import (path + "/packages.nix")) pkgs;
  in (builtins.map (set: {
    name = name + "-" + set.name;
    value = set.value;
  }) (reverseAttrToList package-sets));
in (builtins.listToAttrs (builtins.concatLists (builtins.map build (getSubDirFiles ../src "packages.nix"))))
