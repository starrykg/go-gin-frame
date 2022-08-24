package gf

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"time"
)

const (
	DEFAULT     = "default"
	successCode = 200
)

var gf *Gf
var configOne sync.Once
var contextPool *sync.Pool
var res404 = Response{
	Code: 404,
	Msg:  "页面不存在",
	Data: nil,
}
var resSuccess = Response{
	Code: 200,
	Msg:  "错误信息",
	Data: CodeMapping,
}

func GetContext(g *gin.Context) *Context {
	c := contextPool.Get().(*Context)
	c.Context = g
	return c
}

func init() {
	gf = &Gf{
		Gin:    gin.New(),
		Log:    *log.New(),
		Config: &goconfig.ConfigFile{},
		Db:     make(map[string]*gorm.DB),
		Redis:  make(map[string]*redis.Client),
	}
	gf.Gin.NoRoute(func(context *gin.Context) {
		context.JSON(successCode, res404)
	})
	gf.Gin.GET("/codes", func(context *gin.Context) {
		context.JSON(successCode, resSuccess)
	})
	contextPool = &sync.Pool{New: func() interface{} {
		return &Context{}
	}}
}

func InitConfig(path string) {
	configOne.Do(func() {
		var err error
		GetGf().Config, err = goconfig.LoadConfigFile(path)
		if err != nil {
			panic(fmt.Sprintf("unable to load config file：%s \n", err))
		}
		ServiceName, _ = gf.Config.GetValue("", "SERVICE_NAME")
		gf.Gin.GET("/"+ServiceName+"/code", func(context *gin.Context) {
			context.JSON(200, codeMapping{})
		})
	})
}

func InitDb() {
	cf, err := gf.Config.GetValue("", "MYSQL_DSN")
	if err == nil {
		db, err := gorm.Open("mysql", cf)
		if err == nil {
			GetGf().Db[DEFAULT] = db
		}
	}

	if err != nil {
		log.Error("db init error: %s", err.Error())
	}
}

func InitRedis() {
	cf, err := GetGf().Config.GetSection("redis")
	if err == nil {
		c := redis.NewClient(&redis.Options{
			Network:      "tcp",
			Addr:         cf["host"] + ":" + cf["port"],
			Dialer:       nil,
			OnConnect:    nil,
			Password:     cf["password"],
			DB:           String2Int(cf["db"], 0),
			MaxRetries:   3,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			PoolSize:     String2Int(cf["poolSize"], 20),
			MinIdleConns: String2Int(cf["minIdleConn"], 5),
			MaxConnAge:   20,
			PoolTimeout:  3 * time.Second,
			IdleTimeout:  5 * time.Second,
		})
		if c.Ping().Err() == nil {
			GetGf().Redis[DEFAULT] = c
		}
	}

	if err != nil {
		log.Error("redis init error: %s", err.Error())
	}
}

// 转换默认值
func String2Int(s string, defaultVal int) int {
	a, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}

	return a
}
