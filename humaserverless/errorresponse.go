package humaserverless

import "github.com/aws/aws-lambda-go/events"

func ErrorResponse(code int, err error) events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: code,
		Body:       err.Error(),
	}
}
