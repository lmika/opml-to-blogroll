package models

import "encoding/xml"

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

func (o OPML) FeedItems() []*Outline {
	var fic feedItemCollector
	o.addFeedItem(&fic)
	return fic.outline
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
