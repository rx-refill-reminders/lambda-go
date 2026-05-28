package middleware

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type HandlerFunc func(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)
