/*
Copyright © 2023 Dennis Koehler <dennis.koehler.it@gmail.com>
*/
package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/0x9ef/openai-go"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/textinput"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type chatGPTConfig struct {
	completionOptions *openai.CompletionOptions
	editOptions       *openai.EditOptions
}

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

			if len(args) == 0 {
				userCompletionInput := promptGetInput(instructionPromptContent)
				args = append(args, userCompletionInput)
			}

			apiKey := viper.GetString("application.apiKey")
			if apiKey != "" {
				chatGPT := openai.New(apiKey)

				config := getCompletionConfig()

				ctx := context.Background()
				config.completionOptions.Prompt = args
				coResponse, err := chatGPT.Completion(ctx, config.completionOptions)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				if _, err := json.MarshalIndent(coResponse, "", "  "); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				fmt.Printf("\n%s\n\n\n", coResponse.Choices[0].Text)

				for {
					isExit := promptGetConfirm(exitPromptContent)
					if isExit {
						os.Exit(1)
					}
					userEditInput := promptGetInput(instructionPromptContent)

					config.editOptions.Input = coResponse.Choices[0].Text
					config.editOptions.Instruction = userEditInput
					editResponse, err := chatGPT.Edit(ctx, config.editOptions)
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
						os.Exit(1)
					}
					fmt.Printf("\n%s\n\n\n", editResponse.Choices[0].Text)
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

func getCompletionConfig() chatGPTConfig {

	completionOptions := &openai.CompletionOptions{
		Model: openai.DefaultModel,
	}

	editOptions := &openai.EditOptions{
		Model: openai.DefaultModel,
	}

	dataModel := viper.GetString("application.dataModel")
	if dataModel != "" {
		completionOptions.Model = openai.Model(dataModel)
	}

	maxTokens := viper.GetInt("application.maxTokens")
	if maxTokens > 0 {
		completionOptions.MaxTokens = maxTokens
	}

	temperature := viper.GetFloat64("application.temperature")
	if temperature > 0 {
		editOptions.Temperature = float32(temperature)
		completionOptions.Temperature = float32(temperature)
	}

	maxCompletions := viper.GetInt("application.maxCompletions")
	if maxCompletions > 0 {
		editOptions.N = maxCompletions
		completionOptions.N = maxCompletions
	}

	sequencesStop := viper.GetStringSlice("application.sequencesStop")
	if len(sequencesStop) > 0 {
		completionOptions.Stop = sequencesStop
	}

	config := chatGPTConfig{
		editOptions:       editOptions,
		completionOptions: completionOptions,
	}

	return config
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
