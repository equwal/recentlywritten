---
title: dmenu scripts repo released
date: 2025-09-01
order: 037
tags: articles
source: "therealtruex.com (recovered from web.archive.org; date approximate, from the announced repo's creation)"
---

## My dmenu scripts <a href="https://github.com/equwal/dmenu-scripts" class="external-link">(click here for the git repo)</a>

<a href="https://tools.suckless.org/dmenu/" class="external-link">dmenu</a> is a simple program with a pretty selection window for X. These are my scripts for doing things with that.

## Scripts

- `dbc` <a href="static/dbc" class="external-link">(file)</a>: do `bc -l` on the input so you can do math

- `dbrowsel` <a href="static/dbrowsel" class="external-link">(file)</a>: select a browser for the link in the clipboard to open

- `dcpupower` <a href="static/dcpupower" class="external-link">(file)</a>: select CPU throttling settings

- `dkeymap` <a href="static/dkeymap" class="external-link">(file)</a>: select a keymap

- `dintelbacklight`<a href="static/dintelbacklight" class="external-link">(file)</a>: select backlight with dmenu or wjt (a scroll wheel)

- `ddmenu` <a href="static/ddmenu" class="external-link">(file)</a>: use dmenu to select a dmenu script – dception

- `dmenu_run/dsudo` <a href="static/dmenu_run" class="external-link">(dmenurun)</a>  
  <a href="static/dsudo" class="external-link">(dsudo)</a>: run programs, maybe from a keybind set with xbindkeys or a wm

- `dshow` <a href="static/dshow" class="external-link">(file)</a>: run a command and show output without having to open a terminal

- `dsysact` <a href="static/dsysact" class="external-link">(file)</a>: reboot/shutdown/etc.

- `passmenu` <a href="static/passmenu" class="external-link">(file)</a>: unix <a href="https://www.passwordstore.org/" class="external-link">pass</a> interface

- `passmenu-otp` <a href="static/passmenu-otp" class="external-link">(file)</a>: like passmenu but copies the OTP code to clipboard.

- `y-or-n` <a href="static/y-or-n" class="external-link">(file)</a>: ask if yes or no

<!-- -->

- `dmenupass`<a href="static/dmenupass" class="external-link">(file)</a>: ask for the sudo password as SUDO\_ASKPASS

- `dmount/dunmount`<a href="static/dmount" class="external-link">(dmount)</a>  
  <a href="static/dunmount" class="external-link">(dunmount)</a>: mount/unmount drives

- `dopenrc-runlevel`<a href="static/dopenrc-runlevel" class="external-link">(file)</a>: choose openrc runlevels for various advanced tasks like networking

- `drssadd`<a href="static/drssadd" class="external-link">(file)</a>: add RSS links to the newsboat config from clipboad

- `dtimer`<a href="static/dtimer" class="external-link">(file)</a>: a coffee timer. Check out <a href="https://www.uninformativ.de/git/countty/file/README.html" class="external-link">uniformativ.de</a>'s countty for a non-dmenu one.

- `dtmpl`<a href="static/dtmpl" class="external-link">(file)</a>: copy template files from a directory into the clipboard

### Submodules with even more scripts

- <a href="https://github.com/equwal/dmodurl" class="external-link">dmodurl</a>: change a URL in the clipboard according to sed rules

- <a href="https://github.com/equwal/dpatchmail" class="external-link">dpatchmail</a>: do (code) patches with dmenu

- <a href="https://github.com/equwal/dsad" class="external-link">dsad</a>: manage background audio with Simple Audio Daemon

- <a href="https://github.com/equwal/prompt-mail" class="external-link">prompt-mail</a>: manage your email contacts list (neomutt compatible)

- <a href="https://github.com/equwal/sbm" class="external-link">sbm</a>: manage your bookmarks list

- <a href="https://github.com/equwal/tnot" class="external-link">tnot</a>: like tmux but not (using <a href="https://github.com/deadpixi/mtm" class="external-link">mtm</a> and <a href="https://www.brain-dump.org/projects/abduco/" class="external-link">abduco</a>)

### Quick hack for dmenu with sudo

You *could* do:  
`` ` ``  
`alias program='SUDO_ASKPASS=dgivepass sudo -A program'`  
\`

but it won't come up in dmenu\_run. So move \`program\` to \`program-aux\` and make a new script at program:  
`` ` ``  
`#!/bin/sh`  
`SUDO_ASKPASS=dgivepass sudo -A program`  
\`

### Utilities

- `intel-backlight`: uninformativ.de's intel-backlight script
