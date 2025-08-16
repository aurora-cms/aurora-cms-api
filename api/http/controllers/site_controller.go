package controllers

//import (
//	"aurora-api/internal/dto"
//	"aurora-api/internal/services"
//	"aurora-api/internal/utils/httpx"
//	"github.com/gin-gonic/gin"
//	"github.com/go-playground/validator"
//	"github.com/rs/zerolog/log"
//	"net/http"
//)
//
//type SiteController struct {
//	siteService services.SiteService
//	validate    *validator.Validate
//}
//
//func NewSiteController(siteService services.SiteService, group *gin.RouterGroup) *SiteController {
//	log.Info().Msg("Initializing site controller")
//	c := &SiteController{
//		siteService: siteService,
//		validate:    validator.New(),
//	}
//	c.registerRoutes(group)
//	return c
//}
//
//func (c *SiteController) registerRoutes(r *gin.RouterGroup) {
//	siteApi := r.Group("/sites")
//	siteApi.GET("", c.getAllSites)
//	siteApi.GET("/:id", c.getSiteById)
//	siteApi.POST("", c.createSite)
//	siteApi.PUT("/:id", c.updateSite)
//	siteApi.DELETE("/:id", c.deleteSite)
//}
//
//func (c *SiteController) getAllSites(ctx *gin.Context) {
//	sites, err := c.siteService.GetAll(ctx.Request.Context())
//	if err != nil {
//		log.Error().Err(err).Msg("Failed to retrieve sites")
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	ctx.JSON(http.StatusOK, gin.H{"sites": sites})
//}
//
//func (c *SiteController) getSiteById(ctx *gin.Context) {
//	id, err := httpx.ParseIntID(ctx)
//	if err != nil {
//		return
//	}
//
//	s, err := c.siteService.GetByID(ctx.Request.Context(), id)
//	if err != nil {
//		httpx.HandleError(ctx, err)
//	}
//
//	ctx.JSON(http.StatusOK, s)
//}
//
//func (c *SiteController) createSite(ctx *gin.Context) {
//	var input dto.CreateSiteRequest
//	if err := ctx.ShouldBindJSON(&input); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON", "details": err.Error()})
//		return
//	}
//
//	if err := c.validate.Struct(input); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
//		return
//	}
//
//	if err := c.siteService.Create(ctx.Request.Context(), &input); err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create site"})
//	}
//
//	ctx.JSON(http.StatusCreated, gin.H{"message": "Site created successfully"})
//}
//
//func (c *SiteController) updateSite(ctx *gin.Context) {
//	id, err := httpx.ParseIntID(ctx)
//	if err != nil {
//		return
//	}
//
//	var input dto.UpdateSiteRequest
//	if err := ctx.ShouldBindJSON(&input); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON", "details": err.Error()})
//		return
//	}
//
//	if err := c.validate.Struct(input); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
//		return
//	}
//	if err := c.siteService.Update(ctx.Request.Context(), id, &input); err != nil {
//		httpx.HandleError(ctx, err)
//		return
//	}
//
//	ctx.JSON(http.StatusNoContent, gin.H{})
//}
//
//func (c *SiteController) deleteSite(context *gin.Context) {
//	id, err := httpx.ParseIntID(context)
//	if err != nil {
//		return
//	}
//
//	if err := c.siteService.Delete(context.Request.Context(), id); err != nil {
//		httpx.HandleError(context, err)
//		return
//	}
//
//	context.JSON(http.StatusNoContent, gin.H{})
//}
