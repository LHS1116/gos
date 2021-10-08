package main

import (
	"fmt"
	"gos/web"
	"time"
)

func Format(method, str string) string {
	t := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s] %s  (%s)", method, str, t)
	//[GET] 127.0.0.1:8080/test (2006-01-02 15:04:05)
}

func LogMiddleWare() web.HandleFunc {
	return func(c *web.GosContext) {
		curPath := c.FullPath()
		str := Format(c.Request.Method, curPath)
		fmt.Println(str)
		c.Next()
		fmt.Println("LOGGER OVER")
	}
}

func NextMiddleWare() web.HandleFunc {
	return func(c *web.GosContext) {
		fmt.Println("aha!")
		c.Next()
	}

}
