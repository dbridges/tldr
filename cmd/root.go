package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tldr",
	Short: "tldr is a simple tldr client that implements offline caching",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
	},
}

func Execute() {
	addDefaultCmd()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func addDefaultCmd() {
	cmds := map[string]bool{}
	cmds["help"] = true
	for _, cmd := range rootCmd.Commands() {
		cmds[cmd.Name()] = true
	}
	if len(os.Args) > 1 && !cmds[os.Args[1]] {
		args := append([]string{"view"}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}
}
