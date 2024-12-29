# sovr.io

This is the whole system for sovr.io, which includes various tools for self-hosters.

## Nix

`nix build` should work and produce output at `result/bin/sovr.io`

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
