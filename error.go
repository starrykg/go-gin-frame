package gf

import "strconv"

type Error struct {
	Code int
	Msg  string
	Data interface{}
	Err  interface{}
}

func (err Error) Error() string {
	if err.Msg == "" {
		return CodeMapping.GetCodeInfo(err.GetCode())
	}
	return err.Msg
}

func (err Error) GetMsg() string {
	if err.Msg == "" {
		err.Msg = CodeMapping.GetCodeInfo(err.GetCode())
		if err.Msg == "" {
			return "code:" + strconv.Itoa(err.GetCode())
		}
	}
	return err.Error()
}

func (err Error) GetCode() int {
	if err.Code == 0 {
		return FailCode
	}
	return err.Code
}

func (err Error) IsNil() bool {
	return err.Msg == "" && err.Code == 0
}

func (err Error) GetDetail() interface{} {
	if err.Data != nil {
		return err.Data
	} else if err.IsNil() {
		return err.Err
	}
	return nil
}

func NewError(err error) Error {
	return Error{Err: err}
}

//传语言,I18N国际化
func NewErrorCode(code int) Error {
	return Error{Code: code}
}

//传语言,I18N国际化
func NewResponse(code int, msg string, data interface{}) Error {
	return Error{Code: code, Msg: msg, Data: data}
}

func NewErrorStr(msg string) Error {
	return Error{Code: FailCode, Msg: msg}
}
