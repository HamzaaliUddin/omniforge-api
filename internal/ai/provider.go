package ai

import "context"

type Provider interface {
	GenerateText(
		ctx context.Context,
		prompt string,
	) (string, error)

	StreamText(
		ctx context.Context,
		prompt string,
		onDelta func(string) error,
	) error

	GenerateStructured(
		ctx context.Context,
		prompt string,
	) (*StructuredTextResponse, error)
}