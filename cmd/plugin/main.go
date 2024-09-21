package main

import (
	"fmt"
	"os"

	"atlas-cli-plugin/internal/cli/restore"

	"github.com/spf13/cobra"
)

func main() {
	exampleCmd := &cobra.Command{
		Use:   "tools",
		Short: "MongoDB Tools for the Atlas CLI",
	}

	exampleCmd.AddCommand(
		restore.Builder(),
	)

	completionOption := &cobra.CompletionOptions{
		DisableDefaultCmd:   true,
		DisableNoDescFlag:   true,
		DisableDescriptions: true,
		HiddenDefaultCmd:    true,
	}
	rootCmd := &cobra.Command{
		DisableFlagParsing: true,
		DisableAutoGenTag:  true,
		DisableSuggestions: true,
		CompletionOptions:  *completionOption,
	}
	rootCmd.AddCommand(exampleCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
