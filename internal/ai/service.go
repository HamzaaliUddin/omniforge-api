package ai

import "context"

type Service struct {
	provider Provider
}

func NewService(provider Provider) *Service {
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