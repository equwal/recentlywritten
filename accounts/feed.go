package main

import (
	"encoding/xml"
	"io"
	"time"
)

// The rss types give the RSS 2.0 shape. encoding/xml does the escaping.
type rss struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel channel  `xml:"channel"`
}

type channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Language    string `xml:"language"`
	Items       []item `xml:"item"`
}

type item struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	GUID        guid     `xml:"guid"`
	PubDate     string   `xml:"pubDate,omitempty"`
	Categories  []string `xml:"category"`
	Description string   `xml:"description"`
}

type guid struct {
	IsPermaLink string `xml:"isPermaLink,attr"`
	Value       string `xml:",chardata"`
}

// feedSize is the number of items in a private feed.
const feedSize = 20

// buildFeed makes the private feed for a user who muted the tags in muted.
// body returns the HTML body of a post.
func buildFeed(baseURL string, posts []Post, muted map[string]bool, body func(slug string) string) channel {
	ch := channel{
		Title:       "Recently Written",
		Link:        baseURL + "/",
		Description: "A site full of things which I have recently written.",
		Language:    "en-us",
	}
	for _, p := range posts {
		if len(ch.Items) == feedSize {
			break
		}
		if !wanted(p, muted) {
			continue
		}
		url := baseURL + "/" + p.Slug + ".html"
		ch.Items = append(ch.Items, item{
			Title:       p.Title,
			Link:        url,
			GUID:        guid{IsPermaLink: "true", Value: url},
			PubDate:     feedDate(p.Date),
			Categories:  p.Tags,
			Description: body(p.Slug),
		})
	}
	return ch
}

// feedDate gives the RFC 822 date that RSS requires, in the format that
// build.sh uses. A date that does not parse gives an empty string, and the
// item then has no pubDate.
func feedDate(d string) string {
	t, err := time.Parse("2006-01-02", d)
	if err != nil {
		return ""
	}
	return t.Format(time.RFC1123Z)
}

// writeFeed writes ch as an RSS 2.0 document. encoding/xml escapes the
// text and replaces each character that XML 1.0 does not permit.
func writeFeed(w io.Writer, ch channel) error {
	if _, err := io.WriteString(w, xml.Header); err != nil {
		return err
	}
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(rss{Version: "2.0", Channel: ch}); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}
