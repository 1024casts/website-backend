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
	webUser "github.com/1024casts/1024casts/handler/web/user"
	"github.com/1024casts/1024casts/handler/web/wiki"

	"github.com/1024casts/1024casts/handler/api/v1/video"
	"github.com/1024casts/1024casts/handler/web"
	"github.com/1024casts/1024casts/handler/web/topic"
	"github.com/foolin/gin-template"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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
	//g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// pprof router
	//pprof.Register(g)

	// api for authentication functionalities
	g.POST("/v1/login", user.Login)
	g.POST("/v1/logout", user.Logout)

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

	// Middlewares.
	router.Use(gin.Recovery())
	// 404 Handler.
	//g.NoRoute(func(c *gin.Context) {
	//	c.String(http.StatusNotFound, "The incorrect API route.")
	//})

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
	router.GET("/login", webUser.GetLogin)
	router.POST("/login", webUser.Login)
	router.GET("/register", webUser.GetRegister)
	router.POST("/register", webUser.Register)
	router.GET("/logout", webUser.Logout)
	router.GET("/users/:username", webUser.Profile)             // 个人资料
	router.GET("/users/:username/topics", webUser.Logout)       // 发表过的主题
	router.GET("/users/:username/replies", webUser.Logout)      // 回复过的
	router.GET("/users/:username/favorites", webUser.Logout)    // 收藏过的
	router.GET("/users/:username/following", webUser.Following) // 正在关注的人
	router.GET("/users/:username/followers", webUser.Follower)  // 关注者(粉丝)

	router.Use(middleware.CookieMiddleware()).GET("/settings", webUser.Setting)
	router.Use(middleware.CookieMiddleware()).GET("/notifications", webUser.Notification)

	router.GET("/courses", webCourse.Index)
	router.GET("/courses/:slug", webCourse.Detail)

	router.GET("/topics", topic.Index)
	router.GET("/topics/:id", topic.Detail)
	t := router.Group("/topic")
	t.Use(middleware.CookieMiddleware())
	{
		t.GET("/new", topic.Create)
		t.POST("/new", topic.Create)
	}

	router.GET("/vip", webPlan.Index)
	router.GET("/wiki", wiki.Index)
	router.Run(":8888")

	return router
}
