---
title: How to Beat Emacs M-x dunnet in 5 minutes.
date: 2018-12-29
order: 003
tags: articles
source: "therealtruex.com (recovered from web.archive.org; date exact, from the original .post file in equwal/truex.eu)"
---

## What?

In the emacs text editor there are some goofy commands that play games. M-x (which means Alt+x) is the way to execute a command in the editor, so M-x dunnet starts the game called "dunnet" which is a text adventure. I cheated and won. Defining just one function and setting one variable is sufficient to win. \[Here is a text map of the easy path.\]([https://recentlywritten.com/static/path.txt)](static/path.txt)) The command prompt is a \`&gt;\`.

### Basic Commands

- `>take` something

- `>look`

- `>save` path

- `>restore` path

- Directions Are: n,s,e,w,ne,se,nw,sw,u,d,in,out (so you might `go n`).

### To win

First, save the game so you can restart easily later:

- `>save ~/fresh.dun`

How you need to move around the world a bit. Dunnet will greet you:

&lt;pre&gt;You are at a dead end of a dirt road. The road goes to the east. In the distance you can see that it will eventually fork off. The trees here are very tall royal palms, and they are spaced equidistant from each other. There is a shovel here.&lt;/pre&gt;

- `>take shovel`

- `>go e`

- `>go e`

- `>dig`

- `>look`

- `>take cpu`

- `>go se`

- `>take food`

- `>go se`

- `>feed bear food`

- `>look`

- `>take key`

- `>go nw`

- `>go nw`

- `>go ne`

- `>go ne`

- `>go ne`

- `>go w`

- `>put cpu in computer`

Now in M-: you must login to the VAX console

    Eval: (setq dun-logged-in t)

back to the dungeon:

- `>type`

greeted by a $ prompt:

To get to endgame in the VAX console, define this with M-:

    (defun dun-score (garb) 90)

Now back to the dungeon:  
-

    $ rlogin endgame

- `>go n`

Congrats, now you get questions, here they all are with answers:

- What is your password on the machine called 'pokey'? robert

- What password did you use during anonymous ftp to gamma? foo

- Excluding the endgame, how many places are there where you can put treasures for points? 4 or four

- What is your login name on the 'endgame' machine? toukmond

- What is the nearest whole dollar to the price of the shovel? 20 or twenty

- What is the name of the bus company serving the town? mobytours

- Give either of the two last names in the mailroom, other than your own. collier

- What cartoon character is on the towel? snoopy

- What is the last name of the author of EMACS? stallman

- How many megabytes of memory is on the CPU board for the Vax? 2

- Which street in town is named after a U.S. state? vermont

- How many pounds did the weight weigh? ten

- Name the STREET which runs right over the subway stop. 4th

- How many corners are there in town (excluding the one with the Post Office)? 24

- What type of bear was hiding your key? grizzly

- Name either of the two objects you found by digging. cpu

- What network protocol is used between pokey and gamma? ip

So you might do:

- `>answer cpu`

then:

- `>go n`

there are three quesions in total:

- `>answer` something

- `>go n`

- `>answer` something

- `>take bill`

(naturally, you deserve it you cheater)!

- `>go n`

- `>take Mona`

(you dirty cheater)

Congrats you are in the winner's room.

- `>save ~/cheater.dun`

- `>restore ~/fresh.dun`

Now you can actually play it if you want.
