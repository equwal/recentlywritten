---
title: sbm
date: 2023-09-22
order: 009
tags: projects
source: "therealtruex.com (recovered from web.archive.org; date approximate, first archive crawl on equwal.com)"
---

How to keep track of thousands of bookmarks with a simple tool and your own tags.

  

### Copy or add a bookmark:

![Copy or Add a bookmark.](static/firstscrn.png)

### dmenu fuzzy search dialogue for all your bookmarks. Simple text format.

![Searching for booksmarks](static/secondscrn.png)  
[github.com/equwal/sbm](https://github.com/equwal/sbm)

It needs to be made with make and the config file should be edited too. More docs soon, although it is a simple program. Details in the source.

[Karl Voit's article (the good ideas are from here)](https://karl-voit.at/2022/01/29/How-to-Use-Tags/)  
[download the dmenu package](static/sbm-0.1.tar.gz)  

Grand Generalized Dmenu Tagger coming soon to stores near you.

### Scripts

- 1.  bm -- search & copy, add, or edit/plumb a bookmark (edit/plumb not implemented yet).

### Install

Requirements: xclip, dmenu, make

- Export the $BOOKMARKS variable containing the location of the bookmarks. I like to use ~/.local/share/sbm to store my config.  
  mkdir -p ~/.local/share/sbm/  
  ~/.bashrc  
  export BOOKMARKS="~/.local/share/sbm/bookmarks"
- define your tag choices in USERTAGS and export that too. The file syntax is:  
  &lt;tag&gt;&lt;space&gt;&lt;description&gt;&lt;newline&gt;  
  See the included usertags file for an example.  
  ~/.bashrc  
  export USERTAGS="~/.local/share/sbm/usertags"  
- make install

### Usage

Execute

search-term is an extended regex for grep.

  
  

Bookmark lines are space separated values  
URL description

  

### Tip

The bookmarks can be converted from the web browser formats using netscape-bookmark-converter at

[github.com/jhh/netscape-bookmark-converter](https://github.com/jhh/netscape-bookmark-converter)
