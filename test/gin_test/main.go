package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Person struct {
	Name string `uri:"name" form:"name"`
}

func main() {
	e := gin.Default()
	e.GET("/hi/:name", func(ctx *gin.Context) {
		var person Person
		err := ctx.ShouldBind(&person)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    "500",
				"message": err.Error(),
			})
		}
		ctx.String(http.StatusOK, "hi! %s", person.Name)
	})
	e.Run(":8080")
}
