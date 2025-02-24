package chatgpt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/diegofsousa/explicAI/internal/gateway/summarize"
	"github.com/diegofsousa/explicAI/internal/infrastructure/clients"
)

const (
	basePath                   = "/v1/chat/completions"
	systemPrompt               = "Você é um sistema que recebe um texto transcrito de um áudio e organiza, separando em parágrafos e corrigindo possíveis erros de concordância"
	resumeUserPrompt           = "Preciso de um objeto com título sugerido, descrição sugerida, resumo breve e médio sobre a seguinte transcrição:"
	fullTextOrganizeUserPrompt = "Retorne apenas o texto normalizado para a seguinte transcrição: "
	functionCallName           = "resume"
)

type (
	ChatgptFunctionCallRequest struct {
		Model        string       `json:"model"`
		Messages     []Message    `json:"messages"`
		Functions    []Function   `json:"functions,omitempty"`
		FunctionCall FunctionCall `json:"function_call,omitempty"`
	}

	ChatgptSimpleRequest struct {
		Model    string    `json:"model"`
		Messages []Message `json:"messages"`
	}

	FunctionCall struct {
		Name string `json:"name"`
	}

	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	Function struct {
		Name       string         `json:"name"`
		Parameters FunctionParams `json:"parameters"`
	}

	FunctionParams struct {
		Type       string         `json:"type"`
		Properties FunctionFields `json:"properties"`
	}

	FunctionFields struct {
		Title        FieldSpec `json:"title"`
		Description  FieldSpec `json:"description"`
		BriefResume  FieldSpec `json:"briefResume"`
		MediumResume FieldSpec `json:"mediumResume"`
	}

	FieldSpec struct {
		Type        string `json:"type"`
		Description string `json:"description"`
	}

	ChatResumeCompletionResponse struct {
		Choices []struct {
			Message struct {
				FunctionCall struct {
					Arguments string `json:"arguments"`
				} `json:"function_call"`
			} `json:"message"`
		} `json:"choices"`
	}

	ChatFullTextCompletionResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
)

type Client struct {
	HttpClient  *clients.BaseHTTP
	ApiKey      string
	ServiceName string
	Model       string
}

func NewClient(serviceName, URL, apiKey, model string, timeout int64) *Client {
	return &Client{
		ServiceName: serviceName,
		HttpClient:  clients.NewHttpClient(URL, timeout),
		ApiKey:      apiKey,
		Model:       model,
	}
}

func (c *Client) Resume(ctx context.Context, transcription string) (
	*summarize.ResumeOutput, error,
) {
	clients.Mutex.Lock()
	req := c.HttpClient.Client.
		SetHeader("Authorization", "Bearer "+c.ApiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(c.buildResumeRequest(transcription))

	res, err := req.Post(basePath)
	clients.Mutex.Unlock()

	if res.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("error on chatgpt resume request: response=%s | status=%s", res.Body(), res.Status())
	}

	if err != nil {
		return nil, fmt.Errorf("error on chatgpt resume request: error=%s", err.Error())
	}

	var chatResponse ChatResumeCompletionResponse
	if err = json.Unmarshal(res.Body(), &chatResponse); err != nil {
		return nil, fmt.Errorf("error on chatgpt resume request: error=%s", err.Error())
	}

	choices := chatResponse.Choices
	if len(choices) <= 0 {
		return nil, fmt.Errorf("error on chatgpt resume request: no choices in response")
	}

	var response summarize.ResumeOutput
	if err = json.Unmarshal([]byte(choices[0].Message.FunctionCall.Arguments), &response); err != nil {
		return nil, fmt.Errorf("error on chatgpt resume request: error=%s", err.Error())
	}

	return &response, nil
}

func (c *Client) FullTextOrganize(ctx context.Context, transcription string) (*string, error) {
	clients.Mutex.Lock()

	req := c.HttpClient.Client.
		SetHeader("Authorization", "Bearer "+c.ApiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(c.buildFullTextOrganizeRequest(transcription))

	res, err := req.Post(basePath)
	clients.Mutex.Unlock()

	if res.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("error on chatgpt full text organize request: response=%s | status%s", res.Body(), res.Status())
	}

	if err != nil {
		return nil, fmt.Errorf("error on chatgpt full text organize request: error=%s", err.Error())
	}

	var response ChatFullTextCompletionResponse
	if err = json.Unmarshal(res.Body(), &response); err != nil {
		return nil, fmt.Errorf("error on chatgpt full text organize request: error=%s", err.Error())
	}

	if len(response.Choices) <= 0 || response.Choices[0].Message.Content == "" {
		return nil, fmt.Errorf("error on chatgpt full text organize request: empty response")
	}

	responseText := response.Choices[0].Message.Content
	return &responseText, nil
}

func (c *Client) buildResumeRequest(transcription string) ChatgptFunctionCallRequest {
	fullResumeUserPrompt := resumeUserPrompt + "\n" + transcription

	return ChatgptFunctionCallRequest{
		Model: c.Model,
		Messages: []Message{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: fullResumeUserPrompt,
			},
		},
		Functions: buildFunctionCallRequest(),
		FunctionCall: FunctionCall{
			Name: functionCallName,
		},
	}
}

func buildFunctionCallRequest() []Function {
	return []Function{
		{
			Name: functionCallName,
			Parameters: FunctionParams{
				Type: "object",
				Properties: FunctionFields{
					Title: FieldSpec{
						Type:        "string",
						Description: "Título de até 60 caracteres",
					},
					Description: FieldSpec{
						Type:        "string",
						Description: "Descrição de até 300 caracteres",
					},
					BriefResume: FieldSpec{
						Type:        "string",
						Description: "Resumo breve de até 5 linhas",
					},
					MediumResume: FieldSpec{
						Type:        "string",
						Description: "Resumo médio de até 15 linhas",
					},
				},
			},
		},
	}
}

func (c *Client) buildFullTextOrganizeRequest(transcription string) ChatgptSimpleRequest {
	fullTextOrganizePrompt := fullTextOrganizeUserPrompt + "\n" + transcription

	return ChatgptSimpleRequest{
		Model: c.Model,
		Messages: []Message{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: fullTextOrganizePrompt,
			},
		},
	}
}
