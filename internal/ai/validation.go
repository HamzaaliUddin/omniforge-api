package ai

import (
	"errors"
	"strings"
)

const MaxPromptLength = 10000

var (
	ErrPromptRequired              = errors.New("prompt is required")
	ErrPromptTooLong               = errors.New("prompt is too long")
	ErrInvalidOutputFormat         = errors.New("invalid output format")
	ErrStructuredStreamingNotValid = errors.New("structured output cannot be streamed")

	ErrStructuredTitleRequired = errors.New(
		"structured output title is required",
	)

	ErrStructuredSummaryRequired = errors.New(
		"structured output summary is required",
	)

	ErrStructuredKeyPointsRequired = errors.New(
		"structured output key points are required",
	)
)

func ValidateGenerateTextRequest(
	input *GenerateTextRequest,
) error {
	input.Prompt = strings.TrimSpace(input.Prompt)

	if input.Prompt == "" {
		return ErrPromptRequired
	}

	if len(input.Prompt) > MaxPromptLength {
		return ErrPromptTooLong
	}

	if input.OutputFormat == "" {
		input.OutputFormat = OutputFormatText
	}

	if input.OutputFormat != OutputFormatText &&
		input.OutputFormat != OutputFormatStructured {
		return ErrInvalidOutputFormat
	}

	if input.Stream &&
		input.OutputFormat == OutputFormatStructured {
		return ErrStructuredStreamingNotValid
	}

	return nil
}

func ValidateStructuredOutput(
	output *StructuredTextResponse,
) error {
	output.Title = strings.TrimSpace(output.Title)
	output.Summary = strings.TrimSpace(output.Summary)

	if output.Title == "" {
		return ErrStructuredTitleRequired
	}

	if output.Summary == "" {
		return ErrStructuredSummaryRequired
	}

	if len(output.KeyPoints) == 0 {
		return ErrStructuredKeyPointsRequired
	}

	validKeyPoints := make(
		[]string,
		0,
		len(output.KeyPoints),
	)

	for _, point := range output.KeyPoints {
		point = strings.TrimSpace(point)

		if point != "" {
			validKeyPoints = append(
				validKeyPoints,
				point,
			)
		}
	}

	if len(validKeyPoints) == 0 {
		return ErrStructuredKeyPointsRequired
	}

	output.KeyPoints = validKeyPoints

	return nil
}