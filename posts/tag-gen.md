---
title: tag-gen
date: 2025-09-05
order: 042
tags: articles
source: "therealtruex.com (recovered from web.archive.org; date exact, from the therealtruex.com feed)"
---

<a href="https://github.com/equwal/tag-gen" class="external-link">git hub website link here</a>

## PURPOSE:

Specify which directories to generate tags for, and what kinds of files  
to generate them for in each. Good for tag jumping in an editor.

## Requirements

- An ANSI Common Lisp

- ASDF (for UIOP portable shell facilities)

- A ctags-like program.

- Desire

## Non-features/Non-Requirements

- Not multithreaded!

- Works with any tag program (doesn't do any tagging itself)!

- No special configuration language, just two sexps!

- Does not use Roswell!

- Works on my machine!

## Install

1.  put this folder somewhere

2.  put the config folder somewhere

3.  put the tag-gen.sh and tag-gen-root.sh somewhere

4.  put the cron job in

5.  put your configuration in

Edit the config/ files with the paths you want tagged.

Add a command to cron like:

    */29 * * * * su username -c 'tag-gen.sh      /path/to/tag-gen/local.conf /usr/bin/etags'
    */29 * * * * su root     -c 'tag-gen-root.sh /path/to/tag-gen/root.conf  /usr/bin/etags'

for routine tagging.

tags.conf needs to contain this kind of format:

    "/path/to/ctags/program"

    (("<directory>" ("<ext>"...))
      ...)

The tags program reads the file extensions to infer what language to do tags for. Specify the directories you want to use tag browsing programs on.

See examples in config-examples. Configuration is **mandatory** for the program to work. The format is just read with the lisp reader and is very simple.
