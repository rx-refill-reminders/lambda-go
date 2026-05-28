package humaserverless

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// Convert the API Gateway event into a HTTP request
func EventToRequest(event events.APIGatewayV2HTTPRequest) (*http.Request, error) {
	httpMethod := event.RequestContext.HTTP.Method

	body := strings.NewReader(event.Body)

	requestURL := event.RawPath
	if event.RequestContext.Stage != "" {
		requestURL = strings.TrimPrefix(requestURL,
			fmt.Sprintf("/%s", event.RequestContext.Stage))
	}

	if event.RawQueryString != "" {
		requestURL = fmt.Sprintf("%s?%s", requestURL, event.RawQueryString)
	}

	request, err := http.NewRequest(httpMethod, requestURL, body)
	if err != nil {
		return nil, err
	}

	for headerName, headerVal := range event.Headers {
		request.Header.Set(headerName, headerVal)
	}

	return request, nil
}

// Serve the HTTP request to the API's adapter
// Convert the response to the API gateway response
func ResponseToEvent(response httptest.ResponseRecorder) *events.APIGatewayV2HTTPResponse {
	event := &events.APIGatewayV2HTTPResponse{
		StatusCode:        response.Code,
		Headers:           map[string]string{},
		MultiValueHeaders: map[string][]string{},
	}

	for headerName, headerValues := range response.Header() {
		if len(headerValues) > 0 {
			event.Headers[headerName] = headerValues[len(headerValues)-1]
			event.MultiValueHeaders[headerName] = headerValues
		}
	}

	contentType := event.Headers["Content-Type"]
	if isBinaryContentType(contentType) {
		event.Body = base64.StdEncoding.EncodeToString(response.Body.Bytes())
		event.IsBase64Encoded = true
	} else {
		event.Body = response.Body.String()
	}

	return event
}

func isBinaryContentType(contentType string) bool {
	return strings.HasPrefix(contentType, "image/") ||
		contentType == "application/octet-stream"
}
