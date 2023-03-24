/*
Copyright © 2023 Dennis Koehler <dennis.koehler.it@gmail.com>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/0x9ef/openai-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// promptCmd represents the prompt command
var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Your input for ChatGPT",
	Long: `This command is used to send your input to ChatGPT. For Example:
	openaicmd prompt <text>`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println(args)

		chatGPT := openai.New(viper.GetString("application.apikey"))

		CoResponse, err := chatGPT.Completion(context.Background(), &openai.CompletionOptions{
			// Choose model, you can see list of available models in models.go file
			Model: openai.DefaultModel,
			// Text to completion
			Prompt: []string{"Write a little bit of Wikipedia. What is that?"},
		})
		if err != nil {
			panic(err)
		}

		if b, err := json.MarshalIndent(CoResponse, "", "  "); err != nil {
			panic(err)
		} else {
			fmt.Println(string(b))
		}

		// Wikipedia is a free online encyclopedia, created and edited by volunteers.
		fmt.Println("What is the Wikipedia?", CoResponse.Choices[0].Text)
	},
}

func init() {
	rootCmd.AddCommand(promptCmd)
}
