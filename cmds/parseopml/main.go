package main

import (
	"bytes"
	"encoding/xml"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/lmika/gopkgs/fp/slices"
	"github.com/lmika/opml-to-blogroll/models"
	"io"
	"net/http"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	opmlBodyFile, opmlHeader, err := r.FormFile("opml")
	if err != nil {
		http.Error(w, "cannot read OPML file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer opmlBodyFile.Close()

	if opmlHeader.Header.Get("Content-type") != "text/xml" {
		http.Error(w, "expected XML file", http.StatusBadRequest)
		return
	}

	opmlBodyBytes, err := io.ReadAll(opmlBodyFile)
	if err != nil {
		http.Error(w, "cannot read OPML file", http.StatusInternalServerError)
		return
	}

	var opml models.OPML
	if err := xml.NewDecoder(bytes.NewReader(opmlBodyBytes)).Decode(&opml); err != nil {
		http.Error(w, "invalid XML file", http.StatusBadRequest)
		return
	}

	feedItems := opml.FeedItems()
	feedTitles := strings.Join(slices.Map(feedItems, func(fi *models.Outline) string {
		return "- " + fi.Title
	}), "\n")

	w.Header().Set("Content-type", "text/plain; encoding=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, feedTitles)
}

func main() {
	lambda.Start(httpadapter.New(http.HandlerFunc(handler)).ProxyWithContext)
}
