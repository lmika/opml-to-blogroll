package models

import (
	"encoding/xml"
	"github.com/lmika/gopkgs/fp/maps"
	"sort"
	"strings"
)

const (
	maxFeedItems  = 1000
	maxGroupDepth = 10
)

type OPML struct {
	XMLName xml.Name   `xml:"opml"`
	Version string     `xml:"version,attr"`
	Head    Head       `xml:"head"`
	Body    []*Outline `xml:"body>outline"`
}

func (o OPML) FeedItems() []FeedItem {
	var fic feedItemCollector
	o.addFeedItem(&fic)
	return fic.sortAndDedupe()
}

func (o OPML) addFeedItem(fic *feedItemCollector) {
	for _, o := range o.Body {
		if o.isFeedItem() {
			if !fic.add(o) {
				return
			}
		}
		o.feedItems(fic, 1)
	}
}

type Head struct {
	Title string `xml:"title"`
}

type Outline struct {
	Text    string `xml:"text,attr"`
	Title   string `xml:"title,attr"`
	Type    string `xml:"type,attr"`
	XMLUrl  string `xml:"xmlUrl,attr"`
	HTMLUrl string `xml:"htmlUrl,attr"`

	Children []*Outline `xml:"outline"`
}

type FeedItem struct {
	Title   string
	XMLUrl  string
	HTMLUrl string
}

func (o *Outline) isFeedItem() bool {
	return o.Type == "rss"
}

func (o Outline) feedItems(fic *feedItemCollector, depth int) {
	if depth >= maxGroupDepth {
		return
	}

	for _, o := range o.Children {
		if o.isFeedItem() {
			if !fic.add(o) {
				return
			}
		}
		o.feedItems(fic, depth+1)
	}
}

type feedItemCollector struct {
	outline []*Outline
}

func (fic *feedItemCollector) add(o *Outline) bool {
	if len(fic.outline) >= maxFeedItems {
		return false
	}
	fic.outline = append(fic.outline, o)
	return true
}

func (fic *feedItemCollector) sortAndDedupe() []FeedItem {
	urlsSeen := make(map[string]FeedItem)

	for _, o := range fic.outline {
		u := strings.ToLower(strings.TrimSpace(o.HTMLUrl))
		if _, ok := urlsSeen[u]; !ok {
			urlsSeen[u] = FeedItem{
				Title:   strings.TrimSpace(o.Title),
				HTMLUrl: strings.TrimSpace(o.HTMLUrl),
				XMLUrl:  strings.TrimSpace(o.XMLUrl),
			}
		}
	}

	feedItems := maps.Values(urlsSeen)
	sort.Slice(feedItems, func(i, j int) bool {
		return strings.ToLower(feedItems[i].Title) < strings.ToLower(feedItems[j].Title)
	})

	return feedItems
}
