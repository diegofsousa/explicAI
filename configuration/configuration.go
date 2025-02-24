package configuration

import (
	"context"
	"errors"
	"os"

	"github.com/diegofsousa/explicAI/internal/infrastructure/log"
	"github.com/spf13/viper"
)

var (
	openAiApiKey = os.Getenv("OPEN_AI_API_KEY")
	config       = viper.New()
)

func Init() *viper.Viper {
	defaultConfigs()
	config.AutomaticEnv()

	if openAiApiKey == "" {
		log.LogError(context.Background(), "OPEN_AI_API_KEY is not set",
			errors.New("open ai apiKey is required"))
	}

	return config
}

func defaultConfigs() {
	config.SetDefault("server.host", "0.0.0.0:8080")
	config.SetDefault("app.name", "explicAI")
	config.SetDefault("whisper.name", "whisper")
	config.SetDefault("whisper.url", "api.openai.com")
	config.SetDefault("whisper.host", "https://api.openai.com")
	config.SetDefault("whisper.timeout", 30000)
	config.SetDefault("whisper.model", "whisper-1")
	config.SetDefault("chatgpt.name", "chatgpt")
	config.SetDefault("chatgpt.url", "api.openai.com")
	config.SetDefault("chatgpt.host", "https://api.openai.com")
	config.SetDefault("chatgpt.timeout", 30000)
	config.SetDefault("chatgpt.model", "gpt-4o")
	config.SetDefault("database.url", "postgres://admin:admin@localhost:5432/explicai")
}
