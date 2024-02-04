package renderer

import (
	"bytes"
	"embed"
	"github.com/lmika/opml-to-blogroll/models"
	"html/template"
	"io"
	"sort"
	"strings"
)

//go:embed result.gohtml feeditems.gohtml
var resultTemplate embed.FS

func Outlines(w io.Writer, outlines []*models.Outline) error {
	tmpl, err := template.ParseFS(resultTemplate, "*.gohtml")
	if err != nil {
		return err
	}

	sort.Slice(outlines, func(i, j int) bool {
		return strings.ToLower(outlines[i].Title) < strings.ToLower(outlines[j].Title)
	})

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
