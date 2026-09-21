---
title: Nixtoo Genux™ released
date: 2025-09-04
order: 041
tags: articles
source: "therealtruex.com (recovered from web.archive.org; date approximate, from the announced repo's creation)"
---

<a href="https://github.com/equwal/nixtoo-genux" class="external-link">git hub website link here</a>

## Nix package manager on another OS (no nix-env). 100% Declarative setup.

Including Gentoo of course. Nixtoo Genux!

### Files:

1.  `flake.nix`: Put packages here. Includes a dev shell and flakey.

2.  `aliasrc-additions.zsh`: Aliases to make it easier to use the setup. Very mandatory.

3.  `INSTALL.sh`: Installs this setup. Install it after installing the multi-user installation of nix from the <a href="https://nixos.org/download/" class="external-link">official website using curl into shell</a> or (for Gentoo) from `sys-apps/nix` from <a href="https://github.com/trofi/nix-guix-gentoo" class="external-link">Trofi's nix-guix-gentoo overlay</a>

### Aliases:

`nix-install <package>`: Add a package to the file (uses sed) and rebuild.  
`nix-uninstall <package>`: Remove a package from the file (uses sed) and rebuild.  
`nix-update`: Update and rebuild the system.  
`nix-rebuild`: Rebuild the system but don't update.  
`nix-develop`: Makes a nix dev shell for the system. Redundant when using flakey (default).
