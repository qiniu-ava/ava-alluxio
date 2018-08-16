package cmd

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sync"

	"qiniu.com/app/avio/constants"

	"github.com/spf13/cobra"
	"qiniu.com/app/avio/api"
	util "qiniu.com/app/avio/util"
	"qiniu.com/app/common/typo"
)

var (
	zoneRP = regexp.MustCompile("//(z[0-9]){0, 1}/")
)

type FetchMessage struct {
	path    string
	ok      bool
	message string
}

type WorkerMetrics struct {
	jobs           int64
	doneJobs       int64
	successfulJobs int64
	failedJobs     int64
}

type Preload struct {
	uid          int64
	input        string
	prefix       string
	isJSONList   bool
	client       http.Client
	path         string
	log          string
	listFinished bool
	readWg       sync.WaitGroup
	readCh       chan string
	messageWg    sync.WaitGroup
	messageCh    chan FetchMessage
	metrics      *WorkerMetrics
	metricsMu    sync.Mutex
}

func (c *Preload) init(flags *loadFlags) error {
	c.input = flags.Source
	c.prefix = ""

	c.isJSONList = false
	if flags.IsJSONList == "true" {
		c.isJSONList = true
	}

	c.path = ""
	if len(flags.Args) == 1 {
		c.path = flags.Args[0]
	}

	c.metrics = &WorkerMetrics{
		jobs:           0,
		doneJobs:       0,
		successfulJobs: 0,
		failedJobs:     0,
	}

	u, e := util.GetUID()
	if e != nil {
		fmt.Printf("获取用户信息出错，请确定你是在 ava 的工作台中运行")
		return e
	}
	c.uid = u
	return nil
}

func (c *Preload) start() error {
	if c.path != "" {
		return c.startLocalPath()
	}
	return c.startFileList()
}

func (c *Preload) startLocalPath() error {
	p, e := util.ResolveAlluxioPath(c.path, c.uid)
	if e != nil {
		return e
	}

	job := typo.JobSpec{
		UID:  c.uid,
		Name: util.NewJobName(typo.PreloadJobType),
		Type: typo.PreloadJobType,
		Params: typo.JobParams{
			AlluxioURI:   p,
			Depth:        constants.MAX_DEPTH,
			FromFileList: false,
		},
	}

	if e := api.NewJob(c.client, job); e != nil {
		fmt.Printf("创建任务失败")
		return e
	}

	if e := api.StartJob(c.client, job.Name, c.uid); e != nil {
		fmt.Printf("任务启动失败")
		return e
	}

	fmt.Printf("启动任务成功，任务名:【%s】", job.Name)
	return nil
}

func (c *Preload) startFileList() error {
	return nil
}

type loadFlags struct {
	Source     string
	IsJSONList string
	Depth      uint
	Args       []string
}

func (f *loadFlags) Validate() error {
	if f.Depth > 10 || f.Depth < 1 {
		return util.NewAvioError(util.ARGUMENT_ERROR_PRELOAD, "请设置合法的参数，depth 参数最大值为 10，最小值为 1")
	}

	if len(f.Args) > 1 {
		return util.NewAvioError(util.ARGUMENT_ERROR_PRELOAD, "当前只支持一次加载一个指定目录或者文件，一次加载多个目录或者文件的功能将在后续支持")
	}

	if len(f.Args) == 1 && f.Source != "" {
		return util.NewAvioError(util.ARGUMENT_ERROR_PRELOAD, "-i/--input-file 参数和 path 只能同时设置一个")
	}

	if f.IsJSONList != "false" && f.IsJSONList != "true" {
		return util.NewAvioError(util.ARGUMENT_ERROR_PRELOAD, "--is-jsonlist 参数的可选项为 false/true")
	}

	if f.IsJSONList == "true" && f.Source != "" {
		return util.NewAvioError(util.ARGUMENT_ERROR_PRELOAD, "--is-jsonlist 参数只有在指定了 -i/--input-file 参数时有效")
	}

	if f.Source != "" && f.Depth != 4 {
		fmt.Printf("[WARN] -d/--depth 被设置为 %d，在使用 avio preload -i/--input-file=<list_file> 的模式下将被忽略", f.Depth)
	}
	return nil
}

var mLoadFlags = &loadFlags{}

func init() {
	preloadCmd.Flags().UintVarP(&mLoadFlags.Depth, "depth", "d", 4, "递归加载的最大深度，默认值为 4，最大值为 10")
	preloadCmd.Flags().StringVarP(&mLoadFlags.Source, "input-source", "i", "", "所有需要 preload 的文件的列表文件")
	preloadCmd.Flags().StringVar(&mLoadFlags.IsJSONList, "is-jsonlist", "false", "只在 -i/--input-file 被设置时有效，可选项：true/false，默认值为 false")
	preloadCmd.SetUsageTemplate(`Usage:{{if .Runnable}}
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
	rootCmd.AddCommand(preloadCmd)
}

var preloadAliases = []string{"ld", "pl", "load"}
var preloadCmd = &cobra.Command{
	Use:                "preload",
	Short:              "将指定目录或文件加载到 alluxio 缓存系统中",
	Long:               `将指定目录或文件加载到 alluxio 缓存系统中，此操作可能导致其他最近未使用的数据被踢出 alluxio 缓存系统或从缓存系统中更高级别的缓存降级到更低级别的缓存中。`,
	Aliases:            preloadAliases,
	DisableSuggestions: false,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		mLoadFlags.Args = args
		return mLoadFlags.Validate()
	},
	Run: func(cmd *cobra.Command, args []string) {
		mCMD := &Preload{}
		if e := mCMD.init(mLoadFlags); e != nil {
			fmt.Println(e)
			os.Exit(1)
		}

		if e := mCMD.start(); e != nil {
			fmt.Println(e)
			os.Exit(1)
		}

		fmt.Println("")
	},
}
