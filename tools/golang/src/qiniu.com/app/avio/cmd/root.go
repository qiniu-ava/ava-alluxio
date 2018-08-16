package cmd

import (
	"fmt"
	"os"

	"qiniu.com/app/avio/util"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "avio",
	Short: "avio 是为 ava 用户开发的，便于使用 alluxio 缓存系统的管理工具",
	Long:  `avio 帮助 ava 用户更方便的使用 alluxio 缓存系统，提高预取数据、同步数据、查看数据状态等操作的速度和便利性。`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return util.NewAvioError(util.ARGUMENT_ERROR_ROOT, "请输入子命令")
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
