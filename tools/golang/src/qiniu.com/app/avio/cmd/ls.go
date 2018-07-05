package cmd

import (
	"fmt"
	"os"

	"qiniu.com/app/avio/util"

	"github.com/spf13/cobra"
)

type lsCMD struct {
}

type lsArguments struct {
	Full      bool
	Humanable bool
	Recursive bool
	Count     bool
	Args      []string
}

func (l *lsArguments) validate() error {
	return nil
}

func (l *lsArguments) String() string {
	return fmt.Sprintf("{Full: %t, Humanable: %t, Recursive: %t, Count: %t}", l.Full, l.Humanable, l.Recursive, l.Count)
}

var mLsArguments = &lsArguments{}

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.Flags().BoolVarP(&mLsArguments.Full, "full", "f", false, "当指定的路径下文件或子目录数量超过 5000 时，是否全量拉取，默认值为 true")
	lsCmd.Flags().BoolVarP(&mLsArguments.Humanable, "humanable", "H", false, "以易读的方式显示文件大小")
	lsCmd.Flags().BoolVarP(&mLsArguments.Recursive, "recursive", "r", false, "递归的列出目录下子目录中的文件")
	lsCmd.Flags().BoolVarP(&mLsArguments.Count, "count", "c", false, "只显示目录下的文件或子目录数")
	lsCmd.SetUsageTemplate(`Usage:{{if .Runnable}}
  {{.UseLine}} [path] {{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`)
}

var lsAliases = []string{"l", "list"}
var lsCmd = &cobra.Command{
	Use:                "ls",
	Short:              "列出指定目录或者文件的信息",
	Long:               `列出指定目录或者文件的信息`,
	DisableSuggestions: false,
	Aliases:            lsAliases,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		mLsArguments.Args = args
		return mLsArguments.validate()
	},
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Printf("ls arguments: %v", mLsArguments)
		var path string
		if len(args) > 1 {
			fmt.Printf("[WARN] avio 目前尚不支持一次列出多个路径，本次将只列出[%s]路劲下的文件信息", args[0])
			path = args[0]
		} else if len(args) == 0 {
			ex, err := os.Getwd()
			if err != nil {
				fmt.Printf("[WARN] get current path failed, %v", err)
			}
			path = ex
			fmt.Println(path)
		} else {
			path = args[0]
		}

		var depth uint = 2
		if mLsArguments.Recursive {
			depth = 10
		}

		walker := util.NewDFWalker(path, depth, mLsArguments.Full)
		defer walker.Close()
		g, e := walker.Walk()

		if e != nil {
			fmt.Printf("failed to walk, err: %s", e)
			return
		}

		var count uint32 = 0 // max files to list 4,294,967,295
		msg := g.Next()
		for msg.Err == nil && !msg.EOF {
			count++
			if !mLsArguments.Count {
				handlePath(msg.Info, msg.Path)
			}
			msg = g.Next()
		}
		printTotal(count)
	},
}

func handlePath(info os.FileInfo, filepath string) {
	sizeStr := fmt.Sprintf("%d", info.Size())
	if mLsArguments.Humanable {
		sizeStr = util.ByteCountBinary(info.Size())
	}
	fmt.Printf("%s\t%s\t%s\t%s\n", info.Mode(), sizeStr, info.ModTime().Format("1 02 15:04"), filepath)
}

func printTotal(count uint32) {
	fmt.Printf("total %d\n", count)
}
