package api

import (
	"chatgpt-api-proxy/pkg/httphelper"
	"chatgpt-api-proxy/pkg/middlerware"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"

	"github.com/gin-gonic/gin"
)

func InitCompletionRouter(r *gin.Engine) {
	api := r.Group("/api/openai")
	api.POST("/completions", middlerware.OpenAIUsage(), HandleCompletion)
}

type CompletionRequest struct {
	Model            string         `json:"model"`
	Prompt           any            `json:"prompt,omitempty"`
	Suffix           string         `json:"suffix,omitempty"`
	MaxTokens        int            `json:"max_tokens,omitempty"`
	Temperature      float32        `json:"temperature,omitempty"`
	TopP             float32        `json:"top_p,omitempty"`
	N                int            `json:"n,omitempty"`
	Stream           bool           `json:"stream,omitempty"`
	LogProbs         int            `json:"logprobs,omitempty"`
	Echo             bool           `json:"echo,omitempty"`
	Stop             []string       `json:"stop,omitempty"`
	PresencePenalty  float32        `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32        `json:"frequency_penalty,omitempty"`
	BestOf           int            `json:"best_of,omitempty"`
	LogitBias        map[string]int `json:"logit_bias,omitempty"`
	User             string         `json:"user,omitempty"`
}

// CompletionChoice represents one of possible completions.
type CompletionChoice struct {
	Text         string        `json:"text"`
	Index        int           `json:"index"`
	FinishReason string        `json:"finish_reason"`
	LogProbs     LogprobResult `json:"logprobs"`
}

// LogprobResult represents logprob result of Choice.
type LogprobResult struct {
	Tokens        []string             `json:"tokens"`
	TokenLogprobs []float32            `json:"token_logprobs"`
	TopLogprobs   []map[string]float32 `json:"top_logprobs"`
	TextOffset    []int                `json:"text_offset"`
}

// CompletionResponse represents a response structure for completion API.
type CompletionResponse struct {
	ID      string             `json:"id"`
	Object  string             `json:"object"`
	Created int64              `json:"created"`
	Model   string             `json:"model"`
	Choices []CompletionChoice `json:"choices"`
	Usage   Usage              `json:"usage"`
}

// Usage Represents the total token usage per request to OpenAI.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func HandleCompletion(c *gin.Context) {
	// parse request
	request := CompletionRequest{}

	if err := c.Bind(&request); err != nil {
		httphelper.WrapperError(c, httphelper.ErrInvalidRequestError)
		return
	}

	// call openai client
	client := openai.NewClient(getOpenAIAPIKey(c))
	completion, err := client.CreateCompletion(c, openai.CompletionRequest{
		Model:            request.Model,
		Prompt:           request.Prompt,
		MaxTokens:        request.MaxTokens,
		Temperature:      request.Temperature,
		TopP:             request.TopP,
		N:                request.N,
		Stream:           request.Stream,
		LogProbs:         request.LogProbs,
		Echo:             request.Echo,
		Stop:             request.Stop,
		PresencePenalty:  request.PresencePenalty,
		FrequencyPenalty: request.FrequencyPenalty,
		BestOf:           request.BestOf,
		LogitBias:        request.LogitBias,
		User:             request.User,
	})
	if err != nil {
		if errors.Is(openai.ErrCompletionStreamNotSupported, err) || errors.Is(openai.ErrCompletionUnsupportedModel, err) {
			httphelper.WrapperError(c, httphelper.NewAPIError(http.StatusBadRequest, err.Error()))
			return
		}
		httphelper.WrapperError(c, httphelper.NewAPIError(http.StatusInternalServerError, err.Error()))
		return
	}

	// return response
	httphelper.WrapperSuccess(c, completion)
}
