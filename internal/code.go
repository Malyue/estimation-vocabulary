package internal

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeServerBusy
	CodeInvalidTestId
	CodeWordSelectErr
	CodeWordRepeat
	CodeErrFileFormat
	CodeErrJsonFormat
	CodeErrParseBody
	CodeErrParseInt
	CodeLevelInvalid
	CodeLevelRequire
	CodeWrongLevel
	CodeRecordNotFound
	CodeOops
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:        "success",
	CodeServerBusy:     "服务繁忙",
	CodeInvalidTestId:  "非法测试id",
	CodeWordSelectErr:  "单词查询错误",
	CodeWordRepeat:     "单词获取重复",
	CodeErrFileFormat:  "文件格式错误，需要json文件",
	CodeErrJsonFormat:  "json格式错误，无法解析",
	CodeErrParseBody:   "解析body错误",
	CodeErrParseInt:    "ParseInt错误",
	CodeLevelInvalid:   "非法难度",
	CodeLevelRequire:   "需要正确的level参数",
	CodeWrongLevel:     "错误的level",
	CodeRecordNotFound: "数据库中不存在该单词",
	CodeOops:           "Oops,该组数据无效", //批处理单词在数据库中找不到为数据无效，无效数据过多测试组舍弃
}

// Msg 返回特定的错误提示信息
func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		return codeMsgMap[CodeServerBusy]
	}
	return msg
}
