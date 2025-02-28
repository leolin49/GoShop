package main

import (
	"fmt"
	"goshop/configs"
	errorcode "goshop/pkg/error"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func httpServerStart(cfg *configs.Config) bool {
	router := gin.Default()
	gin.DisableConsoleColor()
	// f, _ := os.Create("gin.log")
	// gin.DefaultWriter = io.MultiWriter(f)

	registerRoute(router)

	addr := fmt.Sprintf("%s:%s", cfg.GatewayCfg.Host, cfg.GatewayCfg.Port)
	go func() {
		if err := router.Run(addr); err != nil {
			glog.Errorln("[Gatewayserver] gin route run failed: ", err.Error())
			return
		}
	}()

	glog.Infof("[Gatewayserver] http server start on [%s]\n", addr)
	return true
}

func registerRoute(r *gin.Engine) {
	r.GET("/hello", func(c *gin.Context) {
		glog.Infof("[Gatewayserver] Get hello message from [%v] !!!\n", c.ClientIP())
		c.JSON(http.StatusOK, gin.H{
			"x": "hello world",
		})
	})

	user := r.Group("/user")
	{
		user.POST("/register", handleRegister)
		user.POST("/login", handleLogin)
	}

	product := r.Group("/product")
	{
		product.POST("/add", handleAddProduct) // TODO: manager
		product.POST("/list", handleListProducts)
		product.POST("/get", handleGetProduct)
		product.POST("/search", handleSearchProducts)
	}

	cart := r.Group("/cart")
	cart.Use(JwtAuthMiddleware())
	{
		cart.POST("/get", handleGetCart)
		cart.POST("/add", handleAddCart)
		cart.POST("/clean", handleCleanCart)
	}

	checkout := r.Group("/checkout")
	checkout.Use(JwtAuthMiddleware())
	{
		checkout.POST("/checkout", handleCheckout)
	}

	stock := r.Group("/stock")
	{
		stock.POST("/get", handleGetStock)
		stock.POST("/add", handleAddStock) // TODO: manager
		stock.POST("/sub", handleSubStock)
	}

	r.POST("/refreshToken", handleRefreshToken)
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

func rpcRequestError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
	})
}
