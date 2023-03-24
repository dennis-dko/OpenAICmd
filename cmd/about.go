/*
Copyright © 2023 Dennis Koehler <dennis.koehler.it@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// aboutCmd represents the about command
var aboutCmd = &cobra.Command{
	Use:   "about",
	Short: "Print the information about OpenAICmd",
	Long:  `This command is used to get information about version, author and license.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\n%s %s, %s, %s\n\n", viper.GetString("general.name"), viper.GetString("general.version"), viper.GetString("general.author"), viper.GetString("general.license"))
	},
}

func init() {
	rootCmd.AddCommand(aboutCmd)
}
