package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	t.Run("should support various file MIME types and file extensions", func(t *testing.T) {
		scenarios := []struct {
			mimeType string
			filename string
		}{
			// All valid formats
			{mimeType: "text/xml", filename: "feeds.xml"},
			{mimeType: "text/xml", filename: "feeds.opml"},
			{mimeType: "text/xml; encoding=utf-8", filename: "feeds.xml"},
			{mimeType: "text/xml; encoding=utf-8", filename: "feeds.opml"},
			{mimeType: "text/x-opml", filename: "feeds.xml"},
			{mimeType: "text/x-opml", filename: "feeds.opml"},
			{mimeType: "application/xml", filename: "feeds.xml"},
			{mimeType: "application/xml", filename: "feeds.opml"},
			{mimeType: "application/octet-stream", filename: "feeds.xml"},
			{mimeType: "application/octet-stream", filename: "feeds.opml"},
		}

		for _, scenario := range scenarios {
			t.Run(fmt.Sprintf("mime=%v, filename=%v", scenario.mimeType, scenario.filename), func(t *testing.T) {
				reqBody := &bytes.Buffer{}

				mp := multipart.NewWriter(reqBody)
				mpf, err := mp.CreatePart(textproto.MIMEHeader{
					"Content-type":        []string{scenario.mimeType},
					"Content-Disposition": []string{fmt.Sprintf(`form-data; name="opml"; filename="%s"`, scenario.filename)},
				})
				assert.NoError(t, err)

				io.WriteString(mpf, sampleOpml)
				mp.Close()

				r := httptest.NewRequest("POST", "/", reqBody)
				r.Header.Set("Content-type", mp.FormDataContentType())

				wh := httptest.NewRecorder()

				Handler(wh, r)

				assert.Equal(t, http.StatusOK, wh.Code)
			})
		}
	})
}

var sampleOpml = `<?xml version="1.0" encoding="UTF-8"?>
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
