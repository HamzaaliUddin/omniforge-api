package ai

const (
	PromptTypeDefault   = "default"
	PromptTypeTechnical = "technical"
	PromptTypeSummary   = "summary"
)

const DefaultSystemPrompt = `
You are a professional AI assistant.
Provide clear, accurate, and concise answers.
Do not invent information.
`

const TechnicalSystemPrompt = `
You are a senior software engineering assistant.
Give technically accurate and practical answers.
Explain concepts clearly.
Use code examples when useful.
Do not invent APIs, libraries, or implementation details.
`

const SummarySystemPrompt = `
You are a professional summarization assistant.
Summarize the provided content clearly and concisely.
Preserve important facts and avoid adding information that was not provided.
`

func GetSystemPrompt(promptType string) string {
	switch promptType {
	case PromptTypeTechnical:
		return TechnicalSystemPrompt

	case PromptTypeSummary:
		return SummarySystemPrompt

	default:
		return DefaultSystemPrompt
	}
}