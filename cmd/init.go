/*
Copyright © 2024 Alexander Sharov <kvendingoldo@gmail.com>, Nikolai Mishin <sanduku.default@gmail.com>, Anastasiia Kozlova <anastasiia.kozlova245@gmail.com>
*/
package cmd

import (
	"github.com/opentofuutils/tenv/pkg/consts/text"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Update environment to use tenv correctly",
	Long:  text.InitCmdLongText + text.SubCmdHelpText,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Starting to init tenv")
		//export PATH="${TOFUENV_ROOT}/bin:${PATH}";
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
