package repository

import (
	"be-porto-v3/models"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	GetAllArticles() ([]models.Article, error)
	GetArticleByID(id uint) (*models.Article, error)
	CreateArticle(article *models.Article) error
	UpdateArticle(article *models.Article) error
	DeleteArticle(id uint) error
	GetArticleBySlug(slug string) (*models.Article, error)
	FilterArticles(title string, page, limit int) ([]models.Article, int, error)
	IncrementViewCount(articleID uint) error
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) GetAllArticles() ([]models.Article, error) {
	var articles []models.Article
	err := r.db.Order("created_at DESC").Find(&articles).Error
	return articles, err
}

func (r *articleRepository) GetArticleByID(id uint) (*models.Article, error) {
	var article models.Article
	err := r.db.First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) FilterArticles(title string, page, limit int) ([]models.Article, int, error) {
	var articles []models.Article
	var totalData int64

	query := r.db.Model(&models.Article{})

	// Filter by title (jika ada)
	if title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}

	// Hitung total data sebelum pagination
	if err := query.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	// Ambil data dengan pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, int(totalData), nil
}

func (r *articleRepository) GetArticleBySlug(slug string) (*models.Article, error) {
	var article models.Article
	result := r.db.Where("slug = ?", slug).First(&article)
	if result.Error != nil {
		return nil, result.Error
	}
	return &article, nil
}

func (r *articleRepository) CreateArticle(article *models.Article) error {
	return r.db.Create(article).Error
}

func (r *articleRepository) UpdateArticle(article *models.Article) error {
	return r.db.Save(article).Error
}

func (r *articleRepository) DeleteArticle(id uint) error {
	return r.db.Delete(&models.Article{}, id).Error
}

func (r *articleRepository) IncrementViewCount(articleID uint) error {
	return r.db.Model(&models.Article{}).
		Where("id = ?", articleID).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).
		Error
}
