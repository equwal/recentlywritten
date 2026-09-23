---
title: Projects
---

Spenser Truex. Free tarballs.

Tools for language learners, e-ink readers and Android, small Unix programs,
and Common Lisp.
Support the work at [ko-fi.com/truex](https://ko-fi.com/truex).

## Sites

* [subread.space](https://subread.space) - Line up an audiobook with its ebook. Get read-along subtitles, EPUB 3 read-along books, and subtitled video. Free in the browser.
* [honjimaku.com](https://honjimaku.com) - A library of Japanese subtitles for audiobooks, and Subrep live captions by share link.
* [sbmsync.com](https://sbmsync.com) - sbm Sync: one plain bookmark file, the same on every device. Free, AGPL.
* [hentaibun.online](https://hentaibun.online) - Reading lists and tools to learn kanbun and kobun.
* [recentlywritten.com](https://recentlywritten.com) - Writing on Unix, Lisp, Nix, Esperanto, and language learning. ([equwal.com](https://equwal.com) goes here.)
* [ko-fi.com/truex](https://ko-fi.com/truex) - Support.

## SubRead

Read along with audiobooks, and mine the words.

* [subread.space](https://github.com/equwal/subread.space) - Read an `.srt` one line at a time with the space bar. Works with Yomitan, furigana, and pitch accent markup.
* [subplz-web](https://github.com/equwal/subplz-web) - The aligner behind subread.space. The subtitle text comes from your book, not from a transcript.
* [subread-android](https://github.com/equwal/subread-android) - SubRead on Android. Times an audiobook against its ebook on the device.
* [subread-extension](https://github.com/equwal/subread-extension) - Show your own `.srt` over a YouTube video, in time. Chrome and Firefox.
* [subread-overlay](https://github.com/equwal/subread-overlay) - Show `.srt` lines over any Android media player. Tap a word to look it up.
* [subread-dictionary](https://github.com/equwal/subread-dictionary) - Pop-up dictionary for Android. Reads Yomitan dictionaries, with local audio.
* [subread.koplugin](https://github.com/equwal/subread.koplugin) - KOReader plugin. The book follows the narration, from an `.srt` made by subread.space.
* [subrep-android](https://github.com/equwal/subrep-android) - Live captions of the phone's sound, made on the phone by Whisper, shared by link.
* [desktop-subtitle-replay](https://github.com/equwal/desktop-subtitle-replay) - Live subtitles for anything on your screen. Replay a sentence and mine it into Anki.
* [honjimaku](https://github.com/equwal/honjimaku) - The honjimaku.com server. Fork of [jimaku](https://github.com/Rapptz/jimaku).
* [kikiyomi](https://github.com/equwal/kikiyomi) - Audiobook player for Japanese learners, in the browser. Fork of [kikiyomi](https://github.com/rtr46/kikiyomi).
* [Hoshi-Reader-Android](https://github.com/equwal/Hoshi-Reader-Android) - Fork with a `subread` branch: audiobook read-along in the Hoshi Reader EPUB reader.
* [chimahon](https://github.com/equwal/chimahon) - Fork with a `readalong` branch for the Mihon immersion reader.
* [dickt.store](https://github.com/equwal/dickt.store) - Bring-your-own-key Yomitan dictionary converter.
* [vibeslop-dickt](https://github.com/equwal/vibeslop-dickt) - Make your own Yomitan dictionary in a few minutes.

## Rebind suite

Hardware buttons for e-ink readers and Android 12+.

* [rebind](https://github.com/equwal/rebind) - Remap volume keys, Power, page-turn buttons, and the Viwoods AI key. Tap, double tap, hold, and combinations. Free beta APK.
* [assistkey](https://github.com/equwal/assistkey) - Key remapper for the Viwoods AiPaper Reader: AI key, volume keys, and Power.
* [awesome-rebind](https://github.com/equwal/awesome-rebind) - Extensions and companion apps for Rebind.
* [ink-recents](https://github.com/equwal/ink-recents) - Recent-apps switcher for e-ink. No animation. Swipe up to close.
* [ink-dim](https://github.com/equwal/ink-dim) - Set the Viwoods frontlight below the lowest system level.
* [ink-update](https://github.com/equwal/ink-update) - Tells you when Rebind, Ink Recents, or Ink Dim has a new version.

## sbm

Fuzzy search thousands of bookmarks, on every device.

* [sbm](https://github.com/equwal/sbm) - dmenu bookmarks. Search and plumb them.
* [sbm-android](https://github.com/equwal/sbm-android) - sbm on Android. Share to add.
* [sbm-extension](https://github.com/equwal/sbm-extension) - sbm for Firefox and Chrome.
* [sbm-sync](https://github.com/equwal/sbm-sync) - Sync server. One small Go program, plain files, AGPL.
* [sbm-suckless](https://github.com/equwal/sbm-suckless) - sbm before the rewrite. Suckless style.

## Tools for AI agents

* [ideamine](https://github.com/equwal/ideamine) - Idea inbox for Claude Code. `/idea` saves ideas at zero tokens.
* [deploy-handoff](https://github.com/equwal/deploy-handoff) - The agent drives the browser to the last deploy, billing, or pull request step. You make the final click.
* [memstate](https://github.com/equwal/memstate) - Versioned memory for AI agents in one SQLite file. Fork; contributor to [map588/memstate](https://github.com/map588/memstate).
* [aintitinit](https://github.com/HACKERMEN-ORG/aintitinit) - Claude initializer.
* [hackermen-neovim-config-file](https://github.com/HACKERMEN-ORG/hackermen-neovim-config-file) - Minimal Neovim config with optional AI.

## Unix and suckless

* [tmpl](https://github.com/equwal/tmpl) - Insanely simple templates.
* [dmenu-scripts](https://github.com/equwal/dmenu-scripts) - Ask for things with dmenu.
* [dsad](https://github.com/equwal/dsad) - dmenu interface for sad, the simple audio daemon.
* [dpatchmail](https://github.com/equwal/dpatchmail) - Patch programs from mail with dmenu.
* [prompt-mail](https://github.com/equwal/prompt-mail) - Mail questions and mail files with dmenu.
* [dmodurl](https://github.com/equwal/dmodurl) - Change URLs with dmenu.
* [tnot](https://github.com/equwal/tnot) - Like tmux, but not.
* [tiny-runlevels](https://github.com/equwal/tiny-runlevels) - Use an init system for many instances of simple programs.
* [ii-lchat-metapackage](https://github.com/equwal/ii-lchat-metapackage) - Get suckless IRC working.
* [pass-simple](https://github.com/equwal/pass-simple) - Minimal pass rewrite, with bug fixes.
* [sacc-ebuild](https://github.com/equwal/sacc-ebuild) - Gentoo ebuild for the sacc gopher client.
* [nixtoo-genux](https://github.com/equwal/nixtoo-genux) - Declarative Nix package manager setup.
* [xinput-persistence](https://github.com/equwal/xinput-persistence) - xinput settings that survive reconnect and boot.
* [encrypted-home-backups](https://github.com/equwal/encrypted-home-backups) - Automatic home backups to an encrypted partition.
* [alt-hack0](https://github.com/equwal/alt-hack0) and [alt-hack](https://github.com/equwal/alt-hack) - The Hack font with a normal zero, made by script.
* [minimak-tty-xkb](https://github.com/equwal/minimak-tty-xkb) - Minimak keyboard layout for the TTY and XKB.

## Common Lisp and Emacs

* [font-lock-cl](https://github.com/font-lock-cl/font-lock-cl) - Better Common Lisp syntax highlighting for Emacs. On MELPA.
* [coleslaw](https://github.com/equwal/coleslaw) - Emacs mode for the Coleslaw static site generator. Also [coleslaw-snippets](https://github.com/equwal/coleslaw-snippets).
* [Forthe](https://github.com/equwal/Forthe) - Forth-like DSL.
* [sl](https://github.com/equwal/sl) and [defsl](https://github.com/equwal/defsl) - Dependency language for Common Lisp.
* [sigil](https://github.com/equwal/sigil) - Documentation preprocessor for Common Lisp.
* [posix-pipes](https://github.com/equwal/posix-pipes) - Portable pipes from CLOCC.
* [tag-gen](https://github.com/equwal/tag-gen) - Make tags for each directory.
* [cores](https://github.com/equwal/cores) and [clapt](https://github.com/equwal/clapt) - Save SBCL cores with libraries loaded.
* [asdf-registration](https://github.com/equwal/asdf-registration) - Register ASDF projects.
* [CSV](https://github.com/equwal/CSV) - Read CSV into Lisp lists.
* [gendocs](https://github.com/equwal/gendocs) - Simple document generation.
* [len-cmp](https://github.com/equwal/len-cmp) - Fast list length comparison.
* [LispBrain](https://github.com/equwal/LispBrain) - Brainfuck debugger on the Common Lisp debugger.
* [lisp-multithreaded-internet](https://github.com/equwal/lisp-multithreaded-internet) - Multithreaded internet for Common Lisp.
* [esperanto](https://github.com/equwal/esperanto) - Compress natural language in Esperanto.
* [librivox](https://github.com/equwal/librivox) - Automatic YouTube upload of audiobooks.
* [cl-yag](https://github.com/equwal/cl-yag) and [clic](https://github.com/equwal/clic) - Mirrors of Bitreich Common Lisp projects.
* Spacemacs layers: [timeclock](https://github.com/equwal/timeclock), [rectangle](https://github.com/equwal/rectangle), [ediff-spacemacs](https://github.com/equwal/ediff-spacemacs), [slurp](https://github.com/equwal/slurp), [quickurl](https://github.com/equwal/quickurl), [common-lisp-sly](https://github.com/equwal/common-lisp-sly).

## Upstream contributions

* [suckless](https://suckless.org) - [slstatus](https://git.suckless.org/slstatus): "Not charging" battery status. [scroll](https://git.suckless.org/scroll): typo fix. Patches to dwm, dmenu, quark, and sacc on the mailing lists.
* [sad](https://git.2f30.org/sad/) (2f30) - `status volume` protocol, pause toggle, volume bug fix, simpler config macros.
* [SLY](https://github.com/joaotavora/sly) - REPL and debugger fixes, xref docs.
* [Spacemacs](https://github.com/syl20bnr/spacemacs) - Layers and bindings, cleanup of deprecated `cl` and `use-package` code.
* [Coleslaw](https://github.com/coleslaw-org/coleslaw) - URL generation, post extensions, versioning plugin, docs.
* [CLX](https://github.com/sharplispers/clx) - Composite extension docs and exports.
* [quickproject](https://github.com/xach/quickproject) - Restarts for existing target files.
* [ivy/swiper](https://github.com/abo-abo/swiper), [emacs-w3m](https://github.com/emacs-w3m/emacs-w3m), [Nameless](https://github.com/Malabarba/Nameless), [vlime](https://github.com/vlime/vlime), [Mezzano](https://github.com/froggey/Mezzano), [usocket](https://github.com/usocket/usocket), [Trucler](https://github.com/s-expressionists/Trucler), [Newport](https://github.com/Sharp-CLOCC/Newport).
* `make uninstall` targets: [mle](https://github.com/adsr/mle), [termbox](https://github.com/termbox/termbox), [wak](https://github.com/raygard/wak), [grabc](https://github.com/muquit/grabc), [sj](https://github.com/younix/sj), [gita](https://github.com/nosarthur/gita), [f3](https://github.com/AltraMayor/f3), [lisgd](https://github.com/PorQ-Pine/lisgd), [moreutils](https://github.com/pgdr/moreutils), [ames](https://github.com/eshrh/ames).
* [voidrice](https://github.com/LukeSmithxyz/voidrice), [the_silver_searcher](https://github.com/ggreer/the_silver_searcher), [github.fish](https://github.com/PatrickF1/github.fish), [asyncomplete.vim](https://github.com/prabirshrestha/asyncomplete.vim), [Akusento](https://github.com/tunamayo04/Akusento).
