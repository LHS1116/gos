package main

import "gos/web"

func Test(c *web.GosContext) {
	c.JSON(200, web.H{
		"title": "Test",
		"data":  "success",
	})
}
