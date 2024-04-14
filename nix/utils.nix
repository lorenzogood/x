{
  # Get all subdirectories of basedir that contain a filename.
  getSubDirFiles = basedir: filename: let
    entries = builtins.readDir basedir;

    procEntry = name: type: let
      path = basedir + "/${name}";
    in
      if type == "directory"
      then
        (
          if builtins.pathExists (path + "/${filename}")
          then [path]
          else []
        )
      else [];
  in
    builtins.concatLists (builtins.attrValues (builtins.mapAttrs procEntry entries));

  # Build a derivation from call package.
  buildCallPackage = pkgs: filename: dir: {
    name = builtins.baseNameOf (toString dir);
    value = pkgs.callPackage (dir + "/${filename}") {};
  };

  # Simple foldl usage to concat attrsets.
  concatAttrsets = sets: builtins.foldl' (acc: set: acc // set) {} sets;

  # Unfold an attrset into a list like { a = 1} -> [{ name = "a"; value = 1}].
  reverseAttrToList = set:
    builtins.map (name: {
      inherit name;
      value = set.${name};
    }) (builtins.attrNames set);
}
