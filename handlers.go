package main

import "gos/web"

func Test(c *web.GosContext) {
	param1 := c.DefaultQuery("test", "Test not found")
	param2 := c.DefaultQuery("param", "Param not found")
	if param1 == "error" {
		panic(param1)
		return
	}
	c.JSON(200, web.H{
		"title":  "Test",
		"data":   "success",
		"params": []string{param1, param2},
	})
}

func TestQQ(c *web.GosContext) {
	param := c.PathParam("qq", "QQ NOT FOUND")
	c.JSON(200, web.H{
		"title":  "TestQQ",
		"data":   "success",
		"params": param,
	})

}

func TestName(c *web.GosContext) {
	param := c.PathParam("name", "Name NOT FOUND")
	c.JSON(200, web.H{
		"title":  "TestName",
		"data":   "success",
		"params": param,
	})

}
