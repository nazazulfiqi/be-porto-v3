package controller

import (
	"be-porto-v3/dto"
	"be-porto-v3/models"
	"be-porto-v3/service"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

type ArticleController struct {
	service service.ArticleService
}

func NewArticleController(articleService service.ArticleService) *ArticleController {
	return &ArticleController{service: articleService}
}

// Get All Articles
func (ac *ArticleController) GetAllArticles(ctx *gin.Context) {
	articles, err := ac.service.GetAllArticles()
	if err != nil {
		dto.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch articles")
		return
	}
	dto.SuccessResponse(ctx, http.StatusOK, "Articles retrieved successfully", articles)
}

// Get Article by ID
func (ac *ArticleController) GetArticleByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		dto.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
		return
	}

	article, err := ac.service.GetArticleByID(uint(id))
	if err != nil {
		dto.ErrorResponse(ctx, http.StatusNotFound, "Article not found")
		return
	}

	go ac.service.IncrementViewCount(article.ID)

	dto.SuccessResponse(ctx, http.StatusOK, "Article retrieved successfully", article)
}

// Get Article by Slug
func (ac *ArticleController) GetArticleBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")

	article, err := ac.service.GetArticleBySlug(slug)
	if err != nil {
		dto.ErrorResponse(ctx, http.StatusNotFound, "Article not found")
		return
	}

	go ac.service.IncrementViewCount(article.ID)

	dto.SuccessResponse(ctx, http.StatusOK, "Article retrieved successfully", article)
}

// Filter Articles
func (pc *ArticleController) FilterArticles(ctx *gin.Context) {
	title := ctx.Query("title")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	articles, totalData, err := pc.service.FilterArticles(title, page, limit)
	if err != nil {
		dto.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch articles")
		return
	}

	maxPage := (totalData + limit - 1) / limit // Hitung total halaman

	dto.SuccessPaginationResponse(ctx, http.StatusOK, "Articles retrieved successfully", articles, page, limit, maxPage, totalData)

}

// Create Article
func (ac *ArticleController) CreateArticle(ctx *gin.Context) {
	baseURL := os.Getenv("BASE_URL")

	var article models.Article
	article.Title = ctx.PostForm("title")
	article.Content = ctx.PostForm("content")
	article.Excerpt = ctx.PostForm("excerpt")
	article.SEOTitle = ctx.PostForm("seo_title")
	article.SEODescription = ctx.PostForm("seo_description")
	article.Status = ctx.PostForm("status")
	authorID := ctx.PostForm("author_id")

	// Convert AuthorID to uint
	var authorIDUint uint
	if _, err := fmt.Sscanf(authorID, "%d", &authorIDUint); err != nil {
		dto.ErrorResponse(ctx, http.StatusBadRequest, "Invalid author_id")
		return
	}
	article.AuthorID = authorIDUint

	// Generate Slug
	article.Slug = slug.Make(article.Title)

	// Parse Tags (JSON array)
	if err := json.Unmarshal([]byte(ctx.PostForm("tags")), &article.Tags); err != nil {
		dto.ErrorResponse(ctx, http.StatusBadRequest, "Invalid tags format")
		return
	}

	// Handle cover image upload
	file, err := ctx.FormFile("cover_image")
	if err == nil {
		uploadDir := "uploads/articles"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			dto.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create upload directory")
			return
		}

		fileName := fmt.Sprintf("%d-%s%s", time.Now().Unix(), slug.Make(file.Filename), filepath.Ext(file.Filename))
		filePath := filepath.Join(uploadDir, fileName)
		ctx.SaveUploadedFile(file, filePath)
		article.CoverImage = fmt.Sprintf("%s/%s", baseURL, filePath)
	}

	// Jika status adalah "published", atur PublishedAt
	if article.Status == "published" {
		now := time.Now()
		article.PublishedAt = &now
	}

	// Simpan ke database
	if err := ac.service.CreateArticle(&article); err != nil {
		dto.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create article")
		return
	}

	dto.SuccessResponse(ctx, http.StatusCreated, "Article created successfully", nil)
}

// Update Article
func (ac *ArticleController) UpdateArticle(ctx *gin.Context) {
	baseURL := os.Getenv("BASE_URL")

	var article models.Article

	// Ambil ID dari parameter URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		dto.ErrorResponse(ctx, http.StatusBadRequest, "Invalid article ID")
		return
	}
	article.ID = uint(id)

	// Bind data dari form
	article.Title = ctx.PostForm("title")
	article.Content = ctx.PostForm("content")
	article.Excerpt = ctx.PostForm("excerpt")
	article.SEOTitle = ctx.PostForm("seo_title")
	article.SEODescription = ctx.PostForm("seo_description")
	article.Status = ctx.PostForm("status")

	// Parse JSON array dari form-data untuk Tags
	if err := json.Unmarshal([]byte(ctx.PostForm("tags")), &article.Tags); err != nil {
		dto.ErrorResponse(ctx, http.StatusBadRequest, "Invalid tags format")
		return
	}

	// Handle upload gambar cover (opsional)
	file, err := ctx.FormFile("cover_image")
	if err == nil {
		// Buat folder jika belum ada
		uploadDir := "uploads/articles"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.MkdirAll(uploadDir, os.ModePerm)
		}

		// Simpan file
		filePath := fmt.Sprintf("%s/%s", uploadDir, file.Filename)
		ctx.SaveUploadedFile(file, filePath)
		article.CoverImage = fmt.Sprintf("%s/%s", baseURL, filePath)
	}

	// Update di database
	if err := ac.service.UpdateArticle(&article); err != nil {
		if err.Error() == "article not found" {
			dto.ErrorResponse(ctx, http.StatusNotFound, "Article not found")
			return
		}
		dto.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update article")
		return
	}

	dto.SuccessResponse(ctx, http.StatusOK, "Article updated successfully", nil)
}

// Delete Article
func (ac *ArticleController) DeleteArticle(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		dto.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
		return
	}

	if err := ac.service.DeleteArticle(uint(id)); err != nil {
		dto.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete article")
		return
	}

	dto.SuccessResponse(ctx, http.StatusOK, "Article deleted successfully", nil)
}
