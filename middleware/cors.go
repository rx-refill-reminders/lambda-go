package middleware

import (
	"context"
	"log"
	"slices"

	"github.com/aws/aws-lambda-go/events"
)

type CORSOptions struct {
	ValidOrigins []string
}

var ValidOrigins = []string{
	"http://localhost:8081",
}

func WithCORS(opts CORSOptions) func(next HandlerFunc) HandlerFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
			requestOrigin := getOrigin(event.Headers)

			response, err := next(ctx, event)
			if err != nil {
				log.Printf("caught error: %v", err)
				return events.APIGatewayV2HTTPResponse{
					StatusCode: 500,
					Body:       err.Error(),
				}, nil
			}

			if response.Headers == nil {
				response.Headers = make(map[string]string)
			}

			response.Headers["Access-Control-Allow-Methods"] = "*"
			response.Headers["Access-Control-Allow-Headers"] = "*"
			response.Headers["Access-Control-Allow-Credentials"] = "true"

			if requestOrigin != "" && slices.Contains(opts.ValidOrigins, requestOrigin) {
				response.Headers["Access-Control-Allow-Origin"] = requestOrigin
			}

			return response, nil
		}
	}
}

func getOrigin(headers map[string]string) string {
	if origin, ok := headers["Origin"]; ok {
		return origin
	}
	if origin, ok := headers["origin"]; ok {
		return origin
	}
	return ""
}
