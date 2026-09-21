---
title: Why YOU need NixOS on your desktop and servers
date: 2025-09-11
order: 033
tags: articles
source: "therealtruex.com (recovered from web.archive.org; date approximate, first archive crawl on therealtruex.com)"
---

## Why?

To stop working on your computer and start working with your computer.

<https://nixos.org/>

NixOS has:

- Declarative linux builds

- Declarative package installs (package manager can be installed on any

linux too <https://github.com/NixOS/nix)>

- Tools to deploy to servers (deploy-rs <https://serokell.io/blog/deploy-rs)>

- Terraform to have a declarative and **shareable** deployment of cloud infrastructrure (terraform <https://nix.dev/tutorials/nixos/deploying-nixos-using-terraform.html)>

- Automated declarative partitioning with disco (<https://nixos.wiki/wiki/Disko)>

- A tool to bootstrap NixOS from any linux and replace it on a server (nixos-anywhere <https://github.com/nix-community/nixos-anywhere)>

- The above tools can also be used on a USB stick to make an instantly

deployable desktop. If you setup a good NixOS

desktop, it is not a special root directory that a lot of work went into. It is a few files which DEPLOY into the  
full-fledged thing. A good NixOS config can even delete the entire  
root directory on shutdown.

## Is Gentoo not good enough

Gentoo makes for a great minimal system. Not a great bloatware system.

## Downsides

Installing your own software or packages not in the repo is extra annoying. All software developers on nix are also  
package maintainers.
