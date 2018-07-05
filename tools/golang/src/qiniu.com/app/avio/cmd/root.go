package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "avio",
	Short: "avio 是为 ava 用户开发的，便于使用 alluxio 缓存系统的管理工具",
	Long: `avio 帮助 ava 用户更方便的使用 alluxio 缓存系统，提高预取数据、
								同步数据、查看数据状态等操作的速度和便利性。`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
