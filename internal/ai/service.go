package ai

import "context"

type Service struct {
	provider Provider
}

func NewService(
	provider Provider,
) *Service {
	return &Service{
		provider: provider,
	}
}

func (s *Service) GenerateText(
	ctx context.Context,
	input GenerateTextRequest,
) (*GenerateTextResponse, error) {
	providerInput := ProviderInput{
		Prompt: input.Prompt,
		SystemPrompt: GetSystemPrompt(
			input.PromptType,
		),
	}

	text, err := s.provider.GenerateText(
		ctx,
		providerInput,
	)

	if err != nil {
		return nil, err
	}

	return &GenerateTextResponse{
		Text: text,
	}, nil
}

func (s *Service) StreamText(
	ctx context.Context,
	input GenerateTextRequest,
	onDelta func(string) error,
) error {
	providerInput := ProviderInput{
		Prompt: input.Prompt,
		SystemPrompt: GetSystemPrompt(
			input.PromptType,
		),
	}

	return s.provider.StreamText(
		ctx,
		providerInput,
		onDelta,
	)
}

func (s *Service) GenerateStructured(
	ctx context.Context,
	input GenerateTextRequest,
) (*StructuredTextResponse, error) {
	providerInput := ProviderInput{
		Prompt: input.Prompt,
		SystemPrompt: GetSystemPrompt(
			input.PromptType,
		),
	}

	result, err := s.provider.GenerateStructured(
		ctx,
		providerInput,
	)

	if err != nil {
		return nil, err
	}

	if err := ValidateStructuredOutput(result); err != nil {
		return nil, err
	}

	return result, nil
}