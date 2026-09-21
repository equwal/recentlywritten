---
title: Generate your OWN internet web site easily
date: 2025-11-13
order: 030
tags: articles
source: "therealtruex.com (recovered from web.archive.org; date approximate, first archive crawl on therealtruex.com)"
---

## S(S/C)Gs (Static (Site/Content) Generators) don't have to be bloat

Wordpress sucks. Hugo is better. But we can do even better. This site is/was generated at one time with  
<a href="https://github.com/coleslaw-org/coleslaw" class="external-link">coleslaw</a>, a Common lisp static site generator. Even these solutions are  
bloat. Introducing: cl-yag, another Common lisp SSG in only 4/20 lines of code (exactly! Well, in one of the commits it hit the magic number).

## How simple can it get?

The simplest is to just hand-write HTML. Check out <a href="https://bitreich.org" class="external-link">bitreich.org</a>. Can you believe that site  
is handwritten HTML? Here's the git repo link: `git clone git://bitreich.org/bitreich-www/`. Crazy.  
<a href="https://github.com/equwal/cl-yag" class="external-link">cl-yag</a> is the next level up. It converts markdown to html. That's it. Coleslaw is a bit more complex, with various plugins and so on. Still very nice and hackable. Hugo is the famously "simple" one. I think it is too complex. Wordpress-likes are the worst, not even generating a site but making everything dynamically on the fly.

## Picking an SSG

It has been a theme that people tend to use SSGs in their favourite higher level languages, (if savvy), or just opting for a dynamically loaded site (if not savvy). So pick an SSG in your favourite higher level language, but make sure it is simple enough.

### Blogware

#### Bash shell + pandoc

Hasn't everyone done it this way before? Some shell scripting is needed to upload the scripts, so it is natural to just build it out until it is hard to tell what it is even doing anymore.

#### <a href="https://gohugo.io/" class="external-link">Hugo</a>

Hugo is about medium bloat level and written in go. Go is not cool enough.

#### <a href="https://jekyllrb.com/" class="external-link">Jekyll</a>

Haskellers like to use <a href="https://jekyllrb.com/" class="external-link">Jekyll</a>, another medium level of bloat.

#### Common lisp ones

- <a href="https://github.com/equwal/cl-yag" class="external-link">cl-yag</a>: ultra-simple. Lisp wins.

- <a href="https://github.com/coleslaw-org/coleslaw" class="external-link">coleslaw</a>: a bit more bloat

### Not for blogs

- <a href="https://codemadness.org/git/stagit/file/README.html" class="external-link">stagit</a> generates a site for a git repo.

- Documentation generators (eg. using GNU make or Quickref for lisp) can sometimes make HTML sites. These tend to be really nice.
