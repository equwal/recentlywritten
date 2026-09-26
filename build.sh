#!/bin/sh
# build.sh — converts posts/*.md and pages/*.md to site/ using pandoc
# Usage: sh build.sh [--no-deploy]

set -e

POSTS_DIR="posts"
PAGES_DIR="pages"
SITE_DIR="site"
STATIC_DIR="static"
SITE_URL="https://recentlywritten.com"
DEPLOY_HOST="root@recentlywritten.com"
DEPLOY_PATH="/var/www/recentlywritten/"
FEED_SIZE=10
DICKT_FEED="https://dickt.store/rss.xml"

DEPLOY=yes
[ "$1" = "--no-deploy" ] && DEPLOY=no

WORK="${TMPDIR:-/tmp}/rw-build.$$"
mkdir -p "$WORK"
trap 'rm -rf "$WORK"' EXIT INT TERM

# ── setup ──────────────────────────────────────────────────
rm -rf "$SITE_DIR"
mkdir -p "$SITE_DIR"
cp style.css "$SITE_DIR/style.css"
# Images and downloads the recovered posts link to as static/...
[ -d "$STATIC_DIR" ] && cp -R "$STATIC_DIR" "$SITE_DIR/$STATIC_DIR"

# pandoc HTML template for individual posts.
# Kept as a relative path in the working directory on purpose: under cygwin
# the pandoc on PATH is the native Windows build, which cannot resolve a
# POSIX path like /tmp/..., but does resolve a path relative to the cwd.
TEMPLATE=".build-template.html"
trap 'rm -rf "$WORK" "$TEMPLATE"' EXIT INT TERM
cat > "$TEMPLATE" << 'TMPL'
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>$title$ — Recently Written</title>
  <meta property="og:type" content="article" />
  <meta property="og:site_name" content="Recently Written" />
  <meta property="og:title" content="$title$" />
  <meta property="og:url" content="$url$" />
  <meta property="og:image" content="https://recentlywritten.com/static/og.png" />
  <meta name="twitter:card" content="summary_large_image" />
  <link rel="stylesheet" href="style.css" />
  <link rel="alternate" type="application/rss+xml" title="Recently Written" href="rss.xml" />
$if(feed)$
  <link rel="alternate" type="application/atom+xml" title="$title$" href="$feed$" />
$endif$
</head>
<body>
<div id="container">

  <div id="top">
    <a class="site-title" href="index.html">Recently Written</a>
    <nav>
      <a href="index.html">Home</a>
      <a href="about.html">about</a>
      <a href="lair.html">code</a>
      <a href="projects.html">projects</a>
      <a href="git/index.html">git</a>
      <a href="links.html">links</a>
      <a href="youtube.html">youtube</a>
      <a href="esperanto.html">esperanto</a>
      <a href="call.html">contact</a>
      <a href="rss.xml">rss</a>
      <a href="account/">subscribe</a>
    </nav>
  </div>

  <h1 class="post-title">$title$</h1>

  $body$

$if(links)$
  <p class="links">$for(links)$$links$$sep$ · $endfor$</p>

$endif$
  <div id="footer">
    <a href="index.html">← Home</a> ·
    <a href="https://github.com/equwal">Github</a>
    <p class="more">More sites and apps:
      <a href="https://subread.space/" title="Read along with an audiobook.">SubRead</a> ·
      <a href="https://booksimulator.com/" title="A reading room for Aozora Bunko and Project Gutenberg books.">Book Simulator</a> ·
      <a href="https://honjimaku.com/" title="Subtitles for Japanese audiobooks.">honjimaku.com</a> ·
      <a href="https://sbmsync.com/" title="Your bookmarks, the same on every device.">sbm Sync</a> ·
      <a href="https://dickt.store/" title="Language-learning tools, flashcards and web toys.">dickt.store</a> ·
      <a href="https://hentaibun.online/" title="Learn kanbun and kobun.">hentaibun.online</a> ·
      <a href="https://mowimypopol.online/" title="Archival study materials for Polish.">mowimypopol.online</a> ·
      <a href="https://alljapanesealltheti.me/" title="All Japanese All The Time.">AJATT mirror</a></p>
    <p class="more">Apps:
      <a href="https://github.com/equwal/subread-android/releases/latest">SubRead for Android</a> ·
      <a href="https://github.com/equwal/subread-overlay/releases/latest">SubRead Overlay</a> ·
      <a href="https://github.com/equwal/subread-dictionary/releases/latest">SubRead Dictionary</a> ·
      <a href="https://github.com/equwal/subread-anki">SubRead Anki</a> ·
      <a href="https://github.com/equwal/subread-extension/releases/latest">SubRead for YouTube</a> ·
      <a href="https://github.com/equwal/subread.koplugin">SubRead for KOReader</a> ·
      <a href="https://github.com/equwal/subrep-android/releases/latest">Subrep</a> ·
      <a href="https://github.com/equwal/sbm-android/releases/latest">sbm for Android</a> ·
      <a href="https://github.com/equwal/sbm-extension/releases/latest">sbm add-on</a> ·
      <a href="https://github.com/equwal/rebind/releases">Rebind</a> ·
      <a href="https://github.com/equwal/ink-recents/releases/latest">Ink Recents</a> ·
      <a href="https://github.com/equwal/ink-dim/releases/latest">Ink Dim</a> ·
      <a href="https://github.com/equwal/ink-update/releases/latest">Ink Update</a> ·
      <a href="projects.html">All projects</a></p>
  </div>

</div>
</body>
</html>
TMPL

# Read one `key: value` line out of a file's front matter, dropping the
# surrounding quotes that titles containing a colon have to be written with.
meta() {
    sed -n "s/^$2: *//p" "$1" | head -1 | sed -e 's/^"//' -e 's/"$//' -e 's/\\"/"/g'
}

# Escape the characters that must not appear raw in generated HTML.
esc() {
    sed -e 's/&/\&amp;/g' -e 's/</\&lt;/g' -e 's/>/\&gt;/g'
}

# render <markdown> <output> <title>
render() {
    pandoc \
        --from markdown \
        --to html5 \
        --template "$TEMPLATE" \
        --metadata title="$3" \
        --metadata url="$SITE_URL/$(basename "$2")" \
        --output "$2" \
        "$1"
}

# ── build each post ─────────────────────────────────────────
: > "$WORK/posts.tsv"
for md in "$POSTS_DIR"/*.md; do
    [ -e "$md" ] || continue
    slug=$(basename "$md" .md)
    title=$(meta "$md" title)
    date=$(meta "$md" date)
    # Recovered posts carry `order`: their position on the old front page,
    # the only surviving record of their sequence. Everything else sorts by
    # date. The two never interleave, since an order is three digits ("051")
    # and a date leads with its year ("2026-..."), so new posts land on top.
    order=$(meta "$md" order)
    [ -z "$order" ] && order="$date"
    render "$md" "$SITE_DIR/${slug}.html" "$title"
    printf '%s\t%s\t%s\t%s\n' "$order" "$date" "$slug" "$title" >> "$WORK/posts.tsv"
done

# Dates order the posts but are never shown on the site: most recovered ones
# are approximations taken from archive crawls, good enough to sort by and
# not good enough to publish. C collation keeps the sort independent of the
# system language.
LC_ALL=C sort -r "$WORK/posts.tsv" > "$WORK/sorted.tsv"

# ── build each standalone page ──────────────────────────────
# Pages are dateless and stay out of the chronological list; these are the
# hub pages the old site linked from its nav (lair, esperanto, call, ...).
: > "$WORK/pages.tsv"
for md in "$PAGES_DIR"/*.md; do
    [ -e "$md" ] || continue
    slug=$(basename "$md" .md)
    title=$(meta "$md" title)
    render "$md" "$SITE_DIR/${slug}.html" "$title"
    printf '%s\t%s\n' "$slug" "$title" >> "$WORK/pages.tsv"
done

# ── list fragments for the index ────────────────────────────
TAB=$(printf '\t')

while IFS="$TAB" read -r order date slug title; do
    [ -z "$slug" ] && continue
    printf '<li><a href="%s.html">%s</a></li>\n' \
        "$slug" "$(printf '%s' "$title" | esc)"
done < "$WORK/sorted.tsv" > "$WORK/post-items.html"

LC_ALL=C sort -f "$WORK/pages.tsv" | while IFS="$TAB" read -r slug title; do
    [ -z "$slug" ] && continue
    printf '<li><a href="%s.html">%s</a></li>\n' \
        "$slug" "$(printf '%s' "$title" | esc)"
done > "$WORK/page-items.html"

# ── build index ─────────────────────────────────────────────
# The fragments are read from disk rather than substituted into an awk
# variable, so titles containing & or \ survive intact.
awk -v postfile="$WORK/post-items.html" -v pagefile="$WORK/page-items.html" '
    /<!-- POST_LIST -->/ {
        print "<ul class=\"post-list\">"
        while ((getline line < postfile) > 0) print line
        close(postfile)
        print "</ul>"
        next
    }
    /<!-- PAGE_LIST -->/ {
        print "<ul class=\"page-list\">"
        while ((getline line < pagefile) > 0) print line
        close(pagefile)
        print "</ul>"
        next
    }
    { print }
' index.html > "$SITE_DIR/index.html"

# ── feed ─────────────────────────────────────────────────────
# Only the newest FEED_SIZE posts go in, and that is deliberate as well as
# conventional: a feed reader shows each item's date, and the newest posts are
# the ones whose dates are exact rather than recovered from a crawl.

# RFC 822, as RSS requires. The C locale matters: without it the day and
# month names come out in the system language. `date -d` is GNU; where it is
# missing the item simply goes out without a pubDate. A post has only a
# date, so its time is midnight; a podcast episode has an exact time.
feed_date() {
    LC_ALL=C date -u -d "$1" '+%a, %d %b %Y %H:%M:%S +0000' 2>/dev/null
}

# A feed reader shows the HTML away from the site, so links relative to it
# need the site put in front of them.
absolute() {
    sed -E \
        -e 's@(href|src)="/([^/"][^"]*)"@\1="'"$SITE_URL"'/\2"@g' \
        -e 's@(href|src)="([^"/#][^":]*)"@\1="'"$SITE_URL"'/\2"@g'
}

{
    printf '<?xml version="1.0" encoding="UTF-8"?>\n'
    printf '<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">\n<channel>\n'
    printf '  <title>Recently Written</title>\n'
    printf '  <link>%s/</link>\n' "$SITE_URL"
    printf '  <description>A site full of things which I have recently written.</description>\n'
    printf '  <language>en-us</language>\n'
    printf '  <atom:link href="%s/rss.xml" rel="self" type="application/rss+xml" />\n' "$SITE_URL"
    head -n "$FEED_SIZE" "$WORK/sorted.tsv" | while IFS="$TAB" read -r order date slug title; do
        url="$SITE_URL/$slug.html"
        printf '  <item>\n'
        printf '    <title>%s</title>\n' "$(printf '%s' "$title" | esc)"
        printf '    <link>%s</link>\n' "$url"
        printf '    <guid isPermaLink="true">%s</guid>\n' "$url"
        pub=$(feed_date "$date") && [ -n "$pub" ] &&
            printf '    <pubDate>%s</pubDate>\n' "$pub"
        printf '    <description><![CDATA['
        pandoc --from markdown --to html5 "$POSTS_DIR/$slug.md" |
            absolute | sed 's/]]>/]]]]><![CDATA[>/g'
        printf ']]></description>\n'
        printf '  </item>\n'
    done
    # Releases on dickt.store go in the same feed. The build copies each
    # item from the live dickt.store feed. If the fetch fails, the build
    # continues without these items, because the posts are more important.
    if curl -fsS -o "$WORK/dickt.xml" "$DICKT_FEED"; then
        awk '/<item>/,/<\/item>/' "$WORK/dickt.xml"
    else
        echo "warning: no items from $DICKT_FEED" >&2
    fi
    printf '</channel>\n</rss>\n'
} > "$SITE_DIR/rss.xml"

# ── podcast ──────────────────────────────────────────────────
# podcast.tsv lists the podcast episodes, newest first, one on each line:
#   date <TAB> file <TAB> bytes <TAB> seconds <TAB> youtube <TAB> title
# date is the UTC time of the upload in ISO 8601. file is the name of the
# audio in static/language-podcast/. The audio stays out of git (see
# .gitignore), so its size and length, which the feed must give, are kept
# here. youtube is the ID of the same episode on YouTube, or "-" if none.
# Every field needs a value, because read merges empty tab-separated fields.
PODCAST_PAGE="polyglot-podcast.html"
AUDIO_DIR="static/language-podcast"

youtube_url() {
    printf 'https://www.youtube.com/watch?v=%s' "$1"
}

{
    printf '<?xml version="1.0" encoding="UTF-8"?>\n'
    printf '<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom"'
    printf ' xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd"'
    printf ' xmlns:podcast="https://podcastindex.org/namespace/1.0">\n<channel>\n'
    printf '  <title>Polyglot Podcast</title>\n'
    printf '  <link>%s/%s</link>\n' "$SITE_URL" "$PODCAST_PAGE"
    printf '  <description>Spenser Truex talks about how to learn languages, and about AI.</description>\n'
    printf '  <language>en-us</language>\n'
    printf '  <atom:link href="%s/podcast.xml" rel="self" type="application/rss+xml" />\n' "$SITE_URL"
    printf '  <itunes:author>Spenser Truex</itunes:author>\n'
    # Podcast directories such as Spotify send the code that proves
    # ownership of the feed to this email address.
    printf '  <itunes:owner>\n'
    printf '    <itunes:name>Spenser Truex</itunes:name>\n'
    printf '    <itunes:email>truex@equwal.com</itunes:email>\n'
    printf '  </itunes:owner>\n'
    printf '  <itunes:image href="%s/static/podcast.png" />\n' "$SITE_URL"
    printf '  <itunes:category text="Education">\n'
    printf '    <itunes:category text="Language Learning" />\n'
    printf '  </itunes:category>\n'
    printf '  <itunes:explicit>false</itunes:explicit>\n'
    printf '  <podcast:funding url="https://ko-fi.com/truex">Support the podcast</podcast:funding>\n'
    while IFS="$TAB" read -r date file bytes seconds youtube title; do
        [ -z "$date" ] && continue
        case $file in
            *.m4a) type=audio/x-m4a ;;
            *.mp3) type=audio/mpeg ;;
            *) echo "podcast.tsv: no MIME type for $file" >&2; exit 1 ;;
        esac
        url="$SITE_URL/$AUDIO_DIR/$file"
        printf '  <item>\n'
        printf '    <title>%s</title>\n' "$(printf '%s' "$title" | esc)"
        printf '    <enclosure url="%s" length="%s" type="%s" />\n' "$url" "$bytes" "$type"
        printf '    <guid isPermaLink="false">%s</guid>\n' "$url"
        pub=$(feed_date "$date") && [ -n "$pub" ] &&
            printf '    <pubDate>%s</pubDate>\n' "$pub"
        printf '    <itunes:duration>%s</itunes:duration>\n' "$seconds"
        if [ "$youtube" != - ]; then
            printf '    <link>%s</link>\n' "$(youtube_url "$youtube")"
            printf '    <description>Video: %s</description>\n' "$(youtube_url "$youtube")"
        fi
        printf '  </item>\n'
    done < podcast.tsv
    printf '</channel>\n</rss>\n'
} > "$SITE_DIR/podcast.xml"

# The podcast page lists the same episodes. The list goes in place of the
# comment <!-- EPISODE_LIST --> in pages/polyglot-podcast.md.
while IFS="$TAB" read -r date file bytes seconds youtube title; do
    [ -z "$date" ] && continue
    printf '<li><a href="%s/%s">%s</a>' "$AUDIO_DIR" "$file" \
        "$(printf '%s' "$title" | esc)"
    [ "$youtube" = - ] ||
        printf ' (<a href="%s">video</a>)' "$(youtube_url "$youtube")"
    printf '</li>\n'
done < podcast.tsv > "$WORK/episode-items.html"
awk -v listfile="$WORK/episode-items.html" '
    /<!-- EPISODE_LIST -->/ {
        print "<ul>"
        while ((getline line < listfile) > 0) print line
        close(listfile)
        print "</ul>"
        next
    }
    { print }
' "$SITE_DIR/$PODCAST_PAGE" > "$WORK/podcast-page.html"
cp "$WORK/podcast-page.html" "$SITE_DIR/$PODCAST_PAGE"

# ── post data for the accounts service ──────────────────────
# accounts/ reads these files to send email about new posts and to make
# the private feed of each user. index.tsv has one line for each post,
# newest first: date, slug, tags, title. postdata/<slug>.html has the body
# of the post, with absolute links.
mkdir -p "$SITE_DIR/postdata"
while IFS="$TAB" read -r order date slug title; do
    [ -z "$slug" ] && continue
    printf '%s\t%s\t%s\t%s\n' "$date" "$slug" "$(meta "$POSTS_DIR/$slug.md" tags)" "$title"
    pandoc --from markdown --to html5 "$POSTS_DIR/$slug.md" < /dev/null |
        absolute > "$SITE_DIR/postdata/$slug.html"
done < "$WORK/sorted.tsv" > "$SITE_DIR/postdata/index.tsv"

echo "Built $(ls "$SITE_DIR"/*.html | wc -l | tr -d ' ') pages → $SITE_DIR/"
echo "  posts: $(wc -l < "$WORK/posts.tsv" | tr -d ' ')   pages: $(wc -l < "$WORK/pages.tsv" | tr -d ' ')   feed: $(grep -c '<item>' "$SITE_DIR/rss.xml") items   podcast: $(grep -c '<item>' "$SITE_DIR/podcast.xml") episodes"

# ── deploy ───────────────────────────────────────────────────
# git/ and git.html share this web root but are published by deploy-git.sh,
# not built here. Excluding them keeps --delete from treating them as stale:
# rsync never deletes an excluded path on the receiving side.
# static/book/ holds book files that stay out of the repo (see .gitignore).
# They exist only on the server, so the same rule protects them.
# static/language-podcast/ holds the podcast audio, which also stays out of
# the repo. The P (protect) rule stops --delete from removing the server
# copy when the working tree has no audio, but the audio that the working
# tree has still uploads.
if [ "$DEPLOY" = yes ]; then
    rsync -avzP --delete --exclude=/git/ --exclude=/git.html \
        --exclude=/static/book/ --filter='P /static/language-podcast/***' \
        "$SITE_DIR/" "$DEPLOY_HOST:$DEPLOY_PATH"
    ssh "$DEPLOY_HOST" "chmod -R a+rX $DEPLOY_PATH"
fi
