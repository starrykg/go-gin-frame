package gf

import (
	"github.com/Unknwon/goconfig"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

var ServiceName string

type Gf struct {
	Gin    *gin.Engine
	Log    log.Logger
	Config *goconfig.ConfigFile
	Db     map[string]*gorm.DB
	Redis  map[string]*redis.Client
	Cache  map[string]ICache
}

// return instance config
func GetConfig() *goconfig.ConfigFile {
	return GetGf().Config
}

// return instance
func GetGf() *Gf {
	return gf
}

// return instance gin core
func GetGin() *gin.Engine {
	return GetGf().Gin
}

// return instance logger
func GetLog() log.Logger {
	return GetGf().Log
}

// get db if exists
func GetDb(name string) *gorm.DB {
	if name == "" {
		name = DEFAULT
	}
	return GetGf().GetDb(name)
}

// alias return  instance gin core
func GetRouter() IRouter {
	return GetGf().Gin
}

// alias gin group
func (gf *Gf) Group(path string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return gf.Gin.Group(path, handlers...)
}

// alias gin Use
func (gf *Gf) Use(middleware ...gin.HandlerFunc) gin.IRoutes {
	return gf.Gin.Use(middleware...)
}

// alias gin Handle
func (gf *Gf) Handle(httpMethod, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return gf.Gin.Handle(httpMethod, relativePath, handlers...)
}

// get default cache, if config set MYSQL_DSN
func (gf *Gf) GetDefaultCache() ICache {
	if db, ok := gf.Cache[DEFAULT]; ok {
		return db
	}
	return nil
}

// get cache with name
// if not exists, fatal err
func (gf *Gf) GetCache(name string) ICache {
	if db, ok := gf.Cache[name]; ok {
		return db
	}

	panic(name + " cache not exists")
}

// get default db, if config set MYSQL_DSN
func (gf *Gf) GetDefaultDb() *gorm.DB {
	if db, ok := gf.Db[DEFAULT]; ok {
		return db
	}
	return nil
}

// get cache with name
// now only support gorm
// if not exists, fatal err
func (gf *Gf) GetDb(name string) *gorm.DB {
	if db, ok := gf.Db[name]; ok {
		return db
	}

	panic(name + " db not exists")
}

// get default redis
func (gf *Gf) GetDefaultRedis() *redis.Client {
	if c, ok := gf.Redis[DEFAULT]; ok {
		return c
	}
	return nil
}

// get cache with name
// now only support gorm
func (gf *Gf) GetRedis(name string) *redis.Client {
	if c, ok := gf.Redis[name]; ok {
		return c
	}
	return nil
}
