package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	version    string = "v0.0.1"
	projectDir string
	bashRC     string

	rootCmd = &cobra.Command{
		Use:   "halfway",
		Short: "等待补充",
		Long:  `等待补充`,
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
