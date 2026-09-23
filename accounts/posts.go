package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Post is one line of the post index that build.sh writes.
type Post struct {
	Date  string // YYYY-MM-DD, or empty
	Slug  string
	Tags  []string
	Title string
}

// validSlug accepts only the file names that build.sh makes. The check
// keeps a bad index line from reading a file outside postdata/.
var validSlug = regexp.MustCompile(`^[A-Za-z0-9_-][A-Za-z0-9._-]*$`)

// parseIndex reads postdata/index.tsv. Each line has four fields that
// tabs separate: date, slug, tags, title. The lines are newest first.
func parseIndex(r io.Reader) ([]Post, error) {
	var posts []Post
	sc := bufio.NewScanner(r)
	for n := 1; sc.Scan(); n++ {
		line := strings.TrimRight(sc.Text(), "\r")
		if line == "" {
			continue
		}
		f := strings.SplitN(line, "\t", 4)
		if len(f) != 4 || !validSlug.MatchString(f[1]) {
			return nil, fmt.Errorf("index line %d is not valid", n)
		}
		posts = append(posts, Post{
			Date:  f[0],
			Slug:  f[1],
			Tags:  splitTags(f[2]),
			Title: f[3],
		})
	}
	return posts, sc.Err()
}

// splitTags reads a front matter tags value such as "articles" or
// "lisp, unix".
func splitTags(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\t'
	})
}

func loadPosts(siteDir string) ([]Post, error) {
	f, err := os.Open(filepath.Join(siteDir, "postdata", "index.tsv"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parseIndex(f)
}

// loadBody returns the HTML body of a post. A missing body gives an
// empty string, because a feed item without a body is still useful.
func loadBody(siteDir, slug string) string {
	b, err := os.ReadFile(filepath.Join(siteDir, "postdata", slug+".html"))
	if err != nil {
		return ""
	}
	return string(b)
}

// allTags returns each tag that a post uses, in the order of first use.
func allTags(posts []Post) []string {
	seen := map[string]bool{}
	var tags []string
	for _, p := range posts {
		for _, t := range p.Tags {
			if !seen[t] {
				seen[t] = true
				tags = append(tags, t)
			}
		}
	}
	return tags
}

// wanted tells if a user who muted the tags in muted gets post p. A post
// without tags always goes out. A post with tags goes out if one or more
// of its tags is not muted.
func wanted(p Post, muted map[string]bool) bool {
	if len(p.Tags) == 0 {
		return true
	}
	for _, t := range p.Tags {
		if !muted[t] {
			return true
		}
	}
	return false
}
