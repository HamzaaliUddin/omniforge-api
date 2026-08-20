package ai

const (
	MessageInvalidRequest          = "invalid request"
	MessagePromptRequired          = "prompt is required"
	MessagePromptTooLong           = "prompt is too long"
	MessageInvalidOutputFormat     = "output format must be text or structured"
	MessageInvalidPromptType       = "prompt type must be default, technical, or summary"
	MessageStructuredStreamInvalid = "structured output cannot be streamed"
	MessageGenerateFailed          = "failed to generate text"
	MessageGenerateSuccess         = "text generated successfully"
)