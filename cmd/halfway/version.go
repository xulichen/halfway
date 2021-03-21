package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "打印版本",
	Long:  `等待补充`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}
