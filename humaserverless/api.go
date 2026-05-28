package humaserverless

import (
	"context"
	"fmt"
	"maps"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

type ApiOpts struct {
	Name string

	Version string

	Servers []*huma.Server

	OAuth2 *OAuth2Opts
}

type OAuth2Opts struct {
	AuthDomain string
	Scopes     map[string]string

	AutoConfig *OAuth2AutoConfig
}

type OAuth2AutoConfig struct {
	ClientID string
	Scopes   []string
}

func NewServerless(opts ApiOpts) huma.API {
	if opts.Version == "" {
		opts.Version = "undefined"
	}

	router := chi.NewMux()

	config := huma.DefaultConfig(opts.Name, opts.Version)
	config.CreateHooks = nil

	api := humachi.New(router, config)

	api.OpenAPI().Servers = opts.Servers

	// Add OAuth 2.0 security scheme if Cognito domain is provided
	if opts.OAuth2 != nil {
		addOAuthSecurityScheme(api.OpenAPI(), opts)
	}

	return api
}

// addOAuthSecurityScheme adds OAuth 2.0 security scheme to the OpenAPI spec
func addOAuthSecurityScheme(spec *huma.OpenAPI, opts ApiOpts) {
	// Initialize Components if not already present
	if spec.Components == nil {
		spec.Components = &huma.Components{}
	}

	if spec.Components.SecuritySchemes == nil {
		spec.Components.SecuritySchemes = make(map[string]*huma.SecurityScheme)
	}

	if spec.Security == nil {
		spec.Security = []map[string][]string{}
	}

	if spec.Extensions == nil {
		spec.Extensions = map[string]any{}
	}

	schemeName := "OAuth2"

	authorizationURL := fmt.Sprintf("https://%s/authorize", opts.OAuth2.AuthDomain)
	tokenURL := fmt.Sprintf("https://%s/oauth2/token", opts.OAuth2.AuthDomain)

	scopeNames := slices.Collect(maps.Keys(opts.OAuth2.Scopes))

	spec.Components.SecuritySchemes[schemeName] = &huma.SecurityScheme{
		Type:        "oauth2",
		Name:        "Authorization",
		In:          "header",
		Scheme:      "OAuth",
		Description: "OAuth 2.0 Authorization Code flow with PKCE. Use AWS Cognito for authentication.",
		Flows: &huma.OAuthFlows{
			AuthorizationCode: &huma.OAuthFlow{
				AuthorizationURL: authorizationURL,
				TokenURL:         tokenURL,
				Scopes:           opts.OAuth2.Scopes,
			},
		},
	}

	spec.Security = append(spec.Security, map[string][]string{
		schemeName: scopeNames,
	})

	if opts.OAuth2.AutoConfig != nil {
		spec.Extensions["x-cli-config"] = huma.AutoConfig{
			Security: "OAuth2",
			Params: map[string]string{
				"client_id":     opts.OAuth2.AutoConfig.ClientID,
				"authorize_url": authorizationURL,
				"token_url":     tokenURL,
				"scopes":        strings.Join(opts.OAuth2.AutoConfig.Scopes, ","),
			},
		}
	}
}

// HttpHandler handles a single request using a pre-built API (caller must register routes).
func HttpHandler(
	ctx context.Context,
	api huma.API,
	event events.APIGatewayV2HTTPRequest,
) (events.APIGatewayV2HTTPResponse, error) {
	request, err := EventToRequest(event)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}
	response := httptest.NewRecorder()
	api.Adapter().ServeHTTP(response, request)
	return *ResponseToEvent(*response), nil
}
