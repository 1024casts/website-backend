package router

import (
	"net/http"

	_ "1024casts/backend/docs" // docs is generated by Swag CLI, you have to import it.
	"1024casts/backend/router/middleware"

	"1024casts/backend/handler/comment"
	"1024casts/backend/handler/course"
	"1024casts/backend/handler/order"
	"1024casts/backend/handler/plan"
	"1024casts/backend/handler/sd"
	"1024casts/backend/handler/user"

	"1024casts/backend/handler/qiniu"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// Load loads the middlewares, routes, handlers.
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Middlewares.
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(mw...)
	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	// swagger api docs
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// pprof router
	pprof.Register(g)

	// api for authentication functionalities
	g.POST("/login", user.Login)
	g.POST("/logout", user.Logout)

	// The user handlers, requiring authentication
	u := g.Group("/v1/users")
	u.Use(middleware.AuthMiddleware())
	{
		u.POST("", user.Create)
		u.DELETE("/:id", user.Delete)
		u.PUT("/:id", user.Update)
		u.GET("", user.List)
		u.GET("/token", user.Get)
		//u.GET("/:id", user.Get)
		u.PUT("/:id/status", user.UpdateStatus)
	}

	// The course handlers, requiring authentication
	c := g.Group("/v1/courses")
	c.Use(middleware.AuthMiddleware())
	{
		c.POST("", course.Create)
		c.DELETE("/:id", user.Delete)
		c.PUT("/:id", course.Update)
		c.GET("", course.List)
		c.GET("/:id", course.Get)
		c.GET("/:id/sections", course.Section)
	}

	// The comment handlers, requiring authentication
	cmt := g.Group("/v1/comments")
	cmt.Use(middleware.AuthMiddleware())
	{
		cmt.GET("", comment.List)
	}

	// The order handlers, requiring authentication
	o := g.Group("/v1/orders")
	o.Use(middleware.AuthMiddleware())
	{
		o.GET("", order.List)
	}

	// The plan handlers, requiring authentication
	p := g.Group("/v1/plans")
	p.Use(middleware.AuthMiddleware())
	{
		p.GET("", plan.List)
		p.GET("/:id", plan.Get)
		//p.GET("/:alias", plan.Get)
	}

	// The plan handlers, requiring authentication
	q := g.Group("/v1/qiniu")
	q.Use(middleware.AuthMiddleware())
	{
		q.GET("", qiniu.List)
	}

	// The health check handlers
	svcd := g.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
	}

	return g
}
