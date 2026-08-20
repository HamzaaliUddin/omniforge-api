package ai

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

type OpenAIProvider struct {
	client openai.Client
}

func NewOpenAIProvider(
	apiKey string,
) *OpenAIProvider {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	return &OpenAIProvider{
		client: client,
	}
}

func (p *OpenAIProvider) GenerateText(
	ctx context.Context,
	prompt string,
) (string, error) {
	result, err := p.client.Responses.New(
		ctx,
		responses.ResponseNewParams{
			Model: openai.ChatModelGPT5_2,
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String(prompt),
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf(
			"failed to generate text: %w",
			err,
		)
	}

	return result.OutputText(), nil
}