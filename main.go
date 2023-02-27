package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/joho/godotenv"

	"github.com/labstack/echo"
)

func main() {
	// 载入env
	err := godotenv.Load(".env")
	if err != nil {
		return
	}

	// GPT地址
	GPTUrl := "https://api.openai.com/v1"

	// 创建 Echo 实例
	ECHO := echo.New()

	// 测试用
	ECHO.GET("/api/test", func(c echo.Context) error {
		GPTKEY := os.Getenv("GPT_KEY")
		return c.String(http.StatusOK, "Hello, Echo GO!"+GPTKEY)
	})

	// 处理图片
	ECHO.POST("/api/gpt/image", func(c echo.Context) error {
		// 地址
		target, setProxyError := url.Parse(GPTUrl + "/images/generations")
		if setProxyError != nil {
			return setProxyError
		}

		// 创建反向代理
		proxy := httputil.NewSingleHostReverseProxy(target)

		// 修改请求头
		c.Request().Header.Set("X-Forwarded-Host", c.Request().Host)

		// 转发请求
		proxy.ServeHTTP(c.Response().Writer, c.Request())

		return nil
	})

	// 启动服务器
	ECHO.Start(":2000")
}
