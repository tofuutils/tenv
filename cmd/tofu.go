/*
Copyright © 2024 Alexander Sharov <kvendingoldo@gmail.com>, Nikolai Mishin <sanduku.default@gmail.com>, Anastasiia Kozlova <anastasiia.kozlova245@gmail.com>
*/
package cmd

import (
	"bytes"
	"fmt"
	"github.com/opentofuutils/tenv/pkg/consts/text"
	"github.com/opentofuutils/tenv/pkg/tool"
	"github.com/opentofuutils/tenv/pkg/utils/fs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os/exec"
)

// tofuCmd represents the tofu command
var tofuCmd = &cobra.Command{
	Use:   "tofu",
	Short: "Use tofuenv wrapper to manager OpenTofu versions",
	Long:  text.TofuCmdLongText + text.SubCmdHelpText,
	Run: func(cmd *cobra.Command, args []string) {
		if !tool.CheckToolInstalled("tofuenv") {
			log.Error("tofuenv is not installed. Please, execute 'tenv upgrade-deps' to use 'tenv tofu' commands")
			return
		}

		toolExec := fs.GetPath("tofuenv_exec")

		exec := exec.Command(toolExec, args...)
		var out bytes.Buffer
		var stderr bytes.Buffer
		exec.Stdout = &out
		exec.Stderr = &stderr
		err := exec.Run()

		if err != nil {
			fmt.Println(stderr.String())
			return
		}
		fmt.Println(out.String())
	},
}

func init() {
	rootCmd.AddCommand(tofuCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tofuCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tofuCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
