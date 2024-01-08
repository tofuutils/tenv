/*
Copyright © 2024 Alexander Sharov <kvendingoldo@gmail.com>, Nikolai Mishin <sanduku.default@gmail.com>, Anastasiia Kozlova <anastasiia.kozlova245@gmail.com>
*/
package cmd

import (
	"github.com/opentofuutils/tenv/pkg/consts/text"
	"github.com/opentofuutils/tenv/pkg/misc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// uninstallDepsCmd represents the uninstallDeps command
var uninstallDepsCmd = &cobra.Command{
	Use:   "uninstallDeps",
	Short: "Uninstall tenv dependencies (tfenv and tofuenv)",
	Long:  text.UninstallDepsCmdLongText + text.SubCmdHelpText,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Starting to uninstall tenv tools")

		err := misc.DeleteFolder(misc.GetPath("bin_dir"))
		if err != nil {
			log.Error("Error removing dependencies directory:", err)
		}

		log.Info("tenv dependencies have been uninstalled successfully")
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
