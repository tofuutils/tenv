/*
Copyright © 2024 Alexander Sharov <kvendingoldo@gmail.com>, Nikolai Mishin <sanduku.default@gmail.com>, Anastasiia Kozlova <anastasiia.kozlova245@gmail.com>
*/
package cmd

import (
	"fmt"
	"github.com/opentofuutils/tenv/pkg/consts/text"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Update environment to use tenv correctly",
	Long:  text.InitCmdLongText + text.SubCmdHelpText,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")
		//export PATH="${TOFUENV_ROOT}/bin:${PATH}";
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
