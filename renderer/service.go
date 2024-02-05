package renderer

import (
	"bytes"
	"embed"
	"github.com/lmika/opml-to-blogroll/models"
	"html/template"
	"io"
)

//go:embed result.gohtml feeditems.gohtml
var resultTemplate embed.FS

func Outlines(w io.Writer, outlines []models.FeedItem) error {
	tmpl, err := template.ParseFS(resultTemplate, "*.gohtml")
	if err != nil {
		return err
	}

	var feedItems bytes.Buffer
	if err := tmpl.ExecuteTemplate(&feedItems, "feeditems.gohtml", map[string]any{
		"FeedItems": outlines,
	}); err != nil {
		return err
	}

	if err := tmpl.ExecuteTemplate(w, "result.gohtml", map[string]any{
		"Output": feedItems.String(),
	}); err != nil {
		return err
	}

	return nil
}
