package data

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Choice struct {
	Role         string  `json:"role"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type OpenApiResponse struct {
	Id      string   `json:"id"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

type Settings struct {
	ApiKey           string
	ApiURL           string
	MaxTokens        int
	Temperature      float32
	FrequencyPenalty float32
	PresencePenalty  float32
	Model            string
	SystemMessage    string
}

type ApiPostBody struct {
	Messages         []Message `json:"messages"`
	Model            string    `json:"model"`
	MaxTokens        int       `json:"max_tokens"`
	Temperature      float32   `json:"temperature"`
	FrequencyPenalty float32   `json:"frequency_penalty"`
	PresencePenalty  float32   `json:"presence_penalty"`
}
