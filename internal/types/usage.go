package types

type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

func (u *Usage) Total() int {
	return u.InputTokens + u.OutputTokens
}
