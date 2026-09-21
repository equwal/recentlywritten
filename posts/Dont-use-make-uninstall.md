---
title: "Don't use make (un)install"
date: 2025-09-06
order: 044
tags: articles
source: "therealtruex.com (recovered from web.archive.org; date exact, from the therealtruex.com feed)"
---

Much like "sending a file over the internet", "installing software" has not been solved generally. There's more than one way to do it. A mixed system that uses both declarative and traditional package managers in unison is the way to go for now. It would be ideal to package a non-specific declarative package install definition together with the source of programs, and make unnecessary the concept of a "package maintainer" for each package manager. The existence of two declarative package managers already suggests the need for an arbitrary standard which can transformed into a definition for any declarative system.

I used an LLM to generate hundreds of uninstall targets for make and submit patches to the authors. Three authors wouldn't patch my changes in. Are they right? The argument goes like this: the software should not be installed using make because make doesn't register file install locations in any way. Therefore make only uninstalles the correct files by coincidence. For installing, it can do a little bit better using the `install` program.

There is a better way which has been around forever: just use the build from the source tree.

Let's consider how the ASDF manager for Common Lisp works (and by extension, quicklisp, which gets ASDF systems from the internet). It has a central registry variable with a list of paths, and it recurs down it to findsystem definitions when the definition loader is called. Then it loads the files into the lisp core. As a result there is no need to have code and binaries in different places. If compiled, it will just store the compiled files together with the source. To "uninstall", just delete that project's directory. Let's call this "in-situ" package management.

With POSIX make and arbitrary build systems like it, the subordinate model can't work, because ASDF can't possibly predetermine what binaries to use and where to put them. The logic can't be generalized in the package manager, and has to be managed per-project. Every package manager now needs thousands of packages with hundreds of maintainers. Typical package managers do installs imperatively and tracking the file location metadata in the system flies. When the binaries or libraries are edited or moved by an outside program or a user without editing the registry, (for example, by trying to re-install gcc manually from source) it can break the system.

The simplest way to get this kind of management is to use the symlinks to the files from directories in the PATH. If the installed files change, but the path stays the same, they can be updated without changing the "installed" data at all. IT's still tenuous – an update to the build system could change the locations of built binaries.

## Declarative and in-situ package managers: Nix and Guix

The Nix package manager uses the model of keeping the binaries in the package manager's repo without "installing" them. This is why nix systems can't break unless impure hacks are used.. Since binaries and libraries are not *moved*, there is no chance of overwriting things or losing track of where they are!

Let's not distribute installers in the code builders like make, but instead build it and use it as-is from the folder using an in-situ package manager. By chance, that's exactly how declarative package managers work, and not how traditional registry managers work.

The benefits of declarative package management are plentiful, but come at the cost of making development of one's own programs and installation of as-yet unpackaged software from the source more difficult. Even simple shell scripts added to the PATH can't be done simply, the whole abstraction of the package manager has to be invoked, or some impure hacks are needed.

The declarative package managers can be installed on another linux which has a traditional registry package manager and impure environment. Why not have both? It is "possible" to install two traditional package managers on the same system, and will almost immediately break it due to out of sync dependency management. With in-situ and system registry managers installed together there are no conflicts, because the in-situ manager stays in its own folder.

Package management solutions:

- In the build system (`make install/uninstall`)

- Using symlinks in the path (`ln -s`)

- In-situ

- Declarative + in-situ

- Traditional system registry package managers

- Declarative meta-package management language (does not exist yet).

## Final note: should there be an uninstall target in the Makefile?

- No, if the makefile doesn't "install" the software but just builds it.

- Yes, if the makefile does "install" the software, it should also have a way to "uninstall" the software.

## Final note: but I want to uninstall something anyway with nothing more than make and not use an uninstall!

### Use make install to a DESTDIR to see what files it installs and remove them

It is the best way, albeit tedious and manual.

"rofl0r" says:

As the maintainer of a <a href="https://github.com/sabotage-linux/sabotage" class="external-link">linux distro</a>, i quickly learned that the DESTDIR &lt;<https://www.gnu.org/prep/standards/html_node/DESTDIR.html>&gt; variable is really \*the\* key to cleanly track the files installed by arbitrary source projects, apart from much more elaborate solutions such as doing installs in a sandbox with before/after comparison or using file-system monitoring (e.g. via ptrace). it allows to install all files into an arbitrary location, which makes it easy to e.g. create a filelist, or in the case of my distro, the DESTDIR is used as the final location where the package actually lives. if we take curl as an example, it would be installed with empty prefix and DESTDIR=/opt/curl, so the curl binary would live in /opt/curl/bin/curl, headers in /opt/curl/include, etc. the package manager, as the last step of the build/install process, then iterates over all files in /opt/curl and symlinks them into the root prefix, e.g. /bin/curl will point to /opt/curl/bin/curl. uninstalling a package then is as easy as to iterate again over the files in /opt/curl (or in my case a filelist made thereof, to reduce the amount of syscalls) to remove the symlinks and finally rm -rf /opt/curl to delete the actual files.
