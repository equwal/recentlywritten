#!/bin/sh
# deploy-git.sh - regenerate the stagit-style pages and publish them.
# Run from the repo root, in Git Bash / any POSIX shell (needs git, gh, python3, ssh, tar).
#   sh deploy-git.sh
set -e
HOST=root@honjimaku.com
WEBROOT=/var/www/recentlywritten

python3 stagit-gen.py equwal git "${TMPDIR:-${TEMP:-/tmp}}/stagit-cache"

# Replace only git/ and git.html; nothing else in the webroot is touched.
ssh "$HOST" "rm -rf $WEBROOT/git"
tar cf - git git.html | ssh "$HOST" "cd $WEBROOT && tar xf - --no-same-owner \
  && chown -R 197108:197121 git git.html \
  && find git -type d -exec chmod 775 {} + && find git -type f -exec chmod 664 {} + \
  && chmod 664 git.html"
echo "published: https://recentlywritten.com/git/"
