/*
Copyright © 2023 Dennis Koehler <dennis.koehler.it@gmail.com>
*/
package cmd

import (
	"fmt"
	"sort"
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
		fmt.Printf("\n(%s)\t\n", viper.ConfigFileUsed())

		fmt.Println("\nName\t Value\t")

		settings := viper.AllSettings()
		for settingKey, settingValues := range settings {
			values, ok := settingValues.(map[string]interface{})
			if !ok || settingKey == "general" {
				continue
			}

			sortNames := make([]string, 0, len(values))
			for name := range values {
				sortNames = append(sortNames, name)
			}
			sort.Strings(sortNames)

			for _, name := range sortNames {
				runes := []rune(name)
				runes[0] = unicode.ToUpper(runes[0])
				fmt.Printf("\n%s\t %v\t\n", string(runes), values[name])
			}
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(settingsCmd)
}
