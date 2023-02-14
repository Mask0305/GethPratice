package router

import (
	"fmt"
	"scp/app/interface/controller"

	"github.com/gin-gonic/gin"
)

type Router interface {
	Set()
}

type router struct {
	g *gin.Engine
	c controller.Controller
	m controller.MetaMaskController
}

func New(g *gin.Engine, c controller.Controller, m controller.MetaMaskController) Router {
	return &router{
		g: g,
		c: c,
		m: m,
	}
}

func (r *router) Set() {
	r.g.GET("/", func(ctx *gin.Context) {
		ctx.AbortWithStatus(200)
	})
	r.g.Use(CORSMiddleware())
	r.c.Router(r.g)
	r.m.Router(r.g)
}

// CORSMiddleware ...
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH,OPTIONS,GET,PUT")

		fmt.Println(c.Request.Method)
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
		}

		c.Next()
	}
}
