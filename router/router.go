package router

import (
	"html/template"
	"net/http"

	_ "github.com/1024casts/1024casts/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/1024casts/1024casts/router/middleware"

	"github.com/1024casts/1024casts/handler/api/v1/comment"
	apiCourse "github.com/1024casts/1024casts/handler/api/v1/course"
	"github.com/1024casts/1024casts/handler/api/v1/order"
	"github.com/1024casts/1024casts/handler/api/v1/plan"
	"github.com/1024casts/1024casts/handler/api/v1/user"
	"github.com/1024casts/1024casts/handler/qiniu"
	"github.com/1024casts/1024casts/handler/sd"
	webCourse "github.com/1024casts/1024casts/handler/web/course"
	webPlan "github.com/1024casts/1024casts/handler/web/plan"
	"github.com/1024casts/1024casts/handler/web/wiki"

	"github.com/1024casts/1024casts/handler/api/v1/video"
	"github.com/1024casts/1024casts/handler/web"
	"github.com/1024casts/1024casts/handler/web/topic"
	"github.com/foolin/gin-template"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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
	//pprof.Register(g)

	// api for authentication functionalities
	g.POST("/login", user.Login)
	g.POST("/logout", user.Logout)

	v1Route := g.Group("/v1")
	v1Route.Use(middleware.AuthMiddleware())
	{
		// 课程
		apiCourse.Endpoint(v1Route)
		// 评论
		comment.Endpoint(v1Route)
		// 视频
		video.Endpoint(v1Route)
		// 用户
		user.Endpoint(v1Route)
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
		q.POST("/upload", qiniu.Upload)
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

func InitWebRouter(g *gin.Engine) *gin.Engine {
	router := gin.Default()

	router.Use(static.Serve("/static", static.LocalFile(viper.GetString("static"), false)))

	//new template engine
	router.HTMLRender = gintemplate.New(gintemplate.TemplateConfig{
		Root:      "templates",
		Extension: ".html",
		Master:    "layouts/master",
		Partials:  []string{},
		Funcs: template.FuncMap{
			"sub": func(a, b int) int {
				return a - b
			},
		},
		DisableCache: true,
	})

	router.GET("/", web.Index)
	router.GET("/courses", webCourse.Index)
	router.GET("/topics", topic.Index)
	router.GET("/topics/:id", topic.Detail)
	router.GET("/vip", webPlan.Index)
	router.GET("/wiki", wiki.Index)
	router.Run(":8099")

	return router
}
