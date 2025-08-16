package controllers

//
//import (
//	"aurora-api/internal/dto"
//	"aurora-api/internal/services"
//	"aurora-api/internal/utils/httpx"
//	"github.com/gin-gonic/gin"
//	"github.com/go-playground/validator"
//	"net/http"
//)
//
//type PageController struct {
//	pageService services.PageService
//	validate    *validator.Validate
//}
//
//func NewPageController(pageService services.PageService, group *gin.RouterGroup) *PageController {
//	c := &PageController{
//		pageService: pageService,
//		validate:    validator.New(),
//	}
//	c.registerRoutes(group)
//	return c
//}
//
//func (c *PageController) registerRoutes(r *gin.RouterGroup) {
//	pageApi := r.Group("/pages")
//	pageApi.GET("", c.GetPages)
//	pageApi.GET("/:id", c.GetPage)
//	pageApi.POST("", c.CreatePage)
//	pageApi.PUT("/:id", c.UpdatePage)
//	pageApi.DELETE("/:id", c.DeletePage)
//}
//
//func (c *PageController) GetPages(ctx *gin.Context) {
//	// Logic to get pages
//	pages, err := c.pageService.GetAll(ctx.Request.Context())
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	ctx.JSON(http.StatusOK, gin.H{"pages": pages})
//}
//
//func (c *PageController) GetPage(ctx *gin.Context) {
//	id, err := httpx.ParseIntID(ctx)
//	if err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
//		return
//	}
//
//	page, err := c.pageService.GetByID(ctx.Request.Context(), id)
//	if err != nil {
//		httpx.HandleError(ctx, err)
//		return
//	}
//
//	ctx.JSON(http.StatusOK, page)
//}
//
//func (c *PageController) CreatePage(ctx *gin.Context) {
//	// Logic to create a page
//	var input dto.CreatePageRequest
//	if err := ctx.ShouldBindJSON(&input); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
//		return
//	}
//
//	// Validate the input
//	if err := c.validate.Struct(&input); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
//		return
//	}
//
//	if err := c.pageService.Create(ctx.Request.Context(), &input); err != nil {
//		httpx.HandleError(ctx, err)
//		return
//	}
//
//	ctx.Status(http.StatusCreated)
//}
//
//func (c *PageController) UpdatePage(ctx *gin.Context) {
//	id, err := httpx.ParseIntID(ctx)
//	if err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
//		return
//	}
//
//	// Validate the ID
//	if id <= 0 {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
//		return
//	}
//
//	var input dto.UpdatePageRequest
//	if err := ctx.ShouldBindJSON(&input); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
//		return
//	}
//	// Validate the input
//	if err := c.validate.Struct(&input); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
//		return
//	}
//
//	if err := c.pageService.Update(ctx.Request.Context(), &input); err != nil {
//		httpx.HandleError(ctx, err)
//		return
//	}
//
//	ctx.Status(http.StatusNoContent)
//}
//
//func (c *PageController) DeletePage(ctx *gin.Context) {
//	id, err := httpx.ParseIntID(ctx)
//	if err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
//		return
//	}
//
//	// Validate the ID
//	if c.validate.Var(id, "gt=0") != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
//		return
//	}
//
//	if err := c.pageService.Delete(ctx.Request.Context(), id); err != nil {
//		httpx.HandleError(ctx, err)
//		return
//	}
//
//	ctx.Status(http.StatusNoContent)
//}
