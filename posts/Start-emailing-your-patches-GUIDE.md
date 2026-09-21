---
title: Start emailing your patches GUIDE
date: 2025-09-11
order: 031
tags: articles
source: "therealtruex.com (recovered from web.archive.org; date approximate, first archive crawl on therealtruex.com)"
---

## Why we need this

Corporate git websites are bloat. All that work can be done from the command line with `git` or `patch` and a mail client like neomutt. Hell, it can be done without a client using sendmail.

## Why this is better for projects

By using a email and static content workflow, and making a mailing list, it is just a better deal for the community working on the program. Discussions can happen on the same mailing list server that patches are sent. How are discussions done on git(hub/lab/tea)? At least sourceforge had its own mailing lists.

## How to send a patch

I have some helper scripts. There are two ways:

- `git format-patch` (my preferred method)

- the `patch` program is another option (more based)

- Have git just do it and send it with a script. Cross your fingers!

- Have git write to a file, then copy it on the command line and paste it in on the command line!

- Have a helper script that does more formatting and whitespace removal. Very good for picky mailing list hawks like when contributing to suckless.org or the linux kernel.

- Don't hand-edit patches. It breaks them.

## How to host your own git static site generator (not Gittea)

Let's not use Gittea. It is bloat. Self-hosting a whole corporate dynamic git website is madness. Let's use <a href="https://codemadness.org/git/stagit/file/README.html" class="external-link">codemadness.org stagit</a>. There is even a gopher edition. All it does is make a simple site like the one linked there.

### my format-patch scripts

These use some dmenu scripts like y-or-n and prompt-mail so they will not be plug-and-play on your system. A good starting place though.

`format-patch`  
Doesn't do anything special. Formats a patch and asks me if I want to send it up without looking for save it to a file.  
If I save it to a file I will use that as the mail message.

    #!/bin/sh
    # also see suckless-format-patch and friends
    [ -z "$1" ] && EMAIL="$(prompt-mail)"
    [ -n "$1" ] && EMAIL="$1"
    echo "Formatting a patch file..."
    TAG="$(git tag --points-at HEAD^)"
    TAGDIFF="${PWD##*/}-$(sed 's:.*/::' .git/HEAD)-$(date +%Y%m%d)-$TAG.diff"
    fp () {
        file="${PWD##*/}-$(sed 's:.*/::' .git/HEAD)-$1.diff"
        git format-patch \
            --to "$EMAIL" \
            --stdout HEAD^ \
            > "$file" \
        && echo "Wrote $file."
    }
    [ -n "$TAG" ] && fp "$TAG"
    [ -z "$TAG" ] && fp "$(date +%Y%m%d)-$(git rev-parse --short HEAD)"
    echo "Patch made."
    y-or-n "git-send-email it right now?" && git-send-email --to "$EMAIL" "$file" \
      && y-or-n "delete patch file?" && rm "$file"

`format-patch-suckless`  
I guess I got this from the suckless mailing list when I sent something without clang formatting.

    #!/bin/sh
    [ -z "$1" ] && EMAIL="$(prompt-mail)"
    [ -n "$1" ] && EMAIL="$1"
    y-or-n "Clang format (for C)?" && clang-format -i *.h *.c && echo 'formatted'
    y-or-n "Stage that?" && git stage *.h *.c
    y-or-n "Amend it to commit silently?" && git cane
    TAG="$(git tag --points-at HEAD^)"
    TAGDIFF="${PWD##*/}-$(sed 's:.*/::' .git/HEAD)-$(date +%Y%m%d)-$TAG.diff"
    fp () {
        file="${PWD##*/}-$(sed 's:.*/::' .git/HEAD)-$1.diff"
        git format-patch \
            --subject-prefix="${PWD##*/}][PATCH" \
            --to "$EMAIL" \
            --stdout HEAD^ \
        > "$file" \
        && echo "Wrote $file."
    }
    [ -n "$TAG" ] && fp "$TAG"
    [ -z "$TAG" ] && fp "$(date +%Y%m%d)-$(git rev-parse --short HEAD)"
    echo "Patch made."

`format-patch-send`  
Beam it up without looking!

    #!/bin/sh
    # also see suckless-format-patch and friends
    EMAIL=""
    [ -z "$1" ] && EMAIL="$(prompt-mail)"
    [ -n "$1" ] && EMAIL="$1"
    echo "Formatting a patch file..."
    TAG="$(git tag --points-at HEAD^)"
    TAGDIFF="${PWD##*/}-$(sed 's:.*/::' .git/HEAD)-$(date +%Y%m%d)-$TAG.diff"
    fp () {
        file="${PWD##*/}-$(sed 's:.*/::' .git/HEAD)-$1.diff"
        git format-patch \
            --to "$EMAIL" \
            --stdout HEAD^ \
            > "$file" \
        && echo "Wrote $file."
    }
    [ -n "$TAG" ] && fp "$TAG"
    [ -z "$TAG" ] && fp "$(date +%Y%m%d)-$(git rev-parse --short HEAD)"
    echo "Patch made."
    git-send-email --to "$EMAIL" "$file"
    rm "$file"
