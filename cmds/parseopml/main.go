package main

import (
	"encoding/xml"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lmika/gopkgs/fp/slices"
	"github.com/lmika/opml-to-blogroll/models"
	"net/http"
	"strings"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	opmlBody := request.Body

	var opml models.OPML
	if err := xml.NewDecoder(strings.NewReader(opmlBody)).Decode(&opml); err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-type": "text/plain; encoding=utf-8",
			},
			Body: err.Error(),
		}, nil
	}

	feedItems := opml.FeedItems()
	feedTitles := strings.Join(slices.Map(feedItems, func(fi *models.Outline) string {
		return "- " + fi.Title
	}), "\n")

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-type": "text/plain; encoding=utf-8",
		},
		Body: feedTitles,
	}, nil
}

func main() {
	lambda.Start(handler)
}
