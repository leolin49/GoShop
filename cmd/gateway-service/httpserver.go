package main

import (
	"fmt"
	errorcode "goshop/pkg/error"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func httpServerStart() bool {
	router := gin.Default()
	gin.DisableConsoleColor()
	// f, _ := os.Create("gin.log")
	// gin.DefaultWriter = io.MultiWriter(f)

	registerRoute(router)

	go func() {
		if err := router.Run("localhost:8080"); err != nil {
			glog.Errorln("[Gatewayserver] gin route run failed: ", err.Error())
			return
		}
	}()

	glog.Infoln("[Gatewayserver] http server start.")
	return true
}

func registerRoute(router *gin.Engine) {
	router.POST("/register", handleRegister)
	router.POST("/login", handleLogin)

	router.POST("/add_product", handleAddProduct)
	router.POST("/list_products", handleListProducts)
	router.POST("/get_product", handleGetProduct)
	router.POST("/search_product", handleSearchProducts)

	router.POST("/get_cart", handleGetCart)
	router.POST("/add_cart", handleAddCart)
	router.POST("/clean_cart", handleCleanCart)
}

func getPostFormInt(c *gin.Context, key string) (int, error) {
	val := c.PostForm(key)
	if val == "" {
		return 0, fmt.Errorf("missing parameter %s.", key)
	}
	return strconv.Atoi(val)
}

func invalidParam(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error_code": errorcode.InvalidParam,
	})
}

func rpcRequestError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error_code": errorcode.RpcRequestFailed,
	})
}
