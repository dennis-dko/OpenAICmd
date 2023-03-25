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
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type chatGPTConfig struct {
	completionOptions *openai.CompletionOptions
	editOptions       *openai.EditOptions
}

type promptContent struct {
	errorMsg string
	label    string
}

var (
	exitPromptContent = promptContent{
		errorMsg: "Please provide an option.",
		label:    "Do you want to exit?",
	}

	instructionPromptContent = promptContent{
		errorMsg: "Please provide a text.",
		label:    "Type your instruction --> ",
	}

	promptCmd = &cobra.Command{
		Use:   "prompt",
		Short: "Your input for ChatGPT",
		Long: `This command is used to send your input to ChatGPT. For Example:
	openaicmd prompt <text>`,
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) == 0 {
				userCompletionInput := promptGetInput(instructionPromptContent, false, true, false)
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
				fmt.Printf("\nChatGPT answer: \n\n%s\n\n", coResponse.Choices[0].Text)

				for {
					promptGetInput(exitPromptContent, true, true, true)
					userEditInput := promptGetInput(instructionPromptContent, false, true, true)

					config.editOptions.Input = coResponse.Choices[0].Text
					config.editOptions.Instruction = userEditInput
					editResponse, err := chatGPT.Edit(ctx, config.editOptions)
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
						os.Exit(1)
					}
					fmt.Printf("\nChatGPT answer: \n\n%s\n\n", editResponse.Choices[0].Text)
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

func promptGetInput(pc promptContent, isConfirm bool, hideEntered bool, isValidate bool) string {
	validate := func(input string) error {
		if isValidate && len(input) <= 0 {
			return errors.New(pc.errorMsg)
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:       pc.label,
		Templates:   templates,
		Validate:    validate,
		IsConfirm:   isConfirm,
		HideEntered: hideEntered,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("\nYour Instruction: \n\n%s\n\n", result)

	return result
}
