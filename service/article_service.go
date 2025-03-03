package service

import (
	"be-porto-v3/models"
	"be-porto-v3/repository"
	"errors"

	"github.com/gosimple/slug"
)

type ArticleService interface {
	GetAllArticles() ([]models.Article, error)
	GetArticleByID(id uint) (*models.Article, error)
	CreateArticle(article *models.Article) error
	UpdateArticle(id uint, updatedData *models.Article) error
	DeleteArticle(id uint) error
	GetArticleBySlug(slug string) (*models.Article, error)
	FilterArticles(title string, page, limit int) ([]models.Article, int, error)
	IncrementViewCount(articleID uint)
}

type articleService struct {
	articleRepo repository.ArticleRepository
}

func NewArticleService(articleRepo repository.ArticleRepository) ArticleService {
	return &articleService{articleRepo: articleRepo}
}

func (s *articleService) GetAllArticles() ([]models.Article, error) {
	return s.articleRepo.GetAllArticles()
}

func (s *articleService) GetArticleByID(id uint) (*models.Article, error) {
	return s.articleRepo.GetArticleByID(id)
}

func (s *articleService) GetArticleBySlug(slug string) (*models.Article, error) {
	return s.articleRepo.GetArticleBySlug(slug)
}

func (s *articleService) FilterArticles(title string, page, limit int) ([]models.Article, int, error) {
	return s.articleRepo.FilterArticles(title, page, limit)
}

func (s *articleService) CreateArticle(article *models.Article) error {
	return s.articleRepo.CreateArticle(article)
}

func (s *articleService) UpdateArticle(id uint, updatedData *models.Article) error {
	article, err := s.articleRepo.GetArticleByID(id)
	if err != nil {
		return errors.New("article not found")
	}

	// Update hanya jika data baru tidak kosong
	if updatedData.Title != "" {
		article.Title = updatedData.Title
		article.Slug = slug.Make(updatedData.Title)
	}
	if updatedData.Content != "" {
		article.Content = updatedData.Content
	}
	if updatedData.Excerpt != "" {
		article.Excerpt = updatedData.Excerpt
	}
	if updatedData.SEOTitle != "" {
		article.SEOTitle = updatedData.SEOTitle
	}
	if updatedData.SEODescription != "" {
		article.SEODescription = updatedData.SEODescription
	}
	if updatedData.CoverImage != "" {
		article.CoverImage = updatedData.CoverImage
	}
	if updatedData.Status != "" {
		article.Status = updatedData.Status
	}
	if updatedData.PublishedAt != nil {
		article.PublishedAt = updatedData.PublishedAt
	}

	return s.articleRepo.UpdateArticle(article)
}

func (s *articleService) DeleteArticle(id uint) error {
	return s.articleRepo.DeleteArticle(id)
}

func (s *articleService) IncrementViewCount(articleID uint) {
	s.articleRepo.IncrementViewCount(articleID)
}
