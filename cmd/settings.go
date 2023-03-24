/*
Copyright © 2023 Dennis Koehler <dennis.koehler.it@gmail.com>
*/
package cmd

import (
	"fmt"
	"unicode"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// settingsCmd represents the settings command
var settingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Print the current settings of OpenAICmd",
	Long:  `This command is used to get the current Settings of OpenAICmd.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\nConfig: %s", viper.ConfigFileUsed())

		settings := viper.AllSettings()
		for settingKey, settingValues := range settings {
			values, ok := settingValues.(map[string]interface{})
			if !ok || settingKey == "general" {
				continue
			}

			for name, value := range values {
				runes := []rune(name)
				runes[0] = unicode.ToUpper(runes[0])
				fmt.Printf("\n\n%s: %s\n\n", string(runes), value)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(settingsCmd)
}
