package main

import (
	"fmt"
	"goshop/configs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func httpServerStart(cfg *configs.Config) bool {
	router := gin.Default()
	gin.DisableConsoleColor()
	// f, _ := os.Create("gin.log")
	// gin.DefaultWriter = io.MultiWriter(f)

	registerRoute(router)

	addr := fmt.Sprintf("%s:%s", cfg.TimeCfg.Host, cfg.TimeCfg.Port)
	go func() {
		if err := router.Run(addr); err != nil {
			glog.Errorln("[TimeServer] gin route run failed: ", err.Error())
			return
		}
	}()

	glog.Infof("[TimeServer] http server start on [%s]\n", addr)
	return true
}

func registerRoute(r *gin.Engine) {
	r.GET("/hello", func(c *gin.Context) {
		glog.Infof("[TimeServer] Get hello message from [%v] !!!\n", c.ClientIP())
		c.JSON(http.StatusOK, gin.H{
			"msg": "hello from time server",
		})
	})
}
