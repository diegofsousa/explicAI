package chatgpt

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/diegofsousa/explicAI/internal/infrastructure/clients"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var (
	config = viper.New()

	//go:embed embed/chatgpt-resume-response.json
	chatgptResumeResponse string

	//go:embed embed/chatgpt-resume-response-fail-marshall.json
	chatgptResumeResponseFailMarshall string

	//go:embed embed/chatgpt-resume-response-fail-marshall-arguments.json
	chatgptResumeResponseFailMarshallArguments string

	//go:embed embed/chatgpt-resume-response-no-choices.json
	chatgptResumeResponseNoChoices string

	//go:embed embed/chatgpt-fulltext-response.json
	chatgptFulltextResponse string

	//go:embed embed/chatgpt-fulltext-response-empty.json
	chatgptFulltextResponseEmpty string
)

type (
	ChatgptClientTestSuite struct {
		suite.Suite
		ctx context.Context

		chatgptClient Client
	}
)

func TestChatgptClient(t *testing.T) {
	suite.Run(t, new(ChatgptClientTestSuite))
}

func (s *ChatgptClientTestSuite) SetupTest() {
	s.ctx = context.Background()
	config.AddConfigPath("embed")
	config.SetConfigName("client_config")

	if err := config.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed conf file read: %w", err))
	}

	s.chatgptClient = getClientConfig(*config.Sub("chatgpt"))
}

func (s *ChatgptClientTestSuite) TearDownTest() {
	mock.AssertExpectationsForObjects(s.T())
}

func (s *ChatgptClientTestSuite) TestChatgptResume() {
	s.Run("successful request/response", func() {
		var response ChatResumeCompletionResponse
		json.Unmarshal([]byte(chatgptResumeResponse), &response)

		httpServerMockParams := clients.HttpServerMockParams{
			ExpectedPath:   basePath,
			ExpectedMethod: http.MethodPost,
			ResponseStatus: http.StatusOK,
			ResponseObject: response,
		}

		server := clients.StartMockServer(httpServerMockParams,
			config.Sub("chatgpt").GetString("host"),
		)

		defer server.Close()

		result, err := s.chatgptClient.Resume(s.ctx, "xpto")
		s.NoError(err)
		s.Equal("title test", result.Title)
		s.Equal("description test.", result.Description)
		s.Equal("brief test", result.BriefResume)
		s.Equal("medium test", result.MediumResume)
	})

	s.Run("fail response with status code error", func() {

		httpServerMockParams := clients.HttpServerMockParams{
			ExpectedPath:   basePath,
			ExpectedMethod: http.MethodPost,
			ResponseStatus: http.StatusInternalServerError,
			ResponseObject: nil,
		}

		server := clients.StartMockServer(httpServerMockParams,
			config.Sub("chatgpt").GetString("host"),
		)

		defer server.Close()

		_, err := s.chatgptClient.Resume(s.ctx, "xpto")
		s.Error(err)
		s.EqualError(err, "error on chatgpt resume request: response= | status=500 Internal Server Error")
	})

	s.Run("fail response with unmarshall error", func() {
		var response []map[string]string
		json.Unmarshal([]byte(chatgptResumeResponseFailMarshall), &response)

		httpServerMockParams := clients.HttpServerMockParams{
			ExpectedPath:   basePath,
			ExpectedMethod: http.MethodPost,
			ResponseStatus: http.StatusOK,
			ResponseObject: response,
		}

		server := clients.StartMockServer(httpServerMockParams,
			config.Sub("chatgpt").GetString("host"),
		)

		defer server.Close()

		_, err := s.chatgptClient.Resume(s.ctx, "xpto")
		s.Error(err)
		s.EqualError(err, "error on chatgpt resume request: error=json: cannot unmarshal array into Go value of type chatgpt.ChatResumeCompletionResponse")
	})

	s.Run("fail response with no choices error", func() {
		var response ChatResumeCompletionResponse
		json.Unmarshal([]byte(chatgptResumeResponseNoChoices), &response)

		httpServerMockParams := clients.HttpServerMockParams{
			ExpectedPath:   basePath,
			ExpectedMethod: http.MethodPost,
			ResponseStatus: http.StatusOK,
			ResponseObject: response,
		}

		server := clients.StartMockServer(httpServerMockParams,
			config.Sub("chatgpt").GetString("host"),
		)

		defer server.Close()

		_, err := s.chatgptClient.Resume(s.ctx, "xpto")
		s.Error(err)
		s.EqualError(err, "error on chatgpt resume request: no choices in response")
	})

	s.Run("fail response with argument unmarshall error", func() {
		var response ChatResumeCompletionResponse
		json.Unmarshal([]byte(chatgptResumeResponseFailMarshallArguments), &response)

		httpServerMockParams := clients.HttpServerMockParams{
			ExpectedPath:   basePath,
			ExpectedMethod: http.MethodPost,
			ResponseStatus: http.StatusOK,
			ResponseObject: response,
		}

		server := clients.StartMockServer(httpServerMockParams,
			config.Sub("chatgpt").GetString("host"),
		)

		defer server.Close()

		_, err := s.chatgptClient.Resume(s.ctx, "xpto")
		s.Error(err)
		s.EqualError(err, "error on chatgpt resume request: error=invalid character 'x' looking for beginning of value")
	})
}

func (s *ChatgptClientTestSuite) TestChatgptFulltext() {
	s.Run("successful request/response", func() {
		var response ChatFullTextCompletionResponse
		json.Unmarshal([]byte(chatgptFulltextResponse), &response)

		httpServerMockParams := clients.HttpServerMockParams{
			ExpectedPath:   basePath,
			ExpectedMethod: http.MethodPost,
			ResponseStatus: http.StatusOK,
			ResponseObject: response,
		}

		server := clients.StartMockServer(httpServerMockParams,
			config.Sub("chatgpt").GetString("host"),
		)

		defer server.Close()

		result, err := s.chatgptClient.FullTextOrganize(s.ctx, "xpto")
		s.NoError(err)
		s.Equal("text", *result)
	})

	s.Run("fail response with unmarshall error", func() {
		var response []map[string]string
		json.Unmarshal([]byte(chatgptResumeResponseFailMarshall), &response)

		httpServerMockParams := clients.HttpServerMockParams{
			ExpectedPath:   basePath,
			ExpectedMethod: http.MethodPost,
			ResponseStatus: http.StatusOK,
			ResponseObject: response,
		}

		server := clients.StartMockServer(httpServerMockParams,
			config.Sub("chatgpt").GetString("host"),
		)

		defer server.Close()

		_, err := s.chatgptClient.FullTextOrganize(s.ctx, "xpto")
		s.Error(err)
		s.EqualError(err, "error on chatgpt full text organize request: error=json: cannot unmarshal array into Go value of type chatgpt.ChatFullTextCompletionResponse")
	})

	s.Run("fail response with empty response", func() {
		var response ChatFullTextCompletionResponse
		json.Unmarshal([]byte(chatgptFulltextResponseEmpty), &response)

		httpServerMockParams := clients.HttpServerMockParams{
			ExpectedPath:   basePath,
			ExpectedMethod: http.MethodPost,
			ResponseStatus: http.StatusOK,
			ResponseObject: response,
		}

		server := clients.StartMockServer(httpServerMockParams,
			config.Sub("chatgpt").GetString("host"),
		)

		defer server.Close()

		_, err := s.chatgptClient.FullTextOrganize(s.ctx, "xpto")
		s.Error(err)
		s.EqualError(err, "error on chatgpt full text organize request: empty response")
	})
}

func getClientConfig(viper viper.Viper) Client {
	client := NewClient(
		viper.GetString("name"),
		viper.GetString("url"),
		viper.GetString("apiKey"),
		viper.GetString("model"),
		viper.GetInt64("timeout"),
	)

	return *client
}
