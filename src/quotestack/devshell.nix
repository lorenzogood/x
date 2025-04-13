{
  mkShell,
  pkgs,
  ...
}:
mkShell {
  buildInputs = with pkgs; [go go-tools gotools gopls];
}
