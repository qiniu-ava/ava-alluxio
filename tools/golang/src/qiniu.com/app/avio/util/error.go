package util

type AvioErrorCode uint

const (
	ARGUMENT_ERROR_ROOT AvioErrorCode = 100000
	ARGUMENT_ERROR_LIST AvioErrorCode = iota + 100100
	ARGUMENT_ERROR_PRELOAD
	ARGUMENT_ERROR_STAT
	ARGUMENT_ERROR_CP
	ARGUMENT_ERROR_MV
	ARGUMENT_ERROR_JOB

	WALK_ERROR_INVALID_STATUS AvioErrorCode = iota + 100200
)

var ErrorMessage map[AvioErrorCode]string = map[AvioErrorCode]string{
	ARGUMENT_ERROR_ROOT:    "root 命令参数错误",
	ARGUMENT_ERROR_LIST:    "ls 命令参数错误",
	ARGUMENT_ERROR_PRELOAD: "preload 命令参数错误",
	ARGUMENT_ERROR_STAT:    "stat 命令参数错误",
	ARGUMENT_ERROR_CP:      "cp 命令参数错误",
	ARGUMENT_ERROR_MV:      "mv 命令参数错误",
	ARGUMENT_ERROR_JOB:     "job 命令参数错误",

	WALK_ERROR_INVALID_STATUS: "invalid walk status to increase",
}

type AvioError struct {
	code AvioErrorCode
	msg  string
}

func NewAvioError(code AvioErrorCode, msg string) *AvioError {
	return &AvioError{code, msg}
}

func (a *AvioError) Error() string {
	return a.msg
}

func (a *AvioError) Code() AvioErrorCode {
	return a.code
}
