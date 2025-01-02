# sovr.cloud

This is the whole system for sovr.cloud, which includes various tools for self-hosters.

## Nix

`nix build` should work and produce output at `result/bin/sovr.cloud`

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
