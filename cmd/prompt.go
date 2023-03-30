/*
Copyright © 2023 Dennis Koehler <dennis.koehler.it@gmail.com>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/textinput"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type promptContent struct {
	errorMsg    string
	label       string
	placeholder string
}

var (
	exitPromptContent = promptContent{
		label: "Do you want to exit?",
	}

	instructionPromptContent = promptContent{
		errorMsg:    "Please provide a text.",
		label:       "Type your instruction --> ",
		placeholder: "Write a little bit of Wikipedia. What is that?",
	}

	promptCmd = &cobra.Command{
		Use:   "prompt",
		Short: "Your input for ChatGPT",
		Long: `This command is used to send your input to ChatGPT. For Example:
	openaicmd prompt <text>`,
		Run: func(cmd *cobra.Command, args []string) {
			apiKey := viper.GetString("application.apiKey")
			if apiKey != "" {
				chatCompletion := getConfig()
				chatGPT := openai.NewClient(apiKey)
				messages := make([]openai.ChatCompletionMessage, 0)

				for {
					userCompletionInput := promptGetInput(instructionPromptContent)
					messages = append(messages, openai.ChatCompletionMessage{
						Role:    openai.ChatMessageRoleUser,
						Content: userCompletionInput,
					})

					chatCompletion.Messages = messages
					coResponse, err := chatGPT.CreateChatCompletion(
						context.Background(),
						chatCompletion,
					)
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
						continue
					}

					content := coResponse.Choices[0].Message.Content
					messages = append(messages, openai.ChatCompletionMessage{
						Role:    openai.ChatMessageRoleAssistant,
						Content: content,
					})
					fmt.Printf("\n%s\n\n\n", content)

					isExit := promptGetConfirm(exitPromptContent)
					if isExit {
						os.Exit(1)
					}
				}
			} else {
				fmt.Fprintln(os.Stderr, errors.New("API-Key not found"))
				os.Exit(1)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(promptCmd)
}

func getConfig() openai.ChatCompletionRequest {
	completionRequest := openai.ChatCompletionRequest{}

	dataModel := viper.GetString("application.dataModel")
	if dataModel != "" {
		completionRequest.Model = dataModel
	}

	maxTokens := viper.GetInt("application.maxTokens")
	if maxTokens > 0 {
		completionRequest.MaxTokens = maxTokens
	}

	temperature := viper.GetFloat64("application.temperature")
	if temperature > 0 {
		completionRequest.Temperature = float32(temperature)
	}

	maxCompletions := viper.GetInt("application.maxCompletions")
	if maxCompletions > 0 {
		completionRequest.N = maxCompletions
	}

	sequencesStop := viper.GetStringSlice("application.sequencesStop")
	if len(sequencesStop) > 0 {
		completionRequest.Stop = sequencesStop
	}

	return completionRequest
}

func promptGetConfirm(pc promptContent) bool {
	confirmInput := confirmation.New(pc.label, confirmation.Undecided)

	confirmResult, err := confirmInput.RunPrompt()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return confirmResult
}

func promptGetInput(pc promptContent) string {
	textInput := textinput.New(pc.label)
	textInput.Placeholder = pc.placeholder
	textInput.Validate = func(input string) error {
		if len(input) <= 0 {
			return errors.New(pc.errorMsg)
		}
		return nil
	}

	textResult, err := textInput.RunPrompt()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return textResult
}
