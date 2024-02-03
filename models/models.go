package models

import "encoding/xml"

type OPML struct {
	XMLName xml.Name   `xml:"opml"`
	Version string     `xml:"version,attr"`
	Head    Head       `xml:"head"`
	Body    []*Outline `xml:"body>outline"`
}

func (o OPML) FeedItems() []*Outline {
	outlines := make([]*Outline, 0)
	for _, o := range o.Body {
		if o.isFeedItem() {
			outlines = append(outlines, o)
		}
		outlines = append(outlines, o.feedItems()...)
	}
	return outlines
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

func (o Outline) feedItems() []*Outline {
	outlines := make([]*Outline, 0)
	for _, o := range o.Children {
		if o.isFeedItem() {
			outlines = append(outlines, o)
		}
		outlines = append(outlines, o.feedItems()...)
	}
	return outlines
}
