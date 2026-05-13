#!/bin/sh
# build.sh — converts posts/*.md to site/ using pandoc
# Usage: sh build.sh

set -e

POSTS_DIR="posts"
SITE_DIR="site"
CSS="../style.css"   # relative path from inside site/

# ── setup ──────────────────────────────────────────────────
rm -rf "$SITE_DIR"
mkdir -p "$SITE_DIR"
cp style.css "$SITE_DIR/style.css"

# pandoc HTML template for individual posts
TEMPLATE="template.html"
cat > "$TEMPLATE" << 'TMPL'
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>$title$ — Recently Written</title>
  <link rel="stylesheet" href="style.css" />
</head>
<body>
<div id="container">

  <div id="top">
    <a class="site-title" href="index.html">Recently Written</a>
    <nav>
      <a href="index.html">Home</a>
    </nav>
  </div>

  <h1>$title$</h1>
  <div class="post-meta">$date$</div>

  $body$

  <div id="footer">
    <a href="index.html">← Home</a> ·
    Powered by <a href="https://neuron.zettel.page">Neuron</a>
  </div>

</div>
</body>
</html>
TMPL

# ── build each post ─────────────────────────────────────────
LIST=""   # accumulate post-list items (newest first)

# collect posts and sort by date descending
for md in "$POSTS_DIR"/*.md; do
    slug=$(basename "$md" .md)
    title=$(grep '^title:' "$md" | sed 's/title: *//')
    date=$(grep '^date:'  "$md" | sed 's/date: *//')

    pandoc \
        --from markdown \
        --to html5 \
        --template "$TEMPLATE" \
        --metadata title="$title" \
        --metadata date="$date" \
        --output "$SITE_DIR/${slug}.html" \
        "$md"

    # prepend to list so newest-first after sort
    LIST="$(printf '%s\n' "${date}|${slug}|${title}")\n${LIST}"
done

# ── build index ─────────────────────────────────────────────
# sort descending by date, emit list items
POST_ITEMS=""
printf '%b' "$LIST" | sort -r | while IFS='|' read -r date slug title; do
    [ -z "$slug" ] && continue
    printf '<li><span class="post-date">%s</span><a href="%s.html">%s</a></li>\n' \
        "$date" "$slug" "$title"
done > /tmp/rw-items.txt

ITEMS=$(cat /tmp/rw-items.txt)

# inject post list into index.html using awk
awk -v items="<ul class=\"post-list\">$(cat /tmp/rw-items.txt)</ul>" \
    '{ gsub(/<!-- POST_LIST -->/, items); print }' \
    index.html > "$SITE_DIR/index.html"

# ── cleanup ──────────────────────────────────────────────────
rm -f "$TEMPLATE" /tmp/rw-items.txt

echo "Built $(ls "$SITE_DIR"/*.html | wc -l | tr -d ' ') pages → $SITE_DIR/"
rsync -avzP --delete site/ root@honjimaku.com:/var/www/recentlywritten/
