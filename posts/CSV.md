---
title: CSV
date: 2019-01-07
order: 004
tags: projects
source: "therealtruex.com (recovered from web.archive.org; date exact, from the original .post file in equwal/truex.eu)"
---

Defines a Common Lisp reader to read in CSV:

    #!"field","field","field"!#

or read in CSV **files** with \`slurp\`:

    (slurp #p"example.csv")

and the result is a list:

    CL-USER> (csv:slurp #p"example.csv")
    ("example" "csv" "list" "result" "like" " this")

This program invokes the lisp reader, making it:

1.  Potentially risky to use if the input it not trusted.

2.  Resource intensive for medium to large files (&gt;5MB).

Use with care.

## Features:

Supports escaping quotes the CSV way, with quad quotes:

    csv   -> lisp

    """" -> "\""

## Installation

\[Clone from github\](<https://github.com/equwal/CSV)> into asdf load directory and use ASDF to install:

    $ git clone https://github.com/equwal/CSV
    $ sbcl

    CL-USER> (asdf:load-system :csv)
    T

## Notes:

By default backslashes are not escaped. CSV is expected to be a list of quoted items, like this:

    #!"example","here"!#
