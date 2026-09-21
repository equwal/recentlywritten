---
title: IRC Bouncer
date: 2018-12-02
order: 002
tags: projects
source: "therealtruex.com (recovered from web.archive.org; date exact, from the original .post file in equwal/truex.eu)"
---

I no longer operate a ZNC bouncer. I now operate an ii logger though.

Currently only plan on supporting libera chat connections, as a replacement for Firrre, who have closed down user signups for the time being.It is available at therealtruex.com:6969 to registered users. There is a web panel to add and remove networks and channels.

### Using the Web Panel

![Your Settings link.](static/yoursettings.png)

![Networks](static/networks.png)

![Channels](static/channels.png)

![Adding a channel](static/channel.png)

### Loggin in from the IRC Client

You might like Hexchat as an IRC client.

Connect to your desired network from the client using the magic username combination:

    username/networkname

the username is your ZNC user (not the nick for the network) and the networkname is the name of the network configured in ZNC (not the network's domain).![Hexchat config panel](static/hexchat.png)
