/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

//	"github.com/opentofuutils/tenv/pkg/utils"

// upgradeDepsCmd represents the upgradeDeps command
var upgradeDepsCmd = &cobra.Command{
	Use:   "upgradeDeps",
	Short: "Upgrade utils dependencies (tfenv and tofuenv)",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("upgradeDeps called")

		//destFolder := "./data"

		//tenv.CreateFolder(destFolder)

		// Download the latest release
		//if err := DownloadLatestRelease("tfutils", "tfenv", destFolder); err != nil {
		//	fmt.Println("Error:", err)
		//	return
		//}

		fmt.Println("Latest release downloaded successfully.")

	},
}

func init() {
	rootCmd.AddCommand(upgradeDepsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upgradeDepsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upgradeDepsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
