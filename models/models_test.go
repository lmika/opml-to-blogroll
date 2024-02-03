package models_test

import (
	"encoding/xml"
	"github.com/lmika/opml-to-blogroll/models"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestModels(t *testing.T) {
	var outline models.OPML

	err := xml.NewDecoder(strings.NewReader(outlineXML)).Decode(&outline)
	assert.NoError(t, err)

	assert.Equal(t, "My Test", outline.Head.Title)
	assert.Equal(t, "Feed 1", outline.Body[0].Title)
	assert.Equal(t, "rss", outline.Body[0].Type)
	assert.Equal(t, "https://example.com/feed", outline.Body[0].XMLUrl)
	assert.Equal(t, "https://example.com/", outline.Body[0].HTMLUrl)

	assert.Equal(t, "Group", outline.Body[1].Title)

	assert.Equal(t, "Group Feed 1", outline.Body[1].Children[0].Text)
	assert.Equal(t, "rss", outline.Body[1].Children[0].Type)
	assert.Equal(t, "https://example.com/group/feed", outline.Body[1].Children[0].XMLUrl)
	assert.Equal(t, "https://example.com/group", outline.Body[1].Children[0].HTMLUrl)
}

func TestModels_FeedItems(t *testing.T) {
	var outline models.OPML

	err := xml.NewDecoder(strings.NewReader(outlineXML)).Decode(&outline)
	assert.NoError(t, err)

	feedItems := outline.FeedItems()
	assert.Len(t, feedItems, 2)

	assert.Equal(t, "Feed 1", feedItems[0].Title)
	assert.Equal(t, "rss", feedItems[0].Type)
	assert.Equal(t, "https://example.com/feed", feedItems[0].XMLUrl)
	assert.Equal(t, "https://example.com/", feedItems[0].HTMLUrl)

	assert.Equal(t, "Group Feed 1", feedItems[1].Text)
	assert.Equal(t, "rss", feedItems[1].Type)
	assert.Equal(t, "https://example.com/group/feed", feedItems[1].XMLUrl)
	assert.Equal(t, "https://example.com/group", feedItems[1].HTMLUrl)
}

var outlineXML = `<?xml version="1.0" encoding="UTF-8"?>
<opml version="1.0">
	<head>
		<title>My Test</title>
	</head>
	<body>
		<outline text="Feed 1" title="Feed 1" type="rss" xmlUrl="https://example.com/feed"
			htmlUrl="https://example.com/" />
		<outline text="Group" title="Group">
			<outline text="Group Feed 1" title="Feed 1" type="rss" xmlUrl="https://example.com/group/feed"
				htmlUrl="https://example.com/group" />
		</outline>
	</body>
</opml>`
