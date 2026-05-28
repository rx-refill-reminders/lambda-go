package logs

import (
	"context"
	"maps"
)

type ctxKeyAnnotations struct{}

func GetAnnotations(ctx context.Context) map[string]any {
	// Check the context for existing annotations
	annotations, ok := ctx.Value(ctxKeyAnnotations{}).(map[string]any)

	// If there are no existing annotations, return an empty map
	if !ok {
		return map[string]any{}
	}

	// Otherwise, return the existing annotations
	return annotations
}

func WithAnnotations(
	ctx context.Context,
	annotations map[string]any,
) context.Context {
	// Fetch the existing annotations from the context
	mergedAnnotations := GetAnnotations(ctx)

	// Copy the new annotations into the existing ones
	maps.Copy(mergedAnnotations, annotations)

	// Return a new context, with the new annotations
	return context.WithValue(ctx, ctxKeyAnnotations{}, mergedAnnotations)
}

func WithAnnotation(
	ctx context.Context,
	key string,
	value any,
) context.Context {
	return WithAnnotations(ctx, map[string]any{
		key: value,
	})
}
