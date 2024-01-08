/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"github.com/opentofuutils/tenv/pkg/misc"
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
		//fmt.Println("tofu called")

		rootDir := misc.GetEnv(misc.RootEnv, "")
		binDir := fmt.Sprintf("%s/bin", rootDir)
		tofuExec := fmt.Sprintf("%s/tofu/bin/tofuenv", binDir)
		tofuExec = "/Users/asharov/go/src/github.com/opentofuutils/tenv/root/bin/tofuenv/tofuutils-tofuenv-e0bec88/bin/tofuenv"
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
