---
title: Install Nix on Gentoo
date: 2025-11-13
order: 040
tags: articles
source: "therealtruex.com (recovered from web.archive.org; date approximate, first archive crawl on therealtruex.com)"
---

Code:  
<a href="https://github.com/trofi/nix-guix-gentoo" class="external-link">nix and guix overlay for gentoo</a>

Guide:  
<a href="https://trofi.github.io/posts/196-nix-on-gentoo-howto.html" class="external-link">nix on gentoo howto by trofi</a>

It is *very* easy to do this way.

## BENEFITS

- Install bloatware with declarable builds and zero effort

- Install basedware with gentoo

- Be able to build stuff from source etc. without the nix abstractions  
  getting in the way

- Learn how nix works for deployable declarative server infrastructure (<a href="Why-YOU-need-NixOS-on-your-desktop-and-servers.html" class="external-link">see post on that here</a>)

## Alternatives

- Use kexec and install both, switch at runtime (complex, can't have both on at the same time)

- Use a gentoo chroot from within nix (complex, requires knowledge of the nix language). Use qemu or whatever. Very smart option. All the same benefits as above but the base system is also more reliable because it is nix. Install portage in NixOS (madness)

## Gotchas

Nix *should* work with the official `curl` into `sudo sh` method, and not break the system. Instead, use the overlay. Let the overlay maintainer deal with those problems.
