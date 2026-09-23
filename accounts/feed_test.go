package main

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"slices"
	"strings"
	"testing"

	"pgregory.net/rapid"
)

// wellFormed reads each token of doc with a strict XML parser.
func wellFormed(doc []byte) error {
	d := xml.NewDecoder(bytes.NewReader(doc))
	for {
		_, err := d.Token()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

// anyText draws any string, also strings that are not valid UTF-8 and
// strings with characters that XML 1.0 does not permit.
func anyText() *rapid.Generator[string] {
	return rapid.OneOf(
		rapid.String(),
		rapid.Map(rapid.SliceOf(rapid.Byte()), func(b []byte) string { return string(b) }),
		rapid.SampledFrom([]string{"]]>", "<![CDATA[", "&amp;", "</item>", "\x00", "\r\n", "\uFFFE"}),
	)
}

// xmlText draws strings of characters that XML 1.0 permits.
func xmlText() *rapid.Generator[string] {
	return rapid.StringOf(rapid.Rune().Filter(func(r rune) bool {
		return r == '\t' || r == '\n' || r == '\r' ||
			r >= 0x20 && r <= 0xD7FF || r >= 0xE000 && r <= 0xFFFD || r >= 0x10000 && r <= 0x10FFFF
	}))
}

func drawChannel(t *rapid.T, text *rapid.Generator[string]) channel {
	ch := channel{Title: text.Draw(t, "title"), Link: text.Draw(t, "link"),
		Description: text.Draw(t, "desc"), Language: "en-us"}
	for range rapid.IntRange(0, 3).Draw(t, "items") {
		ch.Items = append(ch.Items, item{
			Title:       text.Draw(t, "item title"),
			Link:        text.Draw(t, "item link"),
			GUID:        guid{IsPermaLink: "true", Value: text.Draw(t, "guid")},
			PubDate:     text.Draw(t, "date"),
			Categories:  rapid.SliceOfN(text, 0, 2).Draw(t, "tags"),
			Description: text.Draw(t, "body"),
		})
	}
	return ch
}

func TestFeedIsWellFormedProperty(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		var b bytes.Buffer
		if err := writeFeed(&b, drawChannel(t, anyText())); err != nil {
			t.Fatal(err)
		}
		if err := wellFormed(b.Bytes()); err != nil {
			t.Fatalf("feed is not well-formed: %v\n%s", err, b.String())
		}
	})
}

// TestFeedRoundTripProperty checks the escaping: an XML parser gives back
// each text exactly.
func TestFeedRoundTripProperty(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		want := drawChannel(t, xmlText())
		var b bytes.Buffer
		if err := writeFeed(&b, want); err != nil {
			t.Fatal(err)
		}
		var got rss
		if err := xml.Unmarshal(b.Bytes(), &got); err != nil {
			t.Fatalf("unmarshal: %v\n%s", err, b.String())
		}
		g := got.Channel
		if g.Title != want.Title || g.Link != want.Link || g.Description != want.Description ||
			len(g.Items) != len(want.Items) {
			t.Fatalf("channel: got %+v, want %+v", g, want)
		}
		for i, w := range want.Items {
			gi := g.Items[i]
			if gi.Title != w.Title || gi.Link != w.Link || gi.GUID != w.GUID ||
				gi.PubDate != w.PubDate || gi.Description != w.Description ||
				!slices.Equal(gi.Categories, w.Categories) {
				t.Fatalf("item %d: got %+v, want %+v", i, gi, w)
			}
		}
	})
}

func TestBuildFeedFiltersAndLimits(t *testing.T) {
	var posts []Post
	for i := range 30 {
		tag := "articles"
		if i%2 == 1 {
			tag = "projects"
		}
		posts = append(posts, Post{Date: "2026-09-01", Slug: "p" + string(rune('a'+i)), Tags: []string{tag}, Title: "T"})
	}
	posts = append(posts, Post{Date: "bad", Slug: "untagged", Title: "U"})

	ch := buildFeed("https://example.com", posts, map[string]bool{"projects": true}, func(string) string { return "" })
	if len(ch.Items) != 16 { // 15 articles and the untagged post
		t.Fatalf("got %d items, want 16", len(ch.Items))
	}
	for _, it := range ch.Items {
		if slices.Contains(it.Categories, "projects") {
			t.Errorf("muted tag in feed: %+v", it)
		}
	}
	if got := ch.Items[0].PubDate; got != "Tue, 01 Sep 2026 00:00:00 +0000" {
		t.Errorf("pubDate = %q", got)
	}
	if got := ch.Items[15].PubDate; got != "" {
		t.Errorf("pubDate for a bad date = %q, want empty", got)
	}

	ch = buildFeed("https://example.com", posts, nil, func(string) string { return "" })
	if len(ch.Items) != feedSize {
		t.Fatalf("got %d items, want %d", len(ch.Items), feedSize)
	}
}

func TestParseIndex(t *testing.T) {
	in := "2026-09-22\trebind\tprojects\tRebind: a tool\r\n\n2025-01-01\tx\tlisp, unix\tA \"b\" & c\n"
	posts, err := parseIndex(strings.NewReader(in))
	if err != nil {
		t.Fatal(err)
	}
	want := []Post{
		{Date: "2026-09-22", Slug: "rebind", Tags: []string{"projects"}, Title: "Rebind: a tool"},
		{Date: "2025-01-01", Slug: "x", Tags: []string{"lisp", "unix"}, Title: `A "b" & c`},
	}
	if len(posts) != len(want) {
		t.Fatalf("got %+v", posts)
	}
	for i := range want {
		if posts[i].Slug != want[i].Slug || posts[i].Title != want[i].Title ||
			posts[i].Date != want[i].Date || !slices.Equal(posts[i].Tags, want[i].Tags) {
			t.Errorf("post %d = %+v, want %+v", i, posts[i], want[i])
		}
	}

	for _, bad := range []string{"a\t../etc/passwd\tx\tT\n", "a\t.hidden\tx\tT\n", "only\tthree\tfields\n"} {
		if _, err := parseIndex(strings.NewReader(bad)); err == nil {
			t.Errorf("parseIndex(%q) gave no error", bad)
		}
	}
}
