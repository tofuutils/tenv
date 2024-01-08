/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/opentofuutils/tenv/pkg/misc"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/spf13/cobra"
)

// uninstallDepsCmd represents the uninstallDeps command
var uninstallDepsCmd = &cobra.Command{
	Use:   "uninstallDeps",
	Short: "Uninstall tenv dependencies (tfenv and tofuenv)",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Starting to uninstall tenv tools")

		rootDir := misc.GetEnv(misc.RootEnv, "")
		binDir := fmt.Sprintf("%s/bin", rootDir)

		err := os.RemoveAll(binDir)
		if err != nil {
			log.Error("Error removing bin directory:", err)
		}

		log.Info("tenv tools uninstalled successfully")
	},
}

func init() {
	rootCmd.AddCommand(uninstallDepsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uninstallDepsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uninstallDepsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
