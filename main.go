package main

import (
	"gos/web"
)

func main() {
	app := web.Default()
	app.Get("/test", Test)
	app.Use(LogMiddleWare)
	app.Use(NextMiddleWare)
	err := app.Run(8099)
	if err != nil {
		panic(err)
	}
}
