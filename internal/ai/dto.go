package ai

type GenerateTextRequest struct {
	Prompt       string `json:"prompt"`
	Stream       bool   `json:"stream"`
	OutputFormat string `json:"output_format"`
}

type GenerateTextResponse struct {
	Text string `json:"text"`
}

type StructuredTextResponse struct {
	Title     string   `json:"title"`
	Summary   string   `json:"summary"`
	KeyPoints []string `json:"key_points"`
}