package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	binaryName = "bbfw"
)

var rootCMD = &cobra.Command{
	Use:   binaryName,
	Short: "bamboo firewall cli",
	Long: fmt.Sprintf(`BAMBOO Firewall CLI
Description:
  The %s is used to manage global policy,
  to view and manage host endpoint, global network set configuration.`, binaryName),
}

func Execute() {
	rootCMD.AddCommand(createCMD)
	rootCMD.AddCommand(listCMD)
	rootCMD.AddCommand(getCMD)
	rootCMD.AddCommand(deleteCMD)
	rootCMD.AddCommand(validateCommand)
	rootCMD.AddCommand(versionCMD)

	rootCMD.AddCommand(&cobra.Command{
		Use:                   "completion",
		DisableFlagsInUseLine: true,
		Short:                 "Generate a completion script for bash or zsh shell",
		Example: `  # Gen completion for bash shell
  bbfw completion bash

  # Gen completion for zsh shell
  bbfw completion zsh`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				rootCMD.GenBashCompletionV2(os.Stdout, true)
			case "zsh":
				rootCMD.GenZshCompletion(os.Stdout)
			default:
				fmt.Fprintf(os.Stderr, "Unknown shell: %s\n", args[0])
			}
		},
	})

	if err := rootCMD.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
