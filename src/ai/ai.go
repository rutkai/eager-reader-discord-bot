package ai

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"os"
	"strings"
)

var client *openai.Client

func init() {
	client = openai.NewClient(getOpenAiToken())
}

func GetSummary(url string, quote string) (string, error) {
	if quote == "" || !strings.Contains(quote, "%s") {
		quote = "Give me a short summary of the content located under the URL %s. " +
			"Your response should be 3-6 bullet points and should be in the same language as the original content. " +
			"Start the text with stating this is a short summary using the same language as the original content. Do not mention the page URL."
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf(quote, url),
				},
			},
		},
	)

	if err != nil {
		log.Error().Err(err).Str("url", url).Str("quote", quote).Msg("OpenAI chat completion error.")
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func getOpenAiToken() string {
	token := os.Getenv("OPENAI_TOKEN")
	if token == "" {
		panic("OPENAI_TOKEN environment variable is missing")
	}
	return token
}
