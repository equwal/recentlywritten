---
title: Git extra worktrees ROCK
date: 2025-09-11
order: 032
tags: articles
source: "therealtruex.com (recovered from web.archive.org; date approximate, first archive crawl on therealtruex.com)"
---

A friend informed me that now git has explicit support for extra worktrees (rather than just blindly copying stuff). I was pretty incredulous about it at first – why would you need that? I figured it out pretty quick.

For when the local branch of the project is something that is actually used, so it is better to keep the files there as they are (for use). This makes it way easier to commit changes upstream without messing up the current system. For example, I have made a ton of changes to the coleslaw static site generator that makes this website, but haven't pushed many up because I actually use that worktree. Without the git worktree support the options are:

- Stash changes, work on the upstream, and reapply the changes.

- Copy the folder somewhere else.

The first option is error prone (easy to forget about those stashed changes etc.) and second option is better but leads to an issue – what if I want to use git to manage changes across these worktrees?

Git managing worktrees is worth the added complexity. Good change!
