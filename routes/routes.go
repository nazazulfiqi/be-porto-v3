package routes

import (
	"be-porto-v3/controller"
	"be-porto-v3/middleware"
	"be-porto-v3/repository"
	"be-porto-v3/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes mengatur semua route di aplikasi
func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Inisialisasi Repository
	projectRepo := repository.NewProjectRepository(db)
	authRepo := repository.NewAuthRepository(db)
	articleRepo := repository.NewArticleRepository(db)

	// Inisialisasi Service
	projectService := service.NewProjectService(projectRepo)
	authService := service.NewAuthService(authRepo)
	articleService := service.NewArticleService(articleRepo)

	// Inisialisasi Controller
	projectController := controller.NewProjectController(projectService)
	authController := controller.NewAuthController(authService)
	articleController := controller.NewArticleController(articleService)

	// Routes untuk autentikasi
	auth := r.Group("/auth")
	auth.POST("/register", authController.Register)
	auth.POST("/login", authController.Login)

	// Routes untuk project (dengan middleware JWT)
	projects := r.Group("/projects")
	projects.Use(middleware.AuthMiddleware()) // 🔥 Tambahkan middleware di sini
	projects.GET("", projectController.GetProjects)
	projects.GET("/:id", projectController.GetProjectByID)
	projects.GET("/filter", projectController.FilterProjects)
	projects.POST("", projectController.CreateProject)
	projects.PUT("/:id", projectController.UpdateProject)
	projects.DELETE("/:id", projectController.DeleteProject)

	articles := r.Group("/articles")
	articles.Use(middleware.AuthMiddleware()) // 🔥 Middleware JWT
	articles.GET("", articleController.GetAllArticles)
	articles.GET("/:id", articleController.GetArticleByID)
	articles.GET("/slug/:slug", articleController.GetArticleBySlug)
	articles.GET("/filter", articleController.FilterArticles)
	articles.POST("", articleController.CreateArticle)
	articles.PUT("/:id", articleController.UpdateArticle)
	articles.DELETE("/:id", articleController.DeleteArticle)
}
