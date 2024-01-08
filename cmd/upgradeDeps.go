/*
Copyright © 2024 Alexander Sharov <kvendingoldo@gmail.com>, Nikolai Mishin <sanduku.default@gmail.com>, Anastasiia Kozlova <anastasiia.kozlova245@gmail.com>
*/
package cmd

import (
	"github.com/opentofuutils/tenv/pkg/consts/text"
	"github.com/opentofuutils/tenv/pkg/tool"
	"github.com/opentofuutils/tenv/pkg/utils/env"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//	"github.com/opentofuutils/tenv/pkg/github"

// upgradeDepsCmd represents the upgradeDeps command
var upgradeDepsCmd = &cobra.Command{
	Use:   "upgradeDeps",
	Short: "Upgrade tenv dependencies (tfenv and tofuenv)",
	Long:  text.UpgradeDepsCmdLongText + text.SubCmdHelpText,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Starting to upgrade tenv tools")

		rootDir := env.GetEnv(env.RootEnv, "")

		err := tool.PrepareTool("tfutils", "tfenv", rootDir)
		if err != nil {
			return
		}

		err = tool.PrepareTool("opentofuutils", "tofuenv", rootDir)
		if err != nil {
			return
		}

		log.Info("tenv tools upgraded successfully")
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
