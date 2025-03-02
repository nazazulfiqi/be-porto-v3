package controller

import (
	"be-porto-v3/models"
	"be-porto-v3/service"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type ProjectController struct {
	service service.ProjectService
}

func NewProjectController(projectService service.ProjectService) *ProjectController {
	return &ProjectController{service: projectService}
}

func (pc *ProjectController) GetProjects(ctx *gin.Context) {
	projects, err := pc.service.GetProjects()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, projects)
}

func (pc *ProjectController) GetProjectByID(ctx *gin.Context) {
	id := ctx.Param("id")

	project, err := pc.service.GetProjectByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	ctx.JSON(http.StatusOK, project)
}

func (pc *ProjectController) CreateProject(ctx *gin.Context) {
	baseURL := os.Getenv("BASE_URL")

	var project models.Project

	// Bind text fields dari form-data
	project.Title = ctx.PostForm("title")
	project.Description = ctx.PostForm("description")
	project.RepositoryURL = ctx.PostForm("repository_url")
	project.DemoURL = ctx.PostForm("demo_url")
	project.Status = ctx.PostForm("status")

	// Parse JSON array dari form-data
	if err := json.Unmarshal([]byte(ctx.PostForm("technologies")), &project.Technologies); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid technologies format"})
		return
	}
	if err := json.Unmarshal([]byte(ctx.PostForm("tags")), &project.Tags); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tags format"})
		return
	}

	// Handle upload gambar cover
	file, err := ctx.FormFile("picture_cover")
	if err == nil {
		filePath := fmt.Sprintf("uploads/%s", file.Filename)
		ctx.SaveUploadedFile(file, filePath)
		project.PictureCover = fmt.Sprintf("%s/%s", baseURL, filePath) // Buat URL lengkap
	}

	// Handle multiple file uploads untuk pictures
	form, _ := ctx.MultipartForm()
	files := form.File["pictures"]
	for _, file := range files {
		filePath := fmt.Sprintf("uploads/%s", file.Filename)
		ctx.SaveUploadedFile(file, filePath)
		project.Pictures = append(project.Pictures, fmt.Sprintf("%s/%s", baseURL, filePath)) // Buat URL lengkap
	}

	// Simpan ke database
	if err := pc.service.CreateProject(&project); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, project)
}

func (pc *ProjectController) UpdateProject(ctx *gin.Context) {
	var project models.Project
	if err := ctx.ShouldBindJSON(&project); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := pc.service.UpdateProject(&project); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, project)
}

func (pc *ProjectController) DeleteProject(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := pc.service.DeleteProject(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Project deleted"})
}
