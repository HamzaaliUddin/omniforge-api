package ai

import (
	"context"
	"encoding/json"
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
func (p *OpenAIProvider) StreamText(
	ctx context.Context,
	prompt string,
	onDelta func(string) error,
) error {
	stream := p.client.Responses.NewStreaming(
		ctx,
		responses.ResponseNewParams{
			Model: openai.ChatModelGPT5_2,
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String(prompt),
			},
		},
	)

	for stream.Next() {
		event := stream.Current()

		if event.Delta == "" {
			continue
		}

		if err := onDelta(event.Delta); err != nil {
			return err
		}
	}

	if err := stream.Err(); err != nil {
		return fmt.Errorf(
			"failed to stream text: %w",
			err,
		)
	}

	return nil
}
func (p *OpenAIProvider) GenerateStructured(
	ctx context.Context,
	prompt string,
) (*StructuredTextResponse, error) {
	schema, err := GenerateSchema[StructuredTextResponse]()
	if err != nil {
		return nil, err
	}

	result, err := p.client.Responses.New(
		ctx,
		responses.ResponseNewParams{
			Model: openai.ChatModelGPT5_2,
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String(prompt),
			},
			Text: responses.ResponseTextConfigParam{
				Format: responses.ResponseFormatTextConfigUnionParam{
					OfJSONSchema: &responses.ResponseFormatTextJSONSchemaConfigParam{
						Name:   "structured_text",
						Schema: schema,
						Strict: openai.Bool(true),
					},
				},
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to generate structured output: %w",
			err,
		)
	}

	var output StructuredTextResponse

	if err := json.Unmarshal(
		[]byte(result.OutputText()),
		&output,
	); err != nil {
		return nil, fmt.Errorf(
			"failed to decode structured output: %w",
			err,
		)
	}

	return &output, nil
}