---
title: install gentoo hints
date: 2024-05-21
order: 019
tags: articles
source: "therealtruex.com (recovered from web.archive.org; date approximate, first archive crawl on equwal.com)"
---

To install gentoo, we must RTFM. Here is the FM: <https://wiki.gentoo.org/wiki/Handbook:AMD64>

Set your `-march=` to the oldest architechture you might use, and your `-mtune=` to the fanciest architecture you might use. Also, `-O2 and -pipe`. Eg., if we want a Gentoo to work on a based Thinkpad and on a fancy gamerboy machine:

    CFLAGS="-march=core2 -mtune=skylake -O2 -pipe
    CPU_FLAGS_X86="mmx mmxext sse sse2 sse3 sse4_1 ssse3" # common flags between core2 and skylake

`MAKEOPTS="-j2 -l2"` is great for a two-core machine. But what if we want 8 threads when on an 8-core machine? I present the script, `wrap-8-core-kludge`.

    #!/bin/sh
    sed -i 's/-j[0-9] -l[0-9]/-j8 -l8/' /etc/portage/make.conf
    sed -i 's/--jobs=[0-9] --load-average=[0-9]/--jobs=8 --load-average=8/' /etc/portage/make.conf

    "$@"

    sed -i 's/-j8 -l8/-j2 -l2/' /etc/portage/make.conf
    sed -i 's/--jobs=[0-9] --load-average=[0-9]/--jobs=2 --load-average=2/' /etc/portage/make.conf

Call it with `wrap-8-core-kludge emerge --ask vim` when on the fancy machine.

Don't fuck up your profile: choose the latest hardened multilib. Don't switch!

Use the Zen kernel (it's the kernel config intended for desktops).

Don't build modules into your kernel, which is annoying, edit `/etc/conf.d` with the modules you need instead. This is no different than how it is normally done in other distros.

Fancy GPU owners: watch out for NVIDIA. Consider using nouveau if the nvidia drives don't work and crash your Xorg for example.
