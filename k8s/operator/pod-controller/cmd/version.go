package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: 添加 git commit message
		log.Printf("version: v1.0")
	},
}
