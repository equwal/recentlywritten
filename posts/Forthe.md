---
title: Forthe
date: 2019-06-22
order: 006
source: "therealtruex.com (recovered from web.archive.org; date exact, from the original .post file in equwal/truex.eu)"
---

&gt; Go Forth\[e\] and ye shall find.

Version 0.0.2

Docs available as: HTML([https://recentlywritten.com/eforth)](eforth).html) \| ORG-mode(<https://github.com/equwal/Forthe/blob/master/README.org)>

Code on github(<https://github.com/equwal/Forthe)>

This is an implementation of an elisp-based FORTH. This means that elisp is used to write the function definitions. In fact, no FORTH function definitions are provided at all to the Forthe macro from the library. The main benefit comes from reducing the verbosity of elisp in a small Domain Specific Language the programmer defines.

The following is a code sample for generating an eight-panel frame, which requires 37 function calls. See how compact FORTH code is? This is further helped by closing the \`nil\` required by \`switch-to-buffer\` into that lambda.

    (require 'forthe)
    (defun 8pan ()
      (interactive)
      (forthe '((n (lambda () (switch-to-buffer nil)))
                (d split-window-below)
                (r split-window-right)
                (wr windmove-right)
                (wl windmove-left)
                (wd windmove-down)
                (wu windmove-up)
                (del delete-other-windows)
                (cf make-frame-command)
                (nf other-frame))
              (cf del r r wr wr r d wr d wl wl d wl d wr n
                  wr n wr n wl wl wl wd n wr n wr n wr n wl
                  wl wl wu nf)))

That code, when executed, produces this:

<figure>
<a href="static/8pan.png" target="_blank"><img src="static/8pan.png" style="display:block;" alt="The eight panel setup made by the Forthe code snippet." /></a>
</figure>

The eight panel setup made by the Forthe code snippet.

## How does a FORTH work?

Forth code is gone through by the pointer moving through the call stack. When a non-function is found, it is added to the argument stack. When a function is found, it uses up the argument stack if it has arguments, or just runs if it is nullary. Nothing else is provided by the implementation.

## Why?

Needing to simulate a long set of keypresses to create an eight panel frame (37 in total), the best option seemeed to be obviously to create a miniature language with 1 or 2 letter abbreviations for each. Any arguments to a function can be passed before it in the list, and there is no need to define a return stack.

\`M-x 8pan\` will produce a new frame with eight panels (emacs-speak: "windows"). Here is how it works:

\`forthe\` is called with two arguments: the environment, name to function bindings, and the Forthe code list. You can see in this example that there are no arguments used, since all of these are nullary functions. Here is how you might use arguments:

    (forthe '((ff . find-file)) ("~/some-file" ff))

Here we ran \`(find-file "~/some-file")\`. If you need to reuse a FORTH,  
just define an \`\*env\*\` variable:

    (defvar *env* '((>r . to-return-stack) ...))
    (forthe *env* (forth code here ...))

## Benefits

### Concise code.

Elisp's convention of names like:

    org--insert-structure-template-unique-keys

(a real function in the org mode) result in extremely verbose code. Using short, lexical names can reduce code size by a large portion.

### Compiles to Elisp.

Elisp is hated by many, and for good reason. It is an outdated and clunky lisp. Many have tried and failed to port Emacs to Common Lisp, without success. Forth, unlike modern lisps, is designed to be easy to implement. Writing a FORTH compiler is in fact a trivial exercise.

## Drawbacks

### Elispers might not want to read your Forthe code

### Forthe is not stable

FORTH is <a href="http://www.flownet.com/gat/jpl-lisp.html" class="external-link">space-age technology</a>, and Forthe is young.

### Isn't verbose code self-documenting?

NO. Documented code is documented. Using short names with the bindings handily available in the same form is equally explicit as writing them out every time.

## Contributing

Please make sure your contributions are  
melpa(<https://github.com/melpa/melpa/blob/master/CONTRIBUTING.org)-friendly,>  
and documented in the README. You can use  
pandoc-mode(<http://joostkremers.github.io/pandoc-mode/)> to transform  
the README.org into plaintext `;;; Commentary:`.
