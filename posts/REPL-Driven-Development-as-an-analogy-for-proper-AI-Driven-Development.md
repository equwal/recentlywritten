---
title: REPL Driven Development as an analogy for proper AI Driven Development
date: 2025-09-08
order: 048
tags: articles
source: "therealtruex.com (recovered from web.archive.org; date exact, from the therealtruex.com feed)"
---

I recently "discovered" using an LLM to code with, and was looking forward to doing some "REPL Driven Development" on it. Unfortunately, I discovered that all of the vibe-coded vim plugins just use vim as a (terrible) multiplexer for the Claude Code CLI program. The Emacs plugin <a href="https://github.com/manzaltu/claude-code-ide.el" class="external-link">claude-code-ide.el</a> does it right by having text selection and other send-to-repl-like actions like a proper Lisp IDE does.

This was merely a lisp thing to do because the parentheses clearly define selection to send. That's no longer an issue, the LLM can filter out the garbage in a way that a GOFAI symbolic rule-based pattern matcher can't.

The command line program likes to use diffs to edit files, and that is reasonable, but vastly inferior to REPL Driven Development from a scratch file. Below are some of my notes from 2019 for more context on using the lisp REPL. With RDD, you will never type something directly into a REPL again. "Never type into the REPL".

It's nice using grep to lookup some code I might have written 5 years ago that would have been ephemeral and lost at the REPL.

### REPL driven development

Stuard Halloway's 2017 talk at Chicago Clojure:

<a href="https://vimeo.com/223309989" class="external-link">vimeo</a>  
<a href="https://github.com/matthiasn/talk-transcripts/blob/master/Halloway_Stuart/REPLDrivenDevelopment.md" class="external-link">slide notes</a>

- DON'T type *into* the repl. Type into a file that you can search later.

- Use repl-interaction commands (like insert-in-repl, minibuffer, etc.).

- Have a plan to avoid wasting time (since repls are fun).

- Know *all* the jump-to-documentation kinds of commands. Searching for a regex,  
  jumping to the source of a thing, etc.

- Analysing data carefully at the repl (don't just stare at stacktraces).

- Debugging at O(log n)

- Custom REPLs (especially debuggers/inspectors with breakpoints, etc.)

- Do use the debugger.

#### MOST important keystrokes

- Execute a thing in the minibuffer

- Tracing functions (or a short function to execute in the minibuffer like =trace=)

- Sending to repl from behind the point

- Sending to repl from in front of the point

- Sending the whole top-level expression at point to the repl

- Jump to the minibuffer (not for typing! just to see big output)

- A keystroke to jump back to the source file (or something like windmove)

- *all* the documentation commands (search, jump to source, doc-at-point, etc.)

- Execute something and replace it with the output/write the output as

a comment in the scratch buffer.
