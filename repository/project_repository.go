package repository

import (
	"be-porto-v3/models"

	"gorm.io/gorm"
)

type ProjectRepository interface {
	GetProjects() ([]models.Project, error)
	CreateProject(*models.Project) error
	UpdateProject(*models.Project) error
	DeleteProject(id string) error
	GetProjectByID(id string) (*models.Project, error)
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

func (r *projectRepository) CreateProject(project *models.Project) error {
	return r.db.Create(project).Error
}

func (r *projectRepository) UpdateProject(project *models.Project) error {
	return r.db.Save(project).Error
}

func (r *projectRepository) DeleteProject(id string) error {
	return r.db.Delete(&models.Project{}, id).Error
}
