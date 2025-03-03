package service

import (
	"be-porto-v3/models"
	"be-porto-v3/repository"
	"errors"
)

type ArticleService interface {
	GetAllArticles() ([]models.Article, error)
	GetArticleByID(id uint) (*models.Article, error)
	CreateArticle(article *models.Article) error
	UpdateArticle(article *models.Article) error
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

func (s *articleService) UpdateArticle(article *models.Article) error {
	existingArticle, err := s.articleRepo.GetArticleByID(article.ID)
	if err != nil {
		return errors.New("article not found")
	}

	// Pastikan slug tidak berubah
	article.Slug = existingArticle.Slug

	return s.articleRepo.UpdateArticle(article)
}

func (s *articleService) DeleteArticle(id uint) error {
	return s.articleRepo.DeleteArticle(id)
}

func (s *articleService) IncrementViewCount(articleID uint) {
	s.articleRepo.IncrementViewCount(articleID)
}
