package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/h4rdc0m/aurora-api/domain/errors"
	"strconv"
)

type BaseController struct {
}

func (b *BaseController) GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}

	return userID.(string), true
}

func (b *BaseController) GetUserEmail(c *gin.Context) (string, bool) {
	userEmail, exists := c.Get("user_email")
	if !exists {
		return "", false
	}

	return userEmail.(string), true
}

func (b *BaseController) GetUserRoles(c *gin.Context) ([]string, bool) {
	roles, exists := c.Get("user_roles")
	if !exists {
		return nil, false
	}

	return roles.([]string), true
}

func (b *BaseController) HasRole(c *gin.Context, role string) bool {
	roles, exists := b.GetUserRoles(c)
	if !exists {
		return false
	}

	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func (b *BaseController) ParseUIntParam(c *gin.Context, param string) (uint, error) {
	paramID := c.Param(param)
	if paramID == "" {
		return 0, errors.ErrInvalidID
	}

	id, err := strconv.ParseUint(paramID, 10, 32)
	if err != nil {
		return 0, errors.ErrInvalidID
	}

	return uint(id), nil
}
