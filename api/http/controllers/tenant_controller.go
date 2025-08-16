package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/h4rdc0m/aurora-api/application/dto"
	"github.com/h4rdc0m/aurora-api/application/use_cases"
	"github.com/h4rdc0m/aurora-api/domain/common"
	"net/http"
)

// TenantController handles HTTP requests related to tenants.
type TenantController struct {
	BaseController
	tenantUseCase *use_cases.TenantUseCase
	logger        common.Logger
}

// NewTenantController creates a new instance of TenantController with the provided service and logger.
func NewTenantController(tenantUseCase *use_cases.TenantUseCase, logger common.Logger) *TenantController {
	return &TenantController{
		tenantUseCase: tenantUseCase,
		logger:        logger,
	}
}

// GetOneTenant retrieves a single tenant by its ID.
func (t *TenantController) GetOneTenant(c *gin.Context) {
	id, err := t.ParseUIntParam(c, "id")
	if err != nil {
		t.logger.Error("Failed to parse tenant ID", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	tenant, err := t.tenantUseCase.GetTenant(uint64(id))
	if err != nil {
		t.logger.Error("Failed to get tenant", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if tenant == nil {
		t.logger.Error("Tenant not found", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tenant,
	})
}

// GetAllTenants retrieves all tenants from the database.
func (t *TenantController) GetAllTenants(c *gin.Context) {
	tenants, err := t.tenantUseCase.GetAllTenants()
	if err != nil {
		t.logger.Error("Failed to get tenants", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tenants})
}

// GetActiveTenants retrieves only active tenants
func (t *TenantController) GetActiveTenants(c *gin.Context) {
	tenants, err := t.tenantUseCase.GetActiveTenants()
	if err != nil {
		t.logger.Error("Failed to get active tenants", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tenants})
}

// CreateTenant creates a new tenant
func (t *TenantController) CreateTenant(c *gin.Context) {
	var req dto.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		t.logger.Error("Failed to bind JSON to tenant request", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenant, err := t.tenantUseCase.CreateTenant(req.Name, &req.Description)
	if err != nil {
		t.logger.Error("Failed to create tenant", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": tenant})
}

// UpdateTenant updates an existing tenant
func (t *TenantController) UpdateTenant(c *gin.Context) {
	id, err := t.ParseUIntParam(c, "id")
	if err != nil {
		t.logger.Error("Failed to parse tenant ID", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var req dto.UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		t.logger.Error("Failed to bind JSON to tenant request", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenant, err := t.tenantUseCase.UpdateTenant(uint64(id), req.Name, &req.Description)
	if err != nil {
		t.logger.Error("Failed to update tenant", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tenant})
}

// DeleteTenant deletes a tenant by its ID
func (t *TenantController) DeleteTenant(c *gin.Context) {
	id, err := t.ParseUIntParam(c, "id")
	if err != nil {
		t.logger.Error("Failed to parse tenant ID", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	err = t.tenantUseCase.DeleteTenant(uint64(id))
	if err != nil {
		t.logger.Error("Failed to delete tenant", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "Tenant deleted successfully"})
}
