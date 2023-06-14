package internal

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeServerBusy
	CodeInvalidTestId
	CodeWordSelectErr
	CodeWordRepeat
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:       "success",
	CodeServerBusy:    "服务繁忙",
	CodeInvalidTestId: "非法测试id",
	CodeWordSelectErr: "单词查询错误",
	CodeWordRepeat:    "单词获取重复",
}

// Msg 返回特定的错误提示信息
func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		return codeMsgMap[CodeServerBusy]
	}
	return msg
}
