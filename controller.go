package gf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitee.com/kingsingnal/util/gconv"
	"github.com/gin-gonic/gin"
	logs "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type Context struct {
	*gin.Context
}

type ProcessFunc func(controller *Controller) IError
type Controller struct {
	TrackId string
	Ctx     *Context
	Code    int
	Msg     string
	error   IError
	Data    interface{}
	Req     interface{}
	// middleware functions
	Middleware []ProcessFunc
	// logic functions
	ProcessFun ProcessFunc
}

func GetNewController(g *gin.Context, req interface{}) *Controller {
	return &Controller{
		Ctx: &Context{
			Context: g,
		},
		Code: SuccessCode,
		Req:  req,
	}
}

//alias name, If run it, do not use RunProcess
func RunWithRequest(req interface{}, g *gin.Context) {
	ct := GetNewController(g, req)
	c := GetContext(g)
	var e IError
	defer func() {
		if x := recover(); x != nil {
			logs.Error(" panic :", x)
			e = Error{
				Code: FailCode,
				Msg:  "runtime internal error",
				Err:  x,
			}
		}
		if e != nil {
			ct.SetError(e)
		}

		var err error
		c.Status(successCode)
		c.Header("Content-Type", "application/json; charset=utf-8")
		defer func() {
			if err != nil {
				errStr := fmt.Sprintf(`{"code":"%d","msg":"system error: %s","data":null}`, FailInternal, err.Error())
				rStr(c, errStr)
			}
		}()

		res := ct.getResponse()

		rJson(c, res)
	}()

	if e = getTrackId(ct); e != nil {
		return
	}
	if e = runProcess(ct); e != nil {
		return
	}
}

func RunProcess(controller IController, g *gin.Context) {
	c := GetContext(g)
	controller.SetContext(c)
	var e IError
	defer func() {
		if x := recover(); x != nil {
			logs.Error(" panic :", x)
			e = Error{
				Code: FailCode,
				Msg:  "runtime internal error",
				Err:  x,
			}
		}
		if e != nil {
			controller.SetError(e)
		}

		var err error
		c.Status(successCode)
		c.Header("Content-Type", "application/json; charset=utf-8")
		defer func() {
			if err != nil {
				errStr := fmt.Sprintf(`{"code":"%d","msg":"system error: %s","data":null}`, FailInternal, err.Error())
				rStr(c, errStr)
			}
		}()

		res := controller.getResponse()

		rJson(c, res)
	}()

	if e = getTrackId(controller); e != nil {
		return
	}
	if e = runProcess(controller); e != nil {
		return
	}
}

func getTrackId(controller IController) IError {
	//默认采用21位的纳秒时间值(年份仅后两位)
	trackId := gconv.GetRequestKey(controller.GetContext().Request, "track_id")
	if trackId == "" {
		now := time.Now()
		trackId = fmt.Sprintf("%s%09d", now.Format("060102150405"), now.UnixNano()%1e9)
	}
	controller.SetTrackId(trackId)
	return nil
}

//run controller process
func runProcess(controller IController) (err IError) {
	err = controller.Decode()
	if err != nil {
		//controller.SetError(err)
		return
	}

	err = controller.Process()
	if err != nil {
		return
	}

	return
}

func (controller *Controller) Use(fn func(controller *Controller) IError) {
	if controller.Middleware == nil {
		controller.Middleware = make([]ProcessFunc, 1)
	}
	controller.Middleware = append(controller.Middleware, fn)
}

// controller default Decode
func (controller *Controller) Decode() IError {
	controller.Data = nil

	switch controller.Ctx.Context.Request.Method {
	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		fallthrough
	case http.MethodDelete:
		ct := controller.Ctx.Context.Request.Header.Get("Content-Type")
		if strings.Contains(ct, "json") {
			bt, err := controller.Ctx.GetRawData()
			if err == nil {
				if len(bt) == 0 {
					bt = []byte("{}")
				}
				d := json.NewDecoder(bytes.NewReader(bt))
				d.UseNumber()
				if err := d.Decode(&controller.Req); err != nil {
					return NewErrorCode(CodeParamError)
				}
			}
		}
	default:

	}

	return nil
}

// controller default Process
func (controller *Controller) Process() IError {
	if controller.Middleware != nil {
		if len(controller.Middleware) > 0 {
			for _, m := range controller.Middleware {
				err := m(controller)
				if err != nil {
					return err
				}
			}
		}
	}
	if controller.ProcessFun != nil {
		return controller.ProcessFun(controller)
	}
	return nil
}

//ger ready response
func (controller *Controller) getResponse() Response {
	if controller.error != nil {
		return getDefaultErrorResponse(controller.error)
	} else {
		return getResponseWithCode(controller.Code, controller.Msg, controller.Data)
	}
}

//set error
func (controller *Controller) SetError(err IError) {
	logs.Info("set error in %s, track_id: %s, err: %s\n", controller.Ctx.Context.Request.URL.Path, controller.GetTrackId(), err.GetMsg())
	controller.error = err
	return
}

//get trackId
func (controller *Controller) GetTrackId() string {
	return controller.TrackId
}

//set trackId
func (controller *Controller) SetTrackId(id string) {
	controller.TrackId = id
}

//set Context
func (controller *Controller) SetContext(c *Context) {
	if controller.Ctx == nil {
		controller.Ctx = c
	}
}

//get gin Content
func (controller *Controller) GetContext() *Context {
	return controller.Ctx
}

//query things with gin
func (controller *Controller) Query(key string) string {
	return controller.Ctx.Query(key)
}

//router params with gin
func (controller *Controller) Param(key string) string {
	return controller.Ctx.Param(key)
}
