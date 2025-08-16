package controllers

//
//import (
//	"aurora-api/internal/services"
//	"errors"
//	"github.com/gin-gonic/gin"
//	"github.com/rs/zerolog/log"
//	"net/url"
//	"regexp"
//	"strings"
//)
//
//var domainRegex = regexp.MustCompile(`^([a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$`)
//
//type DomainController struct {
//	siteService   services.SiteService
//	tenantService services.TenantService
//	pageService   services.PageService
//}
//
//func NewDomainController(siteService services.SiteService, tenantService services.TenantService, pageService services.PageService, group *gin.RouterGroup) *DomainController {
//	c := &DomainController{
//		siteService:   siteService,
//		tenantService: tenantService,
//		pageService:   pageService,
//	}
//
//	c.registerRoutes(group)
//
//	return c
//}
//
//func (c *DomainController) registerRoutes(group *gin.RouterGroup) {
//	api := group.Group("/domains")
//
//	api.GET("/:domain", c.GetDomain)
//}
//
//func (c *DomainController) GetDomain(ctx *gin.Context) {
//	domain, err := sanitizeDomain(ctx.Param("domain"))
//	if err != nil {
//		ctx.JSON(400, gin.H{"error": err.Error()})
//	}
//
//	if domain == nil {
//		ctx.JSON(400, gin.H{"error": "invalid domain format"})
//		return
//	}
//
//	site, err := c.siteService.GetSiteByDomain(ctx, *domain)
//	if err != nil {
//		ctx.JSON(404, gin.H{"error": "site not found"})
//		return
//	}
//	log.Debug().Msgf("Found site for domain %s", *domain)
//	if site == nil {
//		ctx.JSON(404, gin.H{"error": "site not found"})
//		return
//	}
//	ctx.JSON(200, site)
//}
//
//func sanitizeDomain(domain string) (*string, error) {
//	// Sanitize the domain
//	domain = strings.TrimSpace(domain)     // Remove whitespace
//	domain = strings.ToLower(domain)       // Convert to lowercase
//	domain = strings.Split(domain, "/")[0] // Remove the path
//	domain = strings.Split(domain, "?")[0] // Remove query parameters
//	domain = strings.Split(domain, "#")[0] // Remove fragments
//
//	// Validate the domain format
//	if !domainRegex.MatchString(domain) {
//		// If not valid, try to parse as URL to extract the hostname
//		if u, err := url.Parse("http://" + domain); err == nil {
//			domain = u.Hostname()
//		}
//
//		// If still not valid, return nil and an error
//		if !domainRegex.MatchString(domain) {
//			return nil, errors.New("invalid domain format")
//		}
//	}
//
//	return &domain, nil
//}
