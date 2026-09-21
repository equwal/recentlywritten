---
title: irc relay
source: therealtruex.com (recovered from web.archive.org)
---

oftc over tor

## oftc.net IRC onion

Note: this relay is setup with permission from oftc.net

The relay is at:  
  
` equwalnexhy7aevc7wyiqz6ivwlaegabdwj4s6hvhcnh6w2umdjnsfqd.onion`  
`Port 6667`  

Connecting is done without SSL (since Tor does that already), and the relay does use SSL on the way out.

### Client setup

#### socat tunnel

` socat TCP-LISTEN:6691 \ # note: not SSL!`  
`SOCKS4a:localhost:equwalnexhy7aevc7wyiqz6ivwlaegabdwj4s6hvhcnh6w2umdjnsfqd.onion:6667,socksport=9050`  

#### connect ii to it

This is how I do connect to it with the [ii](https://tools.suckless.org/ii/) client:  
  
` ii -s 127.0.0.1 -p 6691 -n equwal`  

### Server setup

This is how the relay works, if you want to setup your own relay server.  
  
/etc/init.d/socat-oftc  
  
` socat TCP-LISTEN:6667,reuseaddr,fork OPENSSL-CONNECT:irc.oftc.net:6697`  
  
/etc/tor/torrc  
  
` SocksPort 0`  
`HiddenServicePort 6667 127.0.0.1:6667 # Standard stuff`  
`# Make things go fast, but my server non-anonymous`  
`HiddenServiceSingleHopMode 1`  
`HiddenServiceNonAnonymousMode 1`  
  
I would like to be able to forward end-to-end encrypted SSL connections (eg. with SASL for authentication), but I haven't figured out how to do that yet.
