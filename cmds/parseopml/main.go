package main

import (
	"bytes"
	"encoding/xml"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/lmika/opml-to-blogroll/models"
	"github.com/lmika/opml-to-blogroll/renderer"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	opmlBodyFile, opmlHeader, err := r.FormFile("opml")
	if err != nil {
		http.Error(w, "cannot read OPML file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer opmlBodyFile.Close()

	// Unless the file extension is .xml, we cannot rely on the content-type. It may be one of
	// the other extension types.
	if !looksLikeValidOPMLFile(opmlHeader) {
		log.Printf("received unrecognised opml file: filename=%v, mimetype=%v", opmlHeader.Filename, opmlHeader.Header.Get("Content-type"))
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
	lambda.Start(httpadapter.New(http.HandlerFunc(Handler)).ProxyWithContext)
}

func looksLikeValidOPMLFile(h *multipart.FileHeader) bool {
	contentType, _, _ := strings.Cut(h.Header.Get("Content-type"), ";")
	contentType = strings.TrimSpace(strings.ToLower(contentType))

	if contentType == "application/xml" ||
		contentType == "text/xml" ||
		contentType == "text/x-opml" {
		return true
	}

	ext := strings.ToLower(filepath.Ext(h.Filename))
	if contentType == "application/octet-stream" && (ext == ".xml" || ext == ".opml") {
		return true
	}

	return false
}
