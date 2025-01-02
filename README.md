# sovr.cloud

This is the whole system for sovr.cloud, which includes various tools for self-hosters.

## Nix

`nix build` should work and produce output at `result/bin/sovr.cloud`

Also if you want to make an update to a running NixOS server it helps to know that after committing you should do a 'sudo nix flake update sovr-server` before `sudo nixos-rebuild switch` to ensure you get the latest commits.

## Build

Install ninja. `nix-shell` should do this for you on NixOS.

```
ninja
```

## Run

```
./out/sovr-server
```

## Website

Should be running at [localhost:8080](localhost:8080/)
