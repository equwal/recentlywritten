---
title: Hackermen neomutt config
date: 2025-09-08
order: 049
tags: articles
source: "therealtruex.com (recovered from web.archive.org; date exact, from the therealtruex.com feed)"
---

<a href="https://github.com/HACKERMEN-ORG/neomutt-config" class="external-link">link to git hub website repository</a>

## Why this exists

So my friends can send me encrypted mail, bottom post, and send/apply  
patches.

&gt; "I installed neomutt but I haven't figured it out"

Neomutt defaults have legacy lock-in issues.

## Barebones neomutt config

### Features

- optional PGP encryption/signing/both

- Automatic PGP decryption

- hjkl movement: navigate

- n/p movement: read+page mails like in the GNUs client

- open HTML mail in the $BROWSER with a bind

- edit mail in the $EDITOR

- Can setup multiple mailboxes (just source more files like me-TEMPLATE)

- arcane-config-TEMPLATE to keep user-specific config separate

- Lots of colors

- A million default setting changes that don't do much

- signature

### Subjective User Experience

<a href="http://www.booksatoz.com/witsend/tea/orwell.htm" class="external-link">Like this</a>

### Install

- clone to ~/.config/neomutt

- install neomutt

- install isync and sendmail for IMAP and SMTP

- install notmuch to search your mails

- install gpg (optional)

- Edit and rename the TEMPLATE files

- Edit the sources at the top of the neomuttrc to fit your needs

### Good alias to do background mail stuff when opening the client

### TODO

- Automatic WKD for automatic PGP encryption without manually adding keys
