package service

import (
	"be-porto-v3/models"
	"be-porto-v3/repository"
	"errors"
	"sync"
	"time"
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
	LikeArticle(userKey string, articleID uint) error
}

type articleService struct {
	articleRepo repository.ArticleRepository
	likeCache   map[string]time.Time // Menyimpan IP/userID terakhir like
	mu          sync.Mutex           // Mutex untuk menghindari race condition
}

func NewArticleService(articleRepo repository.ArticleRepository) ArticleService {
	return &articleService{
		articleRepo: articleRepo,
		likeCache:   make(map[string]time.Time),
	}
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

func (s *articleService) LikeArticle(userKey string, articleID uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Cek apakah user baru saja like
	lastLike, exists := s.likeCache[userKey]
	if exists && time.Since(lastLike) < 5*time.Second {
		return errors.New("please wait before liking again")
	}

	// Simpan waktu like terbaru
	s.likeCache[userKey] = time.Now()

	// Tambahkan counter like di database
	return s.articleRepo.IncrementLike(articleID)
}
