package main

import (
	"gos/web"
)

func main() {
	app := web.Default()
	app.Get("/test/qqq/*", Test)
	app.Get("/:qq/", TestQQ)
	app.Get("/test/:name", TestName)
	app.PrintRouter()
	v1 := app.Group("/v1", TestName)
	{
		v1.Group("/v2", TestQQ)
	}

	app.Use(LogMiddleWare)
	//app.Use(NextMiddleWare)
	err := app.Run(8099)
	if err != nil {
		panic(err)
	}
}
