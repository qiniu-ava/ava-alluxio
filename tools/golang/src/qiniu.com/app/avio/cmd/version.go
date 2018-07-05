package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	VERSION = "0.1.2" // declare-version
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "输出 avio 命令行工具的版本号",
	Long:  "输出 avio 命令行工具的版本号",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(VERSION)
	},
}
