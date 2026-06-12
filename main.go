package main

import (
	"log"

	"project/configs"
	"project/handlers"
	"project/logger"
	"project/middlewares"
	"project/models"
	"project/services"
	"project/utils"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {

	// ===== init =====

	// 加载配置
	cfg, err := configs.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化logger
	logger.Init(cfg.Server.Mode)
	logger.Info("config loaded", "mode", cfg.Server.Mode)

	// 连接数据库
	db, err := gorm.Open(sqlite.Open(cfg.Database.DBName), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect database", "err", err)
	}
	logger.Info("database connected", "db", cfg.Database.DBName)

	// models迁移
	if err := db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}); err != nil {
		logger.Fatal("Failed to migrate models to database", "err", err)
	}
	logger.Info("database migrated")

	// 服务初始化
	userService := services.NewUserService(db)
	userHandler := handlers.NewUserHandler(userService, []byte(cfg.JWT.Secret))
	postService := services.NewPostService(db)
	postHandler := handlers.NewPostHandler(postService, []byte(cfg.JWT.Secret))
	commentService := services.NewCommentService(db)
	commentHandler := handlers.NewCommentHandler(commentService, []byte(cfg.JWT.Secret))

	// ===== gin =====

	// 创建Gin引擎
	r := gin.New()

	// 全局中间件
	r.Use(middlewares.Recovery())
	r.Use(middlewares.Logger())

	// 公开路由
	r.GET("/health", func(c *gin.Context) {
		utils.Success(c, gin.H{
			"status": "ok",
		})
	})

	public := r.Group("/api/v1")
	{
		// 认证
		public.POST("/auth/register", userHandler.Register)
		public.POST("/auth/login", userHandler.Login)

		// 文章 - 读
		public.GET("/posts", postHandler.GetPostList)
		public.GET("/posts/:id", postHandler.GetPostProfile)

		// 评论 - 读
		public.GET("/posts/:id/comments", commentHandler.GetCommentList)

	}

	// 需鉴权路由
	protected := r.Group("/api/v1")
	protected.Use(middlewares.Auth([]byte(cfg.JWT.Secret)))
	{
		// 文章 - 写
		protected.POST("/posts", postHandler.CreatePost)
		protected.PUT("/posts/:id", postHandler.UpdatePost)
		protected.DELETE("/posts/:id", postHandler.DeletePost)

		// 评论 - 写
		protected.POST("/posts/:id/comments", commentHandler.CreateComment)

	}

	// 启动服务器
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	logger.Info("server starting", "addr", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server", "err", err)
	}

}
