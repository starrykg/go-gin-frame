gin 再次封装
import gf "xxx/go-gin-frame"

	base := gf.GetGf()
	go func() {
		//测试环境开放
		if configure.CmSysConf.GinMode != "release" {
			base.Gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		}
		defer func() {
			if err := recover(); err != nil {
				logger.Error(xxx, "HttpServer server is panic")
			}
		}()
		// 服务连接
		if err := base.Gin.Run( IP + ":8899"); err != nil && err != http.ErrServerClosed {
			logger.Trace(xxx, "HttpServer is listenning")
		}
	}()