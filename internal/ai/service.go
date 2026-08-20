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
	text, err := s.provider.GenerateText(
		ctx,
		input.Prompt,
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
	return s.provider.StreamText(
		ctx,
		input.Prompt,
		onDelta,
	)
}
func (s *Service) GenerateStructured(
	ctx context.Context,
	input GenerateTextRequest,
) (*StructuredTextResponse, error) {
	result, err := s.provider.GenerateStructured(
		ctx,
		input.Prompt,
	)

	if err != nil {
		return nil, err
	}

	if err := ValidateStructuredOutput(result); err != nil {
		return nil, err
	}

	return result, nil
}