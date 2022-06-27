package gf

import (
	"encoding/json"
	"net/http"
)

var ResponseDataNil = struct {
}{}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (res Response) GetBytes() []byte {
	b, _ := json.Marshal(res)
	return b
}

func (res Response) GetString() string {
	b, _ := json.Marshal(res)
	return string(b[:])
}

func getDefaultErrorResponse(err IError) Response {
	var data interface{}

	//I18N  TODO
	if err.Error() == "操作失败" {
		data = err.GetDetail()
	}

	if data == nil {
		data = ResponseDataNil
	}

	return Response{
		err.GetCode(), err.GetMsg(), data,
	}
}

func getResponseWithCode(code int, msg string, data ...interface{}) Response {
	if code == 0 {
		code = SuccessCode
	}
	if msg == "" {
		msg = CodeMapping.GetCodeInfo(code)
	}
	var r = Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	}

	if data == nil {
		return r
	}

	l := len(data)
	if l > 0 {
		if l == 1 {
			r.Data = data[0]
		} else {
			r.Data = data
		}
	}

	return r
}

func ResponseStr(c *Context, str string) {
	rStr(c, str)
}

func ResponseJson(c *Context, data interface{}) {
	rJson(c, data)
}

func rStr(c *Context, str string) {
	c.String(http.StatusOK, str)
}

func rJson(c *Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}
