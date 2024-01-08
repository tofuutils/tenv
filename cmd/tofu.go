/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
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

// tofuCmd represents the tofu command
var tofuCmd = &cobra.Command{
	Use:   "tofu",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		if !misc.CheckToolInstalled("tofuenv") {
			log.Error("tofuenv is not installed. Please, execute 'tenv upgrade-deps' to use 'tenv tofu' commands")
			return
		}

		//fmt.Println("tofu called")

		tofuExec := misc.GetPath("tofuenv_exec")
		fmt.Println(tofuExec)
		//fmt.Println(tofuExec)

		exec := exec.Command(tofuExec, args...)
		var out bytes.Buffer
		var stderr bytes.Buffer
		exec.Stdout = &out
		exec.Stderr = &stderr
		err := exec.Run()

		if err != nil {
			fmt.Println(stderr.String())
			return
		}
		fmt.Println("Result: " + out.String())

		//out, _ := exec.Command(tofuExec, args...).Output()
		////if err != nil {
		////	fmt.Println(err)
		////}
		//fmt.Println(string(out))

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
