package util

type AvioError struct {
	Msg string
}

func (a *AvioError) Error() string {
	return a.Msg
}

type WalkErrorCode uint

var (
	WalkStatusError WalkErrorCode = 0x000001
)

type WalkError struct {
	code WalkErrorCode
	msg  string
}

func newWalkError(code WalkErrorCode, msg string) *WalkError {
	return &WalkError{code, msg}
}

func (w *WalkError) Code() WalkErrorCode {
	return w.code
}

func (w *WalkError) Error() string {
	return w.msg
}
