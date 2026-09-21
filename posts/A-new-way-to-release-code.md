---
title: A new way to release code
date: 2025-09-11
order: 035
tags: articles
source: "therealtruex.com (recovered from web.archive.org; date approximate, first archive crawl on therealtruex.com)"
---

Recently clobbered this together to send over IRC really quick. The power of pipes! Best for April fools.

git-update-simple.sh

<figure>
<a href="static/git-update-simple.sh.webp" target="_blank"><img src="static/git-update-simple.sh.webp" style="display:block;" alt="/static/git-update-simple.sh.webp" /></a>
</figure>

git-update-worker.sh

<figure>
<a href="static/git-update-worker.sh.webp" target="_blank"><img src="static/git-update-worker.sh.webp" style="display:block;" alt="/static/git-update-worker.sh.webp" /></a>
</figure>

    # encode
    base64 -w0 file.sh | qrencode -t ansiutf8 # better aliased to qr...

    # decode 
    zbarimg file.webp | base64 -d >file.sh
