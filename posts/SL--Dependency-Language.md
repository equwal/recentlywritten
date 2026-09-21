---
title: SL--Dependency Language
date: 2019-07-02
order: 007
tags: articles
source: "therealtruex.com (recovered from web.archive.org; date exact, from the original .post file in equwal/truex.eu)"
---

\[image /static/bnf-sl.png\]

\[Github\](<https://github.com/equwal/sl)> \|  
\[Homepage\]([https://recentlywritten.com/sl--dependency-language)](sl--dependency-language).html)

SL is a Domain Specific Language (ie. simple programming language) for dealing with ambiguous dependency situations. It is more expressive than the standard method in Common Lisp: the use of reader macros.

## Examples

The only external macro is `defsl`. There is BNF below, though some examples are in order.

### Simple

This defines `operator-arglist` (in the current package) as  
`slynk:operator-arglist`, unless `slynk` is unavailable, in which case  
`swank:operator-arglist` is used.

`` ` commonlisp ``  
`(defsl operator-arglist :fn :eq slynk swank)`  
\`

The `:fn` shows that a FUNCTION is being defined, as opposed to a `:sym` symbol; to define another namespace a list is required: `(accessor-fn boundp-fn)`. `:eq` denotes that the packages have the same name for this function. For this simplest of examples, consider the reader macro alternative.

\`\`\` commonlisp  
\#+(and slynk swank) (setf (symbol-function 'operator-arglist)

slynk:operator-arglist)

\#+(and slynk (not swank)) (setf (symbol-function 'operator-arglist)

slynk:operator-arglist)

\#+(and (not slynk) swank) (setf (symbol-function 'operator-arglist)

swank:operator-arglist)

\#  
\`\`\`

### Complex example

\`\`\` commonlisp  
(defsl gensymmer (macro-function macro-function)

utils with-unique-names  
alexandria with-gensyms)

\`\`\`

To break this down:

\- `gensymmer` is now defined as `utils:with-unique-names`.  
- `macro-function` was used twice: once as a `setfable` place, and

once as a

predicate.

\- If the `utils` package does not exist, then `sl::gensymmer` is

defined as `alexanrdia:with-gensyms`.

\- to define different names for each package, they must be "qualified"

with the package name.

## Install

Use ASDF to install. Usually this should work:

`` ` bash ``  
`$ cd ~/common-lisp/`  
`$ git clone git@github.com:equwal/sl.git`  
`CL-USER> (asdf:load-system :sl)`  
\`

## BNF

\`\`\` example  
(defsl &lt;sl-name&gt; &lt;fnsym&gt; &lt;package spec&gt;)  
&lt;package spec&gt; ::= &lt;eq&gt;

\| &lt;packages&gt;

&lt;eq&gt; ::= :eq &lt;preferences&gt;  
&lt;fnsym&gt; ::= :fn

\| :sym  
\| &lt;fnpair&gt;

&lt;fnpair&gt; ::= (setfable-place bound-predicate)  
&lt;packages&gt; ::= &lt;package&gt; &lt;fn name&gt; &lt;more packages&gt;  
&lt;package&gt; ::= symbol  
&lt;fn name&gt; ::= symbol  
&lt;preferences&gt; ::= &lt;package&gt; &lt;more packages&gt;  
&lt;more packages&gt; ::= ε

\| &lt;package&gt;  
\| &lt;package&gt; &lt;more packages&gt;

\`\`\`

## Issues:

\- This is a new thing.  
- `defsl` is quite possibly the world's most unhygenic macro: don't

expect anything about evaluation order or number of evaluations to  
be true.
