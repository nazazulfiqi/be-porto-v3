package repository

import (
	"be-porto-v3/models"
	"fmt"

	"gorm.io/gorm"
)

type ProjectRepository interface {
	GetProjects() ([]models.Project, error)
	CreateProject(*models.Project) error
	UpdateProject(*models.Project) error
	DeleteProject(id string) error
	GetProjectByID(id string) (*models.Project, error)
	FilterProjects(title string, page, limit int) ([]models.Project, int, error)
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) GetProjects() ([]models.Project, error) {
	var projects []models.Project
	if err := r.db.Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepository) GetProjectByID(id string) (*models.Project, error) {
	var project models.Project
	err := r.db.First(&project, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) FilterProjects(title string, page, limit int) ([]models.Project, int, error) {
	var projects []models.Project
	var totalData int64

	query := r.db.Model(&models.Project{})

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
	if err := query.Offset(offset).Limit(limit).Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	return projects, int(totalData), nil
}

func (r *projectRepository) CreateProject(project *models.Project) error {
	return r.db.Create(project).Error
}

func (r *projectRepository) UpdateProject(project *models.Project) error {
	var existingProject models.Project
	if err := r.db.First(&existingProject, project.ID).Error; err != nil {
		return fmt.Errorf("project not found")
	}

	// Gunakan data lama jika field kosong
	if project.Title == "" {
		project.Title = existingProject.Title
	}
	if project.Description == "" {
		project.Description = existingProject.Description
	}
	if len(project.Technologies) == 0 {
		project.Technologies = existingProject.Technologies
	}
	if project.PictureCover == "" {
		project.PictureCover = existingProject.PictureCover
	}
	if len(project.Pictures) == 0 {
		project.Pictures = existingProject.Pictures
	}
	if project.RepositoryURL == "" {
		project.RepositoryURL = existingProject.RepositoryURL
	}
	if project.DemoURL == "" {
		project.DemoURL = existingProject.DemoURL
	}
	if project.Status == "" {
		project.Status = existingProject.Status
	}
	if len(project.Tags) == 0 {
		project.Tags = existingProject.Tags
	}

	// Simpan perubahan
	return r.db.Save(project).Error
}

func (r *projectRepository) DeleteProject(id string) error {
	result := r.db.Delete(&models.Project{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("project not found")
	}
	return nil
}
