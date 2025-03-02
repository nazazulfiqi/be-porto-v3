package controller

import (
	"be-porto-v3/dto"
	"be-porto-v3/models"
	"be-porto-v3/service"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

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
		dto.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch projects")
		return
	}
	dto.SuccessResponse(ctx, http.StatusOK, "Projects retrieved successfully", projects)
}

func (pc *ProjectController) GetProjectByID(ctx *gin.Context) {
	id := ctx.Param("id")

	project, err := pc.service.GetProjectByID(id)
	if err != nil {
		dto.ErrorResponse(ctx, http.StatusNotFound, "Project not found")
		return
	}

	dto.SuccessResponse(ctx, http.StatusOK, "Project retrieved successfully", project)
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
		dto.ErrorResponse(ctx, http.StatusBadRequest, "Invalid technologies format")
		return
	}
	if err := json.Unmarshal([]byte(ctx.PostForm("tags")), &project.Tags); err != nil {
		dto.ErrorResponse(ctx, http.StatusBadRequest, "Invalid tags format")
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
		dto.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create project")
		return
	}

	dto.SuccessResponse(ctx, http.StatusCreated, "Project created successfully", nil)
}

func (pc *ProjectController) FilterProjects(ctx *gin.Context) {
	title := ctx.Query("title")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	projects, totalData, err := pc.service.FilterProjects(title, page, limit)
	if err != nil {
		dto.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch projects")
		return
	}

	maxPage := (totalData + limit - 1) / limit // Hitung total halaman

	dto.SuccessPaginationResponse(ctx, http.StatusOK, "Projects retrieved successfully", projects, page, limit, maxPage, totalData)

}

func (pc *ProjectController) UpdateProject(ctx *gin.Context) {
	baseURL := os.Getenv("BASE_URL")

	var project models.Project

	// Ambil ID dari parameter URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		dto.ErrorResponse(ctx, http.StatusBadRequest, "Invalid project ID")
		return
	}
	project.ID = uint(id)

	// Bind data dari form
	project.Title = ctx.PostForm("title")
	project.Description = ctx.PostForm("description")
	project.RepositoryURL = ctx.PostForm("repository_url")
	project.DemoURL = ctx.PostForm("demo_url")
	project.Status = ctx.PostForm("status")

	// Parse JSON array dari form-data
	if err := json.Unmarshal([]byte(ctx.PostForm("technologies")), &project.Technologies); err != nil {
		dto.ErrorResponse(ctx, http.StatusBadRequest, "Invalid technologies format")
		return
	}
	if err := json.Unmarshal([]byte(ctx.PostForm("tags")), &project.Tags); err != nil {
		dto.ErrorResponse(ctx, http.StatusBadRequest, "Invalid tags format")
		return
	}

	// Handle upload gambar cover (opsional)
	file, err := ctx.FormFile("picture_cover")
	if err == nil {
		filePath := fmt.Sprintf("uploads/%s", file.Filename)
		ctx.SaveUploadedFile(file, filePath)
		project.PictureCover = fmt.Sprintf("%s/%s", baseURL, filePath)
	}

	// Handle multiple file uploads untuk pictures (opsional)
	form, _ := ctx.MultipartForm()
	files := form.File["pictures"]
	for _, file := range files {
		filePath := fmt.Sprintf("uploads/%s", file.Filename)
		ctx.SaveUploadedFile(file, filePath)
		project.Pictures = append(project.Pictures, fmt.Sprintf("%s/%s", baseURL, filePath))
	}

	// Update di database
	if err := pc.service.UpdateProject(&project); err != nil {
		if err.Error() == "project not found" {
			dto.ErrorResponse(ctx, http.StatusNotFound, "Project not found")
			return
		}
		dto.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update project")
		return
	}

	dto.SuccessResponse(ctx, http.StatusOK, "Project updated successfully", nil)
}

func (pc *ProjectController) DeleteProject(ctx *gin.Context) {
	id := ctx.Param("id")
	err := pc.service.DeleteProject(id)

	if err != nil {
		if err.Error() == "project not found" {
			dto.ErrorResponse(ctx, http.StatusNotFound, "Project not found")
			return
		}
		dto.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete project")
		return
	}

	dto.SuccessResponse(ctx, http.StatusOK, "Project deleted successfully", nil)
}
