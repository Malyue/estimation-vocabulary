package internal

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponeseCode struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// ResponseError 返回一个code 类型的错误
func ResponseError(c *gin.Context, code ResCode) {
	c.JSON(http.StatusOK, &ResponeseCode{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

//// ResponseErrorWithMsg 返回一个提示信息为自定义的错误
//func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
//	c.JSON(http.StatusOK, &ResponeseCode{
//		Code: code,
//		Msg:  code.Msg(),
//		Data: msg,
//	})
//}

func ResponseErrorWithData(c *gin.Context, code ResCode, data interface{}) {
	c.JSON(http.StatusOK, &ResponeseCode{
		Code: code,
		Msg:  code.Msg(),
		Data: data,
	})
}

// ResponseSuccess 返回成功的信息
func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponeseCode{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	})
}

// ResponseSuccess01 返回系统的状态的信息
func ResponseMsg(c *gin.Context, code ResCode) {
	c.JSON(http.StatusOK, &ResponeseCode{
		Code: code,
		Msg:  code.Msg(),
	})
}
