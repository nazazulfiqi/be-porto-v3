package service

import (
	"be-porto-v3/models"
	"be-porto-v3/repository"
)

type ProjectService interface {
	GetProjects() ([]models.Project, error)
	CreateProject(*models.Project) error
	UpdateProject(*models.Project) error
	DeleteProject(id string) error
	GetProjectByID(id string) (*models.Project, error)
	FilterProjects(title string, page, limit int) ([]models.Project, int, error)
}

type projectService struct {
	projectRepo repository.ProjectRepository
}

func NewProjectService(projectRepo repository.ProjectRepository) ProjectService {
	return &projectService{projectRepo: projectRepo}
}

func (s *projectService) GetProjects() ([]models.Project, error) {
	return s.projectRepo.GetProjects()
}

func (s *projectService) GetProjectByID(id string) (*models.Project, error) {
	return s.projectRepo.GetProjectByID(id)
}

func (s *projectService) FilterProjects(title string, page, limit int) ([]models.Project, int, error) {
	return s.projectRepo.FilterProjects(title, page, limit)
}

func (s *projectService) CreateProject(project *models.Project) error {
	return s.projectRepo.CreateProject(project)
}

func (s *projectService) UpdateProject(project *models.Project) error {
	return s.projectRepo.UpdateProject(project)
}

func (s *projectService) DeleteProject(id string) error {
	return s.projectRepo.DeleteProject(id)
}
