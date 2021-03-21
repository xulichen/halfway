package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "生成项目",
	Long:  `等待补充`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}
