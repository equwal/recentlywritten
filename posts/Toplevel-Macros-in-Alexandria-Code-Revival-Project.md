---
title: "Toplevel Macros in Alexandria: Code Revival Project"
date: 2019-04-29
order: 005
tags: articles
source: "therealtruex.com (recovered from web.archive.org; date exact, from the original .post file in equwal/truex.eu)"
---

\[image /static/lisp-alu-blue.gif, class logo, alt "ALU lisp logo, CC attributed" \]

After finding the old <a href="http://www.lispniks.com/cl-gardeners/" class="external-link">CL Gardeners' site (now defunct)</a>, I thought I'd carry the torch this afternoon. They have four suggestions for "gardening" Common Lisp:

- **Consumer reports**: Write a comparison of different libraries that do the same or similar things. Summarize the strengths and weaknesses of each. Include "under the hood" assessment; would I want to maintain this if the author got hit by a bus?

- **Code revival**: Adopt an abandoned library and clean up the bit rot.

- **Code mining**: Dig through existing open source code bases of large applications and extract bits that can be packaged as a useful stand-alone libraries.

- **Implementation convergence**: Help bring various CL implementations into alignment where there's no good reason for them to differ. This can be done by providing portability libraries or, better yet, by doing the decidedly non-trivial work of finding a common ground that different implementers can actually agree on and then doing whatever it takes to convince them to make the necessary changes to their implementation. ("Whatever it takes" in this case, likely includes patches, test suites, documentation, and civil, egoless, participation in relevant developer forums.)

Browsing the site on the web archive (since they have been defunct for ten years), I <a href="http://wiki.alu.org/Code_mining_and_revival" class="external-link">found a todo item</a>.

John Connors proposes to revivive the Port module from the <a href="http://clocc.sourceforge.net/dist/port.html" class="external-link">CLOCC</a>: it was a small module that provides portable acess to filesystems, threads and sockets.

I was excited to see that two of the functions in <a href="static/ext.lisp" class="external-link">ext.lisp</a> were already in my uitility library (as version spontaneously written by me), so clearly they deserved to be revived!

Cleaned them up and <a href="https://gitlab.common-lisp.net/alexandria/alexandria/merge_requests/12" class="external-link">submitted to alexandria</a>.
