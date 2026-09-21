---
title: Get on mailing lists like it is 1989
date: 2025-09-11
order: 028
tags: articles
source: "therealtruex.com (recovered from web.archive.org; date approximate, first archive crawl on therealtruex.com)"
---

## More folders than INBOX and SPAM

Did you know that you can have as many folders as you want? The main issue is that the filtering needs to be reasonably smart and is usually done server-side. I use runbox.com as my host and it can do it. For mailing lists, the best way is to sign up, read the headers, and see what address that lists sends from. Then filter any mail with that email in the header into the correct folder. This is called a "milter" (mail filter). A spam filter is a kind of milter. Self-hosting email is one way to get very good milting.

## Good clients

Now you need a client that can read those folders. GNUs is an ancient one with a million features (think: emacs) which  
was (is?) useful for usenet and mailing list extraordinaires. Neomutt is the industry standard. If you mail host  
sucks you might really want a good client, so it can milt for you.

## Why YOU need mailing lists

- You can find out when software gets released and go upgrade

- You can develop stuff on the list

- You can find out when your favourite organization is calling for papers or organizing events

- Mailing lists are like IRC but for real work (you're on IRC right? want to <a href="iilchat-setup-too-easy.html" class="external-link">set that up?</a>)

- It is a better way to send small patches than the github.com website (using git/diff/patch etc.)

- You will learn things

- You will learn to notice instead of ignore bugs (a good thing? Maybe)

## How get mail

- Find your project

- Subscribe!

- Verify!

- Go to the archive or wait for a mail to get sent...

- Use the list's email address in the header or other identifying information (some lists prepend something like  
  \[list-name\] to every subject line just for easy miltering)

## How to send mail (important)

It's good to read a few mails from the list to make sure you are sending to the right one. Some lists also have lots of rules. Not following those rules will probably result in messages from annoyed maintainers. Use your mail clients REPLY-ALL function and quote the necessary material. Write your message BELOW the quoted material (not above) or even split up the previous message into pieces and reply under each.

## How to send a patch

git has a patch sending feature. It bypasses the mail client entirely. Another option is to write the patch to a file and send that with the client. DO NOT edit the diff within the patch manually! You can mess it up. I like to use `cat patch.diff | xclip -i` and paste the result directly into neomutt.
