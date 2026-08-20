package ai

type GenerateTextRequest struct {
	Prompt string `json:"prompt" binding:"required"`
}

type GenerateTextResponse struct {
	Text string `json:"text"`
}