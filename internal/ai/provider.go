package ai

import "context"

type ProviderInput struct {
	Prompt       string
	SystemPrompt string
}

type Provider interface {
	GenerateText(
		ctx context.Context,
		input ProviderInput,
	) (string, error)

	StreamText(
		ctx context.Context,
		input ProviderInput,
		onDelta func(string) error,
	) error

	GenerateStructured(
		ctx context.Context,
		input ProviderInput,
	) (*StructuredTextResponse, error)
}