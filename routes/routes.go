package routes

import (
	"caretop-backend/handlers"
	"caretop-backend/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	r.Use(middleware.CORS())

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
			auth.POST("/logout", handlers.Logout)
			auth.GET("/me", middleware.AuthRequired(), handlers.GetMe)
		}

		products := api.Group("/products")
		{
			products.GET("", handlers.GetProducts)
			products.GET("/:slug", handlers.GetProduct)
		}

		blog := api.Group("/blog")
		{
			blog.GET("", handlers.GetBlogPosts)
			blog.GET("/categories", handlers.GetBlogCategories)
			blog.GET("/:slug", handlers.GetBlogPost)
			blog.POST("/:slug/like", middleware.AuthRequired(), handlers.LikeBlogPost)
		}

		adminBlog := api.Group("/admin/blog")
		adminBlog.Use(middleware.AdminRequired())
		{
			adminBlog.GET("", handlers.GetAllBlogPosts)
			adminBlog.POST("", handlers.CreateBlogPost)
			adminBlog.PUT("/:id", handlers.UpdateBlogPost)
			adminBlog.DELETE("/:id", handlers.DeleteBlogPost)
		}

		forum := api.Group("/forum")
		{
			forum.GET("/boards", handlers.GetBoards)
			forum.GET("/boards/:slug/threads", handlers.GetBoardThreads)
			forum.GET("/threads/:id", handlers.GetThread)
			forum.POST("/threads", middleware.AuthRequired(), handlers.CreateThread)
			forum.PUT("/threads/:id", middleware.AuthRequired(), handlers.UpdateThread)
			forum.DELETE("/threads/:id", middleware.AuthRequired(), handlers.DeleteThread)
			forum.POST("/threads/:id/reply", middleware.AuthRequired(), handlers.ReplyThread)
			forum.POST("/threads/:id/like", middleware.AuthRequired(), handlers.LikeThread)
			forum.POST("/threads/:id/collect", middleware.AuthRequired(), handlers.CollectThread)
			forum.GET("/my", middleware.AuthRequired(), handlers.GetMyPosts)
		}

		adminForum := api.Group("/admin/forum")
		adminForum.Use(middleware.AdminRequired())
		{
			adminForum.GET("/stats", handlers.GetForumStats)
		}

		tickets := api.Group("/tickets")
		tickets.Use(middleware.AuthRequired())
		{
			tickets.GET("", handlers.GetTickets)
			tickets.POST("", handlers.CreateTicket)
			tickets.GET("/:id", handlers.GetTicket)
			tickets.POST("/:id/reply", handlers.ReplyTicket)
		}

		adminTickets := api.Group("/admin/tickets")
		adminTickets.Use(middleware.AdminRequired())
		{
			adminTickets.GET("", handlers.GetAllTickets)
			adminTickets.PUT("/:id/status", handlers.UpdateTicketStatus)
		}

		admin := api.Group("/admin")
		admin.Use(middleware.AdminRequired())
		{
			admin.GET("/stats", handlers.GetStats)
			admin.GET("/users", handlers.GetUsers)
			admin.PUT("/users/:id/role", handlers.UpdateUserRole)
			admin.PUT("/users/:id/ban", handlers.BanUser)
		}
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
