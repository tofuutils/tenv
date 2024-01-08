/*
Copyright © 2024 Alexander Sharov <kvendingoldo@gmail.com>, Nikolai Mishin <sanduku.default@gmail.com>, Anastasiia Kozlova <anastasiia.kozlova245@gmail.com>
*/
package cmd

import (
	"bytes"
	"fmt"
	"github.com/opentofuutils/tenv/pkg/misc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os/exec"
)

// tfCmd represents the tf command
var tfCmd = &cobra.Command{
	Use:   "tf",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !misc.CheckToolInstalled("tfenv") {
			log.Error("tfenv is not installed. Please, execute 'tenv upgrade-deps' to use 'tenv tf' commands")
			return
		}

		toolExec := misc.GetPath("tfenv_exec")

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
	rootCmd.AddCommand(tfCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tfCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tfCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
