package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"qiniu.com/app/avio/constants"
	"qiniu.com/app/avio/util"
	"qiniu.com/app/common/typo"

	"github.com/spf13/cobra"
)

type Job struct {
	client http.Client
	uid    int64
}

type jobFlags struct {
	All    bool
	Delete bool
	Args   []string
}

func (j *jobFlags) Validate() error {
	if j.All && len(j.Args) > 0 {
		return util.NewAvioError(util.ARGUMENT_ERROR_JOB, "查看指定任务详情时，--all/-a 无效")
	}

	if j.Delete && len(j.Args) == 0 {
		return util.NewAvioError(util.ARGUMENT_ERROR_JOB, "请指定要删除的任务名称")
	}

	if len(j.Args) > 1 {
		return util.NewAvioError(util.ARGUMENT_ERROR_JOB, "每次只能查看一个任务的详情")
	}

	return nil
}

var mJobFlags jobFlags

func (j *Job) init(flags jobFlags) {
	j.client = http.Client{}
	u, e := util.GetUID()
	if e != nil {
		fmt.Println("获取用户信息出错，请确保你是在 ava 工作台中使用此工具")
		os.Exit(1)
	}

	j.uid = u
}

func (j *Job) start() {
	if len(mJobFlags.Args) == 0 {
		j.handleList()
	} else if mJobFlags.Delete {
		j.deleteItem()
	} else {
		j.handleItem()
	}
}

func (j *Job) handleList() {
	limit := 20
	if mJobFlags.All {
		limit = 1000
	}

	req := http.Request{
		Header: http.Header{
			"X-UID": []string{fmt.Sprintf("%d", j.uid)},
		},
		URL: &url.URL{
			Scheme:   "http",
			Host:     constants.AVIO_SERVICE_HOST,
			Path:     "/jobs",
			RawQuery: fmt.Sprintf("limit=%d&uid=%d", limit, j.uid),
		},
		Method: "GET",
	}
	timeToRetry := 3
retry:
	res, e := j.client.Do(&req)
	if e != nil && timeToRetry > 0 {
		if timeToRetry > 0 {
			goto retry
		} else {
			if res.StatusCode%100 == 5 {
				fmt.Printf("网络请求失败，服务器端故障。")
			} else {
				fmt.Printf("网络请求失败，请确定你的网络状况良好。")
			}
			return
		}
	}
	defer res.Body.Close()
	var result typo.ListJobResult
	if e = json.NewDecoder(res.Body).Decode(&result); e != nil {
		fmt.Printf("获取数据失败")
		return
	}

	if len(result.Items) != result.Total {
		fmt.Printf("total %d, list %d\n", result.Total, len(result.Items))
	}
	if result.Total == 0 {
		fmt.Println("you have no active jobs")
	} else {
		fmt.Printf("%-24s%12s%12s%18s%18s      %-30s\n", "name", "type", "status", "created", "lastUpdate", "message")
	}
	for _, item := range result.Items {
		fmt.Printf("%-24s%12v%12s%18s%18s      %-30s\n", item.Name, item.Type.String(), item.Status, util.DurationString(time.Since(*item.CreateTime)), util.DurationString(time.Since(*item.UpdateTime)), item.Message)
	}
	if len(result.Items) != result.Total {
		fmt.Printf("total %d, list %d\n", result.Total, len(result.Items))
	}
}

func (j *Job) handleItem() {
	jobName := mJobFlags.Args[0]
	req := http.Request{
		Header: http.Header{
			"X-UID": []string{fmt.Sprintf("%d", j.uid)},
		},
		URL: &url.URL{
			Scheme: "http",
			Host:   constants.AVIO_SERVICE_HOST,
			Path:   fmt.Sprintf("/jobs/%s", jobName),
		},
		Method: "GET",
	}
	timeToRetry := 3
retry:
	res, e := j.client.Do(&req)
	if e != nil && timeToRetry > 0 {
		if timeToRetry > 0 {
			goto retry
		} else {
			if res.StatusCode%100 == 5 {
				fmt.Printf("网络请求失败，服务器端故障。")
			} else {
				fmt.Printf("网络请求失败，请确定你的网络状况良好。")
			}
			return
		}
	}
	defer res.Body.Close()
	var info typo.JobInfo
	if e = json.NewDecoder(res.Body).Decode(&info); e != nil {
		fmt.Printf("获取数据失败")
		return
	}

	fmt.Printf(`name:	%s
type:	%s
create time:	%v
update time:	%v
finish time:	%v
status:	%s
message:	%s
path:	%s
depth:	%d
from file list:	%t`,
		info.Name,
		info.Type.String(),
		info.CreateTime,
		info.UpdateTime,
		info.FinishTime,
		info.Status,
		info.Message,
		info.Params.AlluxioURI,
		info.Params.Depth,
		info.Params.FromFileList)
}

func (j *Job) deleteItem() {
	jobName := mJobFlags.Args[0]
	req := http.Request{
		Header: http.Header{
			"X-UID": []string{fmt.Sprintf("%d", j.uid)},
		},
		URL: &url.URL{
			Scheme: "http",
			Host:   constants.AVIO_SERVICE_HOST,
			Path:   fmt.Sprintf("/jobs/%s", jobName),
		},
		Method: "DELETE",
	}
	timeToRetry := 3
retry:
	res, e := j.client.Do(&req)
	if e != nil && timeToRetry > 0 {
		if timeToRetry > 0 {
			goto retry
		} else {
			if res.StatusCode%100 == 5 {
				fmt.Printf("网络请求失败，服务器端故障。")
			} else {
				fmt.Printf("网络请求失败，请确定你的网络状况良好。")
			}
			return
		}
	}
	fmt.Printf("删除【%s】任务成功", jobName)
}

func init() {
	jobCmd.Flags().BoolVarP(&mJobFlags.All, "all", "a", false, "是否显示所有的任务，默认只显示最近 20 个任务")
	jobCmd.Flags().BoolVarP(&mJobFlags.Delete, "delete", "d", false, "删除指定的任务，需要指定任务名称")
	jobCmd.SetUsageTemplate(`Usage:{{if .Runnable}}
	{{.UseLine}} [job_name] {{end}}{{if .HasAvailableSubCommands}}
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
	rootCmd.AddCommand(jobCmd)
}

var jobAliases = []string{"j", "jb", "job"}
var jobCmd = &cobra.Command{
	Use:                "jobs",
	Short:              "查看我的后台任务",
	Long:               `查看我的后台任务列表或指定任务的详情。`,
	Aliases:            jobAliases,
	DisableSuggestions: false,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		mJobFlags.Args = args
		return mJobFlags.Validate()
	},
	Run: func(cmd *cobra.Command, args []string) {
		mCMD := &Job{}
		mCMD.init(mJobFlags)
		mCMD.start()
		fmt.Println("")
	},
}
