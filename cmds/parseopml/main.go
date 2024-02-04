package main

import (
	"bytes"
	"encoding/xml"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/lmika/opml-to-blogroll/models"
	"github.com/lmika/opml-to-blogroll/renderer"
	"io"
	"net/http"
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

	var bfr bytes.Buffer
	if err := renderer.Outlines(&bfr, feedItems); err != nil {
		http.Error(w, "cannot render template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "text/html; encoding=utf-8")
	w.WriteHeader(http.StatusOK)
	io.Copy(w, &bfr)
}

func main() {
	lambda.Start(httpadapter.New(http.HandlerFunc(handler)).ProxyWithContext)
}
