package gos

import (
	"gos/web"
)

func main() {
	app := web.Default()
	app.Get("/test", func(ctx *web.GosContext) {

	})

}
