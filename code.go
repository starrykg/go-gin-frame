package gf

type codeMapping map[int]string

func (cm codeMapping) GetCodeInfo(code int) string {
	if v, ok := cm[code]; ok {
		return v
	}
	return ""
}

func (cm codeMapping) AddCodeInfo(code int, msg string) {
	cm[code] = msg
}

const (
	SuccessCode    = 200
	FailInternal   = 400
	FailCode       = 500
	CodeParamError = 400001
)

var CodeMapping = codeMapping{
	FailCode:       "operation failed",
	FailInternal:   "network error",
	SuccessCode:    "successful",
	CodeParamError: "invalid parameter",
}
