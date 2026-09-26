---
title: "\"Mailing list\" in one line"
date: 2025-09-06
order: 045
tags: articles
source: "therealtruex.com (recovered from web.archive.org; date exact, from the therealtruex.com feed)"
links:
- '[Hackermen neomutt config](Hackermen-neomutt-config.html)'
- '[Add my RSS feed to your feed reader](Add-my-RSS-feed-to-your-feed-reader.html)'
- '[Get on mailing lists like it is 1989](Get-on-mailing-lists-like-it-is-1989.html)'
- '[Start emailing your patches GUIDE](Start-emailing-your-patches-GUIDE.html)'
- '[github.com/HACKERMEN-ORG/neomutt-config](https://github.com/HACKERMEN-ORG/neomutt-config)'
---

# ==> list-of-emails.txt <==
    Your Name <youremail@example.com>
    ...
    ...
    # ==> mailing-list.sh <==
    #!/bin/sh
    while read email; do neomutt -H "$MESSAGEFILE" "$email" < /dev/null; done < /dev/stdin

Make an mbox email with neomutt of whatever.  
Now run it.

    cat list-of-emails.txt | mailing-list.sh email-file.mbox

Now make an alias in (zsh) for extreme FAST mailing

    alias mail2list='f(){ [ -n "$1" ] && cat /home/jose/.config/mailing-list/the-mailing-list | mailing-list "$1" }; f'
