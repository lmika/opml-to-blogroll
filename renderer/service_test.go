package renderer_test

import (
	"bytes"
	"encoding/xml"
	"github.com/lmika/opml-to-blogroll/models"
	"github.com/lmika/opml-to-blogroll/renderer"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestOutlines(t *testing.T) {
	var outline models.OPML

	err := xml.NewDecoder(strings.NewReader(outlineXML)).Decode(&outline)
	assert.NoError(t, err)

	var bfr bytes.Buffer

	err = renderer.Outlines(&bfr, outline.FeedItems())
	assert.NoError(t, err)

	t.Log(bfr.String())
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
