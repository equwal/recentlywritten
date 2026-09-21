---
title: "Don't Self-Host: Self-Locally"
date: 2024-05-21
order: 014
tags: articles
source: "therealtruex.com (recovered from web.archive.org; date approximate, first archive crawl on equwal.com)"
---

My most recent posts are about how to setup a few services locally, which I think *would* be interesting to the consumers of the idea of "self-hosting".

It seems to me that the desire to self-host services actually comes, for many, from the problem of having an insufficiently power computing environment locally to do the necessary tasks.

As far as I can tell: there is no reason to self-host unless one of these factors is met:

- the content must be stored away from the user's main machine (backup server)

- the content must be deliverable to multiple devices for some reason (password manager)

- the hosting isn't really "self" hosting, but actually "family" or  
  "team" hosting (gittea).

So why do people self-host things that they can just run locally?

Disease \#1: Windows or MacOS installed on main machine

This **disease** leads the self-hosting for no better reason than because the host machines are capable of running a decent free operating system. If you're interested in a self-hosthing project because "I can't do that on my local machine", this is probably why. Install Gentoo instead.

Disease \#2: Muh devices

If you're installing some self-hosting service to deliver content to restrictive devices like cell-phones, consider installing an HFS or FTP server instead to transfer the files around. Consider replacing your worst devices with more free alternatives.

Disease \#3: Muh FREEDOM

If you are trying to avoid using services run by others in order to take back some freedom (eg. git server or mail server), then that is a great thing to be doing. But often we find that self-hosting such things is so cumbersome that it leads to massive time-wasting in order to setup something which is already available from a third-party service. How do we deal with this problem? I have made this rule of thumb: if it can't be installed *and working* after partitioning a low-resource fresh linux VPS and running a single command, it is not ready for this use case.

Any project which aims to fill this niche must be as easy to install as making an account on the competing service, and must run on a low-resource VPS to avoid costing the user a significant amount of money.

A final solution, worth considering, which I have tested out without success, is to have a group of friends who each dedicate themselves to hosting different parts of the services you use. This rarely works so well. Some people will neglect their service, or underconfigure it, or let it go down without bringing it back up, let it get out of date, or just stop paying for the hosting or domain name.

The lesson for project maintainers in this area should be clear: you project should install instantly with no configuration and cannot be a resource hog by default either.

I'd like to see a list of such projects keeping track of this:

1.  installs in one command

2.  uses &lt;=8GB of disk space

3.  uses &lt;=1GB of ram max

4.  requires no configuration to start using

And I'd like to see more such projects in the future. If only hosting a mail server were that easy...
