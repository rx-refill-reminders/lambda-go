package middleware

func DefaultChain(next HandlerFunc) HandlerFunc {
	withCORS := WithCORS(CORSOptions{
		ValidOrigins: []string{
			"http://localhost:8081",
		},
	})

	return WithLogger(
		withCORS(
			next,
		),
	)
}
