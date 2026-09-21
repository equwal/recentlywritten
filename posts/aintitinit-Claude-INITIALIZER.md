---
title: "aintitinit: Claude INITIALIZER"
date: 2025-09-09
order: 051
tags: articles
source: "therealtruex.com (recovered from web.archive.org; date exact, from the therealtruex.com feed)"
---

<a href="https://github.com/HACKERMEN-ORG/aintitinit" class="external-link">git hub website repo</a>

## Claude INITIALIZER

INITIALIZES a new git project in the directory if there isn't one already,  
various files for claude in the directory to vibe code really hard,  
and hacks git to not track them.

The bare minimum necessary of course.

<figure>
<a href="static/gifmatrix.gif" target="_blank"><img src="static/gifmatrix.gif" style="display:block;" alt="gifmatrix.gif" /></a>
</figure>

### CLAUDE.md

This is the prompt claude will use in your project. Includes the prompt  
I had it generate, you can put anything you want in here.

### exclude (file)

Appends this to the git exclusions to not track the claude-related files

### .mcp.json

The mcp servers you will use. There are a million already added by the AI if you want to just use it as-ise.

### llms-full.txt

Tells claude about the LLMs in .mcp.jsonquire"cmp.utils.feedkeys".run(2)  
,

### regular use

    claude-init ; claude init

### alias

`ci` claude init.  
`cli` try to continue, otherwise initialize in the directory.

    alias ci='claude-init ; claude init'
    alias cli='f(){ mkdir -p "$1" ; cd ${1:-.} ; claude-init ; claude init }; claude --continue || f'

### TODO

- Automatically remove whitespace and base64 encode the CLAUDE.md from

the readable version to save tokens.
